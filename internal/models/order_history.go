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
	TotalPrice  float64     `bson:"TotalPrice" json:"total_price"`
	OrderStatus OrderStatus `bson:"OrderStatus" json:"order_status"`
}

type OrderItem struct {
	MenuID     int     `bson:"MenuID" json:"menu_id"`
	MenuNameTH string  `bson:"MenuNameTH" json:"menu_name_th"`
	MenuNameEN string  `bson:"MenuNameEN" json:"menu_name_en"`
	Price      float64 `bson:"Price" json:"price"`
	Category   int     `bson:"Category" json:"category"`
	Status     int     `bson:"Status" json:"status"`
	Image      string  `bson:"Image" json:"image"`
	Quantity   float64 `bson:"Quantity" json:"quantity"`
}
