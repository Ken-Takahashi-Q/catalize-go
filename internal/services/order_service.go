package services

import (
	"catalize-go/internal/db"
	"catalize-go/internal/models"
	"catalize-go/internal/models/body"
	"catalize-go/internal/repositories"
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type OrderService struct {
	orderRepo *repositories.OrderRepo
}

func NewOrderService(repo *repositories.OrderRepo) *OrderService {
	return &OrderService{
		orderRepo: repo,
	}
}

func CreateTableVisitService(generateTableVisitID models.OrderTable) (int, error) {
	collection := db.GetCollection(models.GetCollection{DBName: "order", Collection: "table"})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	truncatedDate := generateTableVisitID.Date.Truncate(24 * time.Hour)

	var latestOrderTable models.OrderTable
	err := collection.FindOne(ctx, bson.M{
		"TableID": generateTableVisitID.TableID,
	}, options.FindOne().SetSort(bson.M{"Date": -1})).Decode(&latestOrderTable)

	newOrderTable := models.OrderTable{
		Date:       generateTableVisitID.Date,
		TableID:    generateTableVisitID.TableID,
		TableVisit: int(latestOrderTable.TableVisit) + 1,
	}

	if err != nil && err != mongo.ErrNoDocuments {
		return 0, err
	}
	if err == mongo.ErrNoDocuments || truncatedDate.After(latestOrderTable.Date.Truncate(24*time.Hour)) {
		newOrderTable.TableVisit = 1
	}

	_, err = collection.InsertOne(ctx, newOrderTable)
	if err != nil {
		return 0, err
	}

	return newOrderTable.TableVisit, nil
}

func GetTableVisitService(generateTableVisitID models.OrderTable) (int, error) {
	collection := db.GetCollection(models.GetCollection{DBName: "order", Collection: "table"})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var latestOrderTable models.OrderTable
	err := collection.FindOne(ctx, bson.M{
		"TableID": generateTableVisitID.TableID,
	}, options.FindOne().SetSort(bson.M{"Date": -1})).Decode(&latestOrderTable)

	if err != nil {
		return 0, err
	}

	return int(latestOrderTable.TableVisit), nil
}

func CreateOrderService(order *models.Order) error {
	orderCount, err := CreateOrderCountService(models.UserOrder{
		TableID:    order.TableID,
		TableVisit: order.TableVisit,
		OrderCount: order.OrderCount,
	})
	if err != nil {
		return fmt.Errorf("failed to generate order ID: %w", err)
	}
	order.OrderCount = orderCount
	order.CreatedAt = time.Now()
	order.OrderDate = time.Now()
	order.TotalPrice = CalculateOrderPriceService(order)
	order.OrderStatus = models.Preparing

	collection := db.GetCollection(models.GetCollection{DBName: "order", Collection: "history"})
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := collection.InsertOne(ctx, order)
	if err != nil {
		return fmt.Errorf("failed to insert order: %w", err)
	}

	order.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

func CalculateOrderPriceService(order *models.Order) float64 {
	var totalPrice float64

	for _, item := range order.Items {
		totalPrice += item.Price * float64(item.Quantity)
	}

	return totalPrice
}

func CreateOrderCountService(generateOrderID models.UserOrder) (int, error) {
	collection :=
		db.GetCollection(models.GetCollection{DBName: "order", Collection: "user"})

	var userOrder models.UserOrder
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := collection.FindOne(ctx, bson.M{
		"LastOrder":  bson.M{"$gte": time.Now().Truncate(24 * time.Hour)},
		"TableID":    generateOrderID.TableID,
		"TableVisit": generateOrderID.TableVisit,
	}).Decode(&userOrder)

	if err == mongo.ErrNoDocuments {
		userOrder = models.UserOrder{
			LastOrder:  time.Now(),
			TableID:    generateOrderID.TableID,
			TableVisit: generateOrderID.TableVisit,
			OrderCount: 1,
		}

		if _, insertErr := collection.InsertOne(ctx, userOrder); insertErr != nil {
			return 0, insertErr
		}
		return userOrder.OrderCount, nil
	} else if err != nil {
		return 0, err
	}
	userOrder.OrderCount++

	_, err = collection.UpdateOne(ctx, bson.M{
		"TableID":    generateOrderID.TableID,
		"TableVisit": generateOrderID.TableVisit,
	}, bson.M{
		"$set": bson.M{
			"LastOrder":  time.Now(),
			"OrderCount": userOrder.OrderCount,
		},
	})
	if err != nil {
		return 0, err
	}

	return userOrder.OrderCount, nil
}

func GetOrdersService(getOrders body.GetOrders) ([]models.Order, error) {
	collection := db.GetCollection(models.GetCollection{DBName: "order", Collection: "history"})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{}
	if !getOrders.ID.IsZero() {
		filter["_id"] = getOrders.ID
	}
	if getOrders.UserID != 0 {
		filter["UserID"] = getOrders.UserID
	}
	if getOrders.TableID != 0 {
		filter["TableID"] = getOrders.TableID
	}
	if getOrders.TableVisit != 0 {
		filter["TableVisit"] = getOrders.TableVisit
	}
	if getOrders.OrderID != 0 {
		filter["OrderID"] = getOrders.OrderID
	}
	if !getOrders.OrderDate.IsZero() {
		filter["OrderDate"] = bson.M{
			"$gte": getOrders.OrderDate.Truncate(24 * time.Hour),
			"$lt":  getOrders.OrderDate.Truncate(24 * time.Hour).Add(24 * time.Hour),
		}
	}
	if getOrders.OrderStatus != 0 {
		filter["OrderStatus"] = getOrders.OrderStatus
	}

	// Filter the latest TableVisit for each TableID
	pipeline := []bson.M{
		// Match the filter criteria
		{"$match": filter},
		// Identify the latest TableVisit for each TableID
		{
			"$group": bson.M{
				"_id":         "$TableID",
				"latestVisit": bson.M{"$max": "$TableVisit"}, // Find the latest TableVisit
			},
		},
		// Join back with the original collection to get all documents for the latest TableVisit
		{
			"$lookup": bson.M{
				"from": "history",
				"let":  bson.M{"tableID": "$_id", "latestVisit": "$latestVisit"},
				"pipeline": []bson.M{
					{"$match": bson.M{
						"$expr": bson.M{
							"$and": []bson.M{
								{"$eq": []interface{}{"$TableID", "$$tableID"}},
								{"$eq": []interface{}{"$TableVisit", "$$latestVisit"}},
							},
						},
					}},
				},
				"as": "orders",
			},
		},
		// Unwind to return individual orders instead of grouped arrays
		{"$unwind": "$orders"},
		// Replace the root to show the actual order document
		{"$replaceRoot": bson.M{"newRoot": "$orders"}},
	}

	if getOrders.OrderStatus != 0 {
		pipeline = append(pipeline, bson.M{
			"$match": bson.M{
				"OrderStatus": getOrders.OrderStatus,
			},
		})
	}

	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var orders []models.Order
	if err := cursor.All(ctx, &orders); err != nil {
		return nil, err
	}

	return orders, nil
}

func KitchenDoneOrderService(getOrders body.GetOrders) (models.Order, error) {
	collection := db.GetCollection(models.GetCollection{DBName: "order", Collection: "history"})

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	orderID, err := primitive.ObjectIDFromHex(getOrders.ID.Hex())
	if err != nil {
		return models.Order{}, fmt.Errorf("invalid order ID: %w", err)
	}
	update := bson.M{
		"$set": bson.M{
			"OrderStatus": models.Ready,
		},
	}

	result := collection.FindOneAndUpdate(ctx, bson.M{"_id": orderID}, update, options.FindOneAndUpdate().SetReturnDocument(options.After))
	if result.Err() != nil {
		return models.Order{}, fmt.Errorf("failed to update order status: %w", result.Err())
	}

	var updatedOrder models.Order
	if err := result.Decode(&updatedOrder); err != nil {
		return models.Order{}, fmt.Errorf("failed to decode updated order: %w", err)
	}

	return updatedOrder, nil
}

func ClearAllOrderService() error {
	collection :=
		db.GetCollection(models.GetCollection{DBName: "order", Collection: "history"})

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := collection.DeleteMany(ctx, bson.M{})
	if err != nil {
		return err
	}

	return nil
}
