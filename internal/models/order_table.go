package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OrderTable struct {
	ID         primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Date       time.Time          `bson:"Date" json:"date"`
	TableID    int                `bson:"TableID" json:"table_id"`
	TableVisit int                `bson:"TableVisit" json:"table_visit"`
}
