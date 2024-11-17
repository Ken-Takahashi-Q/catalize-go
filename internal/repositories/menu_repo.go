package repositories

import (
	"catalize-go/internal/models"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type MenuRepo struct {
	collection *mongo.Collection
}

func (r *MenuRepo) CreateOrder(ctx context.Context, order *models.Order) error {
	order.CreatedAt = time.Now()
	_, err := r.collection.InsertOne(ctx, order)
	return err
}

func (r *MenuRepo) GetOrder(ctx context.Context, id string) (*models.Order, error) {
	var order models.Order
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&order)
	return &order, err
}
