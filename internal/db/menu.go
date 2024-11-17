package db

import (
	"catalize-go/internal/models"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// GetMenu retrieves the menu from the database
func GetMenu(ctx context.Context, collection *mongo.Collection) ([]models.Menu, error) {
	var menuItems []models.Menu
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var menuItem models.Menu
		if err := cursor.Decode(&menuItem); err != nil {
			return nil, err
		}
		menuItems = append(menuItems, menuItem)
	}

	return menuItems, cursor.Err()
}
