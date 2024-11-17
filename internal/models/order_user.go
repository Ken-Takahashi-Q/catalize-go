package models

import (
	"time"
)

type UserOrder struct {
	UserID     int       `bson:"UserID" json:"user_id"`
	LastOrder  time.Time `bson:"LastOrder" json:"last_order"`
	TableID    int       `bson:"TableID" json:"table_id"`
	TableVisit int       `bson:"TableVisit" json:"table_visit"`

	OrderCount  int         `bson:"OrderCount" json:"order_count"`
	OrderStatus OrderStatus `bson:"OrderStatus" json:"order_status"`
}
