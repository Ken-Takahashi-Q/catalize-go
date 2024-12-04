package handlers

import (
	"catalize-go/internal/services"
	"catalize-go/internal/utils"
	"net/http"

	"go.mongodb.org/mongo-driver/mongo"
)

func GetMenu(w http.ResponseWriter, r *http.Request) {
	menu, err := services.GetMenuService()
	if err != nil {
		if err == mongo.ErrNoDocuments {
			http.Error(w, "No menu found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to fetch menu", http.StatusInternalServerError)
		}
		return
	}

	utils.JSONResponse(w, "success", menu, http.StatusOK)
}

func GetMenuCategory(w http.ResponseWriter, r *http.Request) {
	menu, err := services.GetMenuCategoryService()
	if err != nil {
		if err == mongo.ErrNoDocuments {
			http.Error(w, "No menu category found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to fetch menu", http.StatusInternalServerError)
		}
		return
	}

	utils.JSONResponse(w, "success", menu, http.StatusOK)
}
