package services

import (
	"catalize-go/internal/db"
	"catalize-go/internal/models"
	"catalize-go/internal/repositories"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type MenuService struct {
	menuRepo *repositories.MenuRepo
}

func NewMenuService(repo *repositories.MenuRepo) *MenuService {
	return &MenuService{
		menuRepo: repo,
	}
}

func GetMenuService() (*models.Menu, error) {
	collection :=
		db.GetCollection(models.GetCollection{DBName: "menu", Collection: "menu"})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var menu models.Menu
	err := collection.FindOne(ctx, bson.M{}).Decode(&menu)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, mongo.ErrNoDocuments
		}
		return nil, err
	}

	return &menu, nil
}
