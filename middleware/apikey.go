package middleware

import (
	"context"
	"naturedopesApi/endpoints"
	"net/http"

	"github.com/jackc/pgx/v4"
)

// ApiKeyMiddleware validates the X-API-Key header
func ApiKeyMiddleware(connectDB func() (*pgx.Conn, error)) func(http.Handler) http.Handler {

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// 1. Extract API key from header
			apiKey := r.Header.Get("X-API-Key")

			// 2. If no key, reject immediately
			if apiKey == "" {
				http.Error(w, "Missing X-API-Key header", http.StatusUnauthorized)
				return
			}

			// 3. Connect to DB and validate
			conn, err := connectDB()
			if err != nil {
				http.Error(w, "Database connection error", http.StatusInternalServerError)
				return
			}
			defer conn.Close(context.Background())

			valid, err := endpoints.ValidateApiKey(conn, apiKey)
			if err != nil {
				http.Error(w, "Error validating API key", http.StatusInternalServerError)
				return
			}

			if !valid {
				http.Error(w, "Invalid or revoked API key", http.StatusUnauthorized)
				return
			}

			// 5. If valid, continue to your handler
			next.ServeHTTP(w, r)
		})
	}
}
