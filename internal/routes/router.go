package routes

import (
	"catalize-go/internal/handlers"

	"github.com/gorilla/mux"
)

func SetupRouter() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/menu", handlers.GetMenu).Methods("GET")
	router.HandleFunc("/order/create_table_visit", handlers.CreateTableVisit).Methods("GET")
	router.HandleFunc("/order/get_table_visit", handlers.GetTableVisit).Methods("GET")

	router.HandleFunc("/order/create_order", handlers.CreateOrder).Methods("POST")
	router.HandleFunc("/order/get_orders", handlers.GetOrders).Methods("GET")

	router.HandleFunc("/order/clear_all_order", handlers.ClearAllOrder).Methods("DELETE")
	return router
}
