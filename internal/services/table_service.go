package services

import (
	"catalize-go/internal/db"
	"catalize-go/internal/models"
	"catalize-go/internal/models/table"
	"catalize-go/internal/repositories"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type TableService struct {
	tableRepo *repositories.TableRepo
}

func NewTableService(repo *repositories.TableRepo) *TableService {
	return &TableService{
		tableRepo: repo,
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

func CheckBillService(orderTable models.OrderTable) (table.TableCheckBill, error) {
	// collection := db.GetCollection(models.GetCollection{DBName: "order", Collection: "table"})

	// ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	// defer cancel()

	// return orderTable, nil
}
