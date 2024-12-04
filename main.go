package main

import (
	"catalize-go/config"
	"catalize-go/internal/db"
	"catalize-go/internal/routes"
	"log"
	"net/http"
	"os"

	"github.com/rs/cors"
)

func main() {
	config.Load()

	err := db.Connect(config.MongoURI)
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}

	router := routes.SetupRouter()

	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	}).Handler(router)

	env := os.Getenv("ENV")
	var host string
	if env == "local" {
		host = "localhost"
	} else {
		host = ""
	}

	log.Println("Server is running on port", config.Port)
	log.Fatal(http.ListenAndServe(host+":"+config.Port, corsHandler))
}
