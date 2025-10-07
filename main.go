package main

import (
	"context"
	"fmt"
	"log"
	"os"

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

	router := mux.NewRouter()
	SetupRoutes(router)
	fmt.Println("Server running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
