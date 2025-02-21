package table

import (
	"catalize-go/internal/models"
	"time"
)

type TableCheckBill struct {
	Date       time.Time      `bson:"Date" json:"date"`
	TableID    int            `bson:"TableID" json:"table_id"`
	TableVisit int            `bson:"TableVisit" json:"table_visit"`
	Orders     []models.Order `bson:"Orders" json:"orders"`
	TotalPrice float64        `bson:"TotalPrice" json:"total_price"`
}
