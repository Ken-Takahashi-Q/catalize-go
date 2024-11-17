package handlers

import (
	"catalize-go/internal/db"
	"catalize-go/internal/models"
	"catalize-go/internal/models/body"
	"catalize-go/internal/services"
	"catalize-go/internal/utils"
	"context"
	"encoding/json"
	"net/http"
	"time"
)

type OrderHandler struct {
	orderService *services.OrderService
}

func NewOrderHandler(svc *services.OrderService) *OrderHandler {
	return &OrderHandler{orderService: svc}
}

func CreateTableVisit(w http.ResponseWriter, r *http.Request) {
	var orderTable models.OrderTable
	if err := json.NewDecoder(r.Body).Decode(&orderTable); err != nil {
		utils.JSONResponse(w, "fail", "Invalid input", http.StatusBadRequest)
		return
	}

	if orderTable.Date.IsZero() {
		orderTable.Date = time.Now()
	}

	tableVisit, err := services.CreateTableVisitService(models.OrderTable{Date: orderTable.Date, TableID: orderTable.TableID})
	if err != nil {
		utils.JSONResponse(w, "fail", "Failed to generate table visit", http.StatusInternalServerError)
		return
	}

	collection := db.GetCollection(models.GetCollection{DBName: "order", Collection: "history"})
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = collection.InsertOne(ctx, orderTable)
	if err != nil {
		utils.JSONResponse(w, "fail", "Failed to serve table ID", http.StatusInternalServerError)
		return
	}

	orderTable.TableVisit = tableVisit

	utils.JSONResponse(w, "success", orderTable, http.StatusCreated)
}

func GetTableVisit(w http.ResponseWriter, r *http.Request) {
	var orderTable models.OrderTable
	if err := json.NewDecoder(r.Body).Decode(&orderTable); err != nil {
		utils.JSONResponse(w, "fail", "Invalid input", http.StatusBadRequest)
		return
	}

	if orderTable.Date.IsZero() {
		orderTable.Date = time.Now()
	}

	tableVisit, err := services.GetTableVisitService(models.OrderTable{Date: orderTable.Date, TableID: orderTable.TableID})
	if err != nil {
		utils.JSONResponse(w, "fail", "Failed to get table visit", http.StatusInternalServerError)
		return
	}

	orderTable.TableVisit = tableVisit

	utils.JSONResponse(w, "success", orderTable, http.StatusCreated)
}

func CreateOrder(w http.ResponseWriter, r *http.Request) {
	var order models.Order
	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		utils.JSONResponse(w, "fail", "Invalid input", http.StatusBadRequest)
		return
	}

	if err := services.CreateOrderService(&order); err != nil {
		utils.JSONResponse(w, "fail", err.Error(), http.StatusInternalServerError)
		return
	}

	utils.JSONResponse(w, "success", order, http.StatusCreated)
}

func GetOrders(w http.ResponseWriter, r *http.Request) {
	var getOrders body.GetOrders
	if err := json.NewDecoder(r.Body).Decode(&getOrders); err != nil {
		utils.JSONResponse(w, "fail", "Invalid input", http.StatusBadRequest)
		return
	}

	orders, err := services.GetOrdersService(getOrders)
	if err != nil {
		utils.JSONResponse(w, "fail", "Failed to get orders", http.StatusInternalServerError)
		return
	}

	utils.JSONResponse(w, "success", orders, http.StatusCreated)
}

func ClearAllOrder(w http.ResponseWriter, r *http.Request) {
	err := services.ClearAllOrderService()
	if err != nil {
		utils.JSONResponse(w, "fail", err.Error(), http.StatusInternalServerError)
		return
	}

	utils.JSONResponse(w, "success", "", http.StatusCreated)
}
