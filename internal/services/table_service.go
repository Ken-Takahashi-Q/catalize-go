package services

import (
	"catalize-go/internal/db"
	"catalize-go/internal/models"
	"catalize-go/internal/models/table"
	"catalize-go/internal/repositories"
	"context"
	"fmt"
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

func CreateTableVisitService(generateTableVisitID models.Table) (int, error) {
	collection := db.GetCollection(models.GetCollection{DBName: "order", Collection: "table"})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	truncatedDate := generateTableVisitID.CreatedAt.Truncate(24 * time.Hour)

	var latestOrderTable models.Table
	err := collection.FindOne(ctx, bson.M{
		"TableID": generateTableVisitID.TableID,
	}, options.FindOne().SetSort(bson.M{"Date": -1})).Decode(&latestOrderTable)

	newTable := models.Table{
		CreatedAt:       generateTableVisitID.CreatedAt,
		TableID:    generateTableVisitID.TableID,
		TableVisit: int(latestOrderTable.TableVisit) + 1,
	}

	if err != nil && err != mongo.ErrNoDocuments {
		return 0, err
	}
	if err == mongo.ErrNoDocuments || truncatedDate.After(latestOrderTable.CreatedAt.Truncate(24*time.Hour)) {
		newTable.TableVisit = 1
	}

	_, err = collection.InsertOne(ctx, newTable)
	if err != nil {
		return 0, err
	}

	return newTable.TableVisit, nil
}

func GetTableVisitService(generateTableVisitID models.Table) (int, error) {
	collection := db.GetCollection(models.GetCollection{DBName: "order", Collection: "table"})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var latestTable models.Table
	err := collection.FindOne(ctx, bson.M{
		"TableID": generateTableVisitID.TableID,
	}, options.FindOne().SetSort(bson.M{"Date": -1})).Decode(&latestTable)

	if err != nil {
		return 0, err
	}

	return int(latestTable.TableVisit), nil
}

func CheckBillService(orderTable models.Table) (table.TableCheckBill, error) {
	// collection := db.GetCollection(models.GetCollection{DBName: "order", Collection: "table"})

	// ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	// defer cancel()

	return table.TableCheckBill{}, nil
}

func GetTablesService() ([]models.Table, error) {
	collection := db.GetCollection(models.GetCollection{DBName: "order", Collection: "table"})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var tables []models.Table
	if err := cursor.All(ctx, &tables); err != nil {
		return nil, err
	}

	return tables, nil
}

func UpdateTableService(table models.Table) (models.Table, error) {
	collection := db.GetCollection(models.GetCollection{DBName: "order", Collection: "table"})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"TableID": table.TableID, "TableVisit": table.TableVisit}
	update := bson.M{"$set": table}

	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return models.Table{}, fmt.Errorf("failed to update table: %v", err)
	}

	if result.MatchedCount == 0 {
		return models.Table{}, fmt.Errorf("no table found with TableID: %s and TableVisit: %d", table.TableID, table.TableVisit)
	}

	var updatedTable models.Table
	if err := collection.FindOne(ctx, filter).Decode(&updatedTable); err != nil {
		return models.Table{}, fmt.Errorf("failed to retrieve updated table: %v", err)
	}

	return updatedTable, nil
}