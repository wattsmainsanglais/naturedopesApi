package main

import (
	"context"
	"fmt"
	"log"
	"naturedopesApi/middleware"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4"
	"net/http"
	//"github.com/joho/godotenv"
)

func connectToDB() (*pgx.Conn, error) {
	/*err := godotenv.Load()   Commenting this line for railway deployment
	if err != nil {
		log.Fatal(err)
	}*/

	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		return nil, fmt.Errorf("DATABASE_URL environment variable is not set")
	}
	conn, err := pgx.Connect(context.Background(), connStr)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func SetupRoutes(router *mux.Router) {
	// Protected image endpoints (require API key)
	apiKeyAuth := middleware.ApiKeyMiddleware(connectToDB)

	imageRouter := router.PathPrefix("/images").Subrouter()
	imageRouter.Use(apiKeyAuth)
	imageRouter.HandleFunc("", getImagesHandler).Methods("GET")
	imageRouter.HandleFunc("/{id}", getImageHandler).Methods("GET")

	// API key management endpoints (unprotected for now)
	router.HandleFunc("/api/keys", createApiKeyHandler).Methods("POST")
	router.HandleFunc("/api/keys", getApiKeysHandler).Methods("GET")
	router.HandleFunc("/api/keys/{id}", revokeApiKeyHandler).Methods("DELETE")
}

func main() {
	// Initialize rate limiter
	rateLimiter := middleware.NewRateLimiter()
	defer rateLimiter.Stop()

	router := mux.NewRouter()

	// Apply rate limiting to all routes
	router.Use(rateLimiter.RateLimitMiddleware)

	SetupRoutes(router)

	// Configure CORS for open access
	corsOptions := handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}), // Allow all origins for open access
		handlers.AllowedMethods([]string{"GET", "POST", "DELETE", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Content-Type", "X-API-Key"}), // Headers clients can send
	)

	fmt.Println("Server running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", corsOptions(router)))
}
