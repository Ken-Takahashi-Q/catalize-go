package models

import (
	"time"
)

type Table struct {
	TableID    	int			`bson:"table_id" json:"table_id"`
	TableVisit 	int			`bson:"table_visit" json:"table_visit"`
	Status  	string		`bson:"status" json:"status"`
	CreatedAt  	time.Time	`bson:"created_at" json:"created_at"`
}
