package main

import (
	"context"
	"curd-web-go/config"
	"curd-web-go/routes"
	"curd-web-go/tracing"
	"log"
	"net/http"

	"github.com/joho/godotenv"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func main() {

	// Load .env
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, relying on container env variables")
	}

	// Init Tracer
	ctx := context.Background()
	shutdown := tracing.InitTracer(ctx, "crud-web-go", "otel-collector:4317")
	defer shutdown(ctx)

	// Connect Database
	config.DatabaseConnection(ctx)
	if config.DB == nil {
		log.Fatal("Database connection failed!")
	}

	if _, err := config.InitRedis(); err != nil {
		log.Fatalf("Redis initialization failed: %v", err)
	}

	handler := otelhttp.NewHandler(routes.Routes(), "http-server")

	log.Println(" Server running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}
