package handlers

import (
	"catalize-go/internal/db"
	"catalize-go/internal/models"
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
	var table models.Table
	if err := json.NewDecoder(r.Body).Decode(&table); err != nil {
		utils.JSONResponse(w, "fail", "Invalid input", http.StatusBadRequest)
		return
	}

	if table.CreatedAt.IsZero() {
		table.CreatedAt = time.Now()
	}

	tableVisit, err := services.CreateTableVisitService(models.Table{CreatedAt: table.CreatedAt, TableID: table.TableID})
	if err != nil {
		utils.JSONResponse(w, "fail", "Failed to generate table visit", http.StatusInternalServerError)
		return
	}

	collection := db.GetCollection(models.GetCollection{DBName: "order", Collection: "history"})
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = collection.InsertOne(ctx, table)
	if err != nil {
		utils.JSONResponse(w, "fail", "Failed to serve table ID", http.StatusInternalServerError)
		return
	}

	table.TableVisit = tableVisit

	utils.JSONResponse(w, "success", table, http.StatusCreated)
}

func GetTableVisit(w http.ResponseWriter, r *http.Request) {
	var table models.Table
	if err := json.NewDecoder(r.Body).Decode(&table); err != nil {
		utils.JSONResponse(w, "fail", "Invalid input", http.StatusBadRequest)
		return
	}

	if table.CreatedAt.IsZero() {
		table.CreatedAt = time.Now()
	}

	tableVisit, err := services.GetTableVisitService(models.Table{CreatedAt: table.CreatedAt, TableID: table.TableID})
	if err != nil {
		utils.JSONResponse(w, "fail", "Failed to get table visit", http.StatusInternalServerError)
		return
	}

	table.TableVisit = tableVisit

	utils.JSONResponse(w, "success", table, http.StatusCreated)
}

func CheckBill(w http.ResponseWriter, r *http.Request) {
	var table models.Table
	if err := json.NewDecoder(r.Body).Decode(&table); err != nil {
		utils.JSONResponse(w, "fail", "Invalid input", http.StatusBadRequest)
		return
	}

	if table.CreatedAt.IsZero() {
		table.CreatedAt = time.Now()
	}

	_, err := services.CheckBillService(table)
	if err != nil {
		utils.JSONResponse(w, "fail", "Failed to get table visit", http.StatusInternalServerError)
		return
	}

	// table.TableVisit = tableVisit

	utils.JSONResponse(w, "success", table, http.StatusCreated)
}

func GetTables(w http.ResponseWriter, r *http.Request) {
	tables, err := services.GetTablesService()
	if err != nil {
		utils.JSONResponse(w, "fail", "Failed to get tables", http.StatusInternalServerError)
		return
	}

	utils.JSONResponse(w, "success", tables, http.StatusCreated)
}

func UpdateTable(w http.ResponseWriter, r *http.Request) {
	var updateTable models.Table
	if err := json.NewDecoder(r.Body).Decode(&updateTable); err != nil {
		utils.JSONResponse(w, "fail", "Invalid input", http.StatusBadRequest)
		return
	}

	updatedTable, err := services.UpdateTableService(updateTable)
	if err != nil {
		utils.JSONResponse(w, "fail", "Failed to update table", http.StatusInternalServerError)
		return
	}

	utils.JSONResponse(w, "success", updatedTable, http.StatusCreated)
}