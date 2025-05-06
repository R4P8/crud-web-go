package main

import (
	"curd-web-go/config"
	"curd-web-go/routes"
	"log"
	"net/http"

	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, relying on container env variables")
	}

	config.DatabaseConnection()
	if config.DB == nil {
		log.Fatal("Database connection failed!")
	}

	router := routes.Routes()

	log.Println(" Server running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
