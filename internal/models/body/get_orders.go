package body

import (
	"catalize-go/internal/models"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GetOrders struct {
	ID         primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	UserID     int                `bson:"UserID" json:"user_id"`
	TableID    int                `bson:"TableID" json:"table_id"`
	TableVisit int                `bson:"TableVisit" json:"table_visit"`

	OrderID     int                `bson:"OrderID" json:"order_id"`
	OrderDate   time.Time          `bson:"OrderDate" json:"order_date"`
	CreatedAt   time.Time          `bson:"CreatedAt" json:"created_at"`
	OrderStatus models.OrderStatus `bson:"OrderStatus" json:"order_status"`
}
