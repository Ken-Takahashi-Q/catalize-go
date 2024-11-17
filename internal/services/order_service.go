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
	order.OrderStatus = models.Received

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
	collection :=
		db.GetCollection(models.GetCollection{DBName: "order", Collection: "history"})

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
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

	cursor, err := collection.Find(ctx, filter, options.Find().SetSort(bson.M{"OrderDate": -1}))
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
