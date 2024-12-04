package services

import (
	"catalize-go/internal/db"
	"catalize-go/internal/models"
	"catalize-go/internal/repositories"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

type MenuService struct {
	menuRepo *repositories.MenuRepo
}

func NewMenuService(repo *repositories.MenuRepo) *MenuService {
	return &MenuService{
		menuRepo: repo,
	}
}

func GetMenuService() (*[]models.MenuItem, error) {
	collection := db.GetCollection(models.GetCollection{
		DBName:     "menu",
		Collection: "menu",
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var menu models.Menu
	err := collection.FindOne(ctx, bson.M{}).Decode(&menu)
	if err != nil {
		return nil, err
	}

	if menu.MenuString != "" {
		var decodedMenu struct {
			Menus []models.MenuItem `json:"menus"`
		}

		err := json.Unmarshal([]byte(menu.MenuString), &decodedMenu)
		if err != nil {
			return nil, fmt.Errorf("failed to decode MenuString: %v", err)
		}
		menu.Menu = decodedMenu.Menus
	}

	return &menu.Menu, nil
}

func GetMenuCategoryService() ([]models.MenuCategory, error) {
	collection := db.GetCollection(models.GetCollection{
		DBName:     "menu",
		Collection: "category",
	})

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	var result struct {
		Category string `bson:"Category"`
	}

	err := collection.FindOne(ctx, bson.M{}).Decode(&result)
	if err != nil {
		return nil, err
	}

	var menuCategories []models.MenuCategory
	if err := json.Unmarshal([]byte(result.Category), &menuCategories); err != nil {
		return nil, err
	}

	return menuCategories, nil
}
