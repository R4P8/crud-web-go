package main

import (
	"context"
	"curd-web-go/config"
	"curd-web-go/routes"
	"curd-web-go/tracing"
	"log"
	"net/http"

	"github.com/joho/godotenv"
)

func main() {

	// Load .env
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, relying on container env variables")
	}

	// Init Tracer
	ctx := context.Background()
	shutdown := tracing.InitTracer(ctx, "crud-web-go", "jaeger.rizqifathirafa.my.id:4317")
	defer shutdown(ctx)

	// Connect Database
	config.DatabaseConnection(ctx)
	if config.DB == nil {
		log.Fatal("Database connection failed!")
	}

	router := routes.Routes()

	log.Println(" Server running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
