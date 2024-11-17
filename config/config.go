package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

var MongoURI string
var Port string
var Host string

func Load() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	Port = os.Getenv("PORT")
	if Port == "" {
		Port = "8080"
	}

	Host = os.Getenv("HOST")
	if Host == "" {
		Host = "localhost"
	}

	MongoURI = os.Getenv("MONGO_URI")
	if MongoURI == "" {
		MongoURI = "mongodb://localhost:27017"
	}
}
