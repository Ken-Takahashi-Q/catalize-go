package db

import (
	"catalize-go/internal/models"
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

func Connect(url string) error {
	var err error

	client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(url))
	if err != nil {
		log.Println("Error connecting to MongoDB:", err)
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Println("Failed to ping MongoDB:", err)
		return err
	}

	return nil
}

func GetClient() *mongo.Client {
	return client
}

func GetCollection(getCollection models.GetCollection) *mongo.Collection {
	return client.Database(getCollection.DBName).Collection(getCollection.Collection)
}
