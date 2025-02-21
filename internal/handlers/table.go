package handlers

import (
	"catalize-go/internal/db"
	"catalize-go/internal/models"
	"catalize-go/internal/models/table"
	"catalize-go/internal/services"
	"catalize-go/internal/utils"
	"context"
	"encoding/json"
	"net/http"
	"time"
)

type TableHandler struct {
	tableService *services.TableService
}

func NewTableHandler(svc *services.TableService) *TableHandler {
	return &TableHandler{tableService: svc}
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

func CheckBill(w http.ResponseWriter, r *http.Request) {
	var orderTable models.OrderTable
	if err := json.NewDecoder(r.Body).Decode(&orderTable); err != nil {
		utils.JSONResponse(w, "fail", "Invalid input", http.StatusBadRequest)
		return
	}

	if orderTable.Date.IsZero() {
		orderTable.Date = time.Now()
	}

	tableVisit, err := services.CheckBillService(table.TableCheckBill{Date: orderTable.Date, TableID: orderTable.TableID})
	if err != nil {
		utils.JSONResponse(w, "fail", "Failed to get table visit", http.StatusInternalServerError)
		return
	}

	orderTable.TableVisit = tableVisit

	utils.JSONResponse(w, "success", orderTable, http.StatusCreated)
}
