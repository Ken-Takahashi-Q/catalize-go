package main

import (
	"catalize-go/config"
	"catalize-go/internal/db"
	"catalize-go/internal/routes"
	"log"
	"net/http"
	"os"
)

func main() {
	config.Load()

	err := db.Connect(config.MongoURI)
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}

	router := routes.SetupRouter()

	env := os.Getenv("ENV")
	var host string
	if env == "local" {
		host = "localhost"
	} else {
		host = ""
	}

	log.Println("Server is running on port", config.Port)
	log.Fatal(http.ListenAndServe(host+":"+config.Port, router))
}
