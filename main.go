package main

import (
	"curd-web-go/config"
	"curd-web-go/routes"
	"log"
	"net/http"
)

func main() {
	config.DatabaseConnection()
	if config.DB == nil {
		log.Fatal("Database connection failed!")
	}

	router := routes.Routes()

	log.Println(" Server running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
