package utils

import (
	"catalize-go/internal/models"
	"encoding/json"
	"net/http"
)

func JSONResponse(w http.ResponseWriter, code string, data interface{}, statusCode int) {
	response := models.Response{
		Code: code,
		Data: data,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}
