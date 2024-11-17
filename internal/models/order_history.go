package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Order struct {
	ID         primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	UserID     int                `bson:"UserID" json:"user_id"`
	TableID    int                `bson:"TableID" json:"table_id"`
	TableVisit int                `bson:"TableVisit" json:"table_visit"`
	OrderCount int                `bson:"OrderCount" json:"order_count"`

	OrderDate   time.Time   `bson:"OrderDate" json:"order_date"`
	CreatedAt   time.Time   `bson:"CreatedAt" json:"created_at"`
	Items       []OrderItem `bson:"Items" json:"items"`
	OrderStatus OrderStatus `bson:"OrderStatus" json:"order_status"`
}

type OrderItem struct {
	MenuID   int     `bson:"MenuID" json:"menu_id"`
	Quantity float64 `bson:"Quantity" json:"quantity"`
}
