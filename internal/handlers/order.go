package handlers

import (
	"catalize-go/internal/models"
	"catalize-go/internal/models/body"
	"catalize-go/internal/services"
	"catalize-go/internal/utils"
	"encoding/json"
	"net/http"
)

type OrderHandler struct {
	orderService *services.OrderService
}

func NewOrderHandler(svc *services.OrderService) *OrderHandler {
	return &OrderHandler{orderService: svc}
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

func KitchenDoneOrder(w http.ResponseWriter, r *http.Request) {
	var kitchenDoneOrder body.GetOrders
	if err := json.NewDecoder(r.Body).Decode(&kitchenDoneOrder); err != nil {
		utils.JSONResponse(w, "fail", "Invalid input", http.StatusBadRequest)
		return
	}

	orders, err := services.KitchenDoneOrderService(kitchenDoneOrder)
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
