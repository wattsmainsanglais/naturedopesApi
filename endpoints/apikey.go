package endpoints

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/jackc/pgx/v4"
)

type ApiKey struct {
	Id         int       `json:"id"`
	Key        string    `json:"key"`
	Name       string    `json:"name"`
	Created_at time.Time `json:"created_at"`
	Last_used  *time.Time `json:"last_used"`
	Revoked    bool      `json:"revoked"`
}

// GenerateApiKey creates a new API key with a random 64-character hex string
func GenerateApiKey(conn *pgx.Conn, name string) (*ApiKey, error) {
	// Generate 32 random bytes (64 hex characters)
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return nil, fmt.Errorf("error generating random key: %v", err)
	}
	key := hex.EncodeToString(bytes)

	var apiKey ApiKey
	err := conn.QueryRow(
		context.Background(),
		"INSERT INTO api_keys (key, name, created_at, revoked) VALUES ($1, $2, NOW(), false) RETURNING id, key, name, created_at, revoked",
		key, name,
	).Scan(&apiKey.Id, &apiKey.Key, &apiKey.Name, &apiKey.Created_at, &apiKey.Revoked)

	if err != nil {
		return nil, fmt.Errorf("error creating API key: %v", err)
	}

	return &apiKey, nil
}

// ValidateApiKey checks if an API key is valid (exists, not revoked) and updates last_used
func ValidateApiKey(conn *pgx.Conn, key string) (bool, error) {
	var revoked bool
	err := conn.QueryRow(
		context.Background(),
		"SELECT revoked FROM api_keys WHERE key = $1",
		key,
	).Scan(&revoked)

	if err == pgx.ErrNoRows {
		return false, nil // Key doesn't exist
	}
	if err != nil {
		return false, fmt.Errorf("error validating API key: %v", err)
	}

	if revoked {
		return false, nil // Key is revoked
	}

	// Update last_used timestamp
	_, err = conn.Exec(
		context.Background(),
		"UPDATE api_keys SET last_used = NOW() WHERE key = $1",
		key,
	)
	if err != nil {
		// Log error but don't fail validation
		fmt.Printf("Warning: failed to update last_used for key: %v\n", err)
	}

	return true, nil
}

// GetApiKeys retrieves all API keys
func GetApiKeys(conn *pgx.Conn) ([]ApiKey, error) {
	rows, err := conn.Query(
		context.Background(),
		"SELECT id, key, name, created_at, last_used, revoked FROM api_keys ORDER BY created_at DESC",
	)
	if err != nil {
		return nil, fmt.Errorf("error getting API keys: %v", err)
	}
	defer rows.Close()

	var apiKeys []ApiKey
	for rows.Next() {
		var apiKey ApiKey
		err := rows.Scan(&apiKey.Id, &apiKey.Key, &apiKey.Name, &apiKey.Created_at, &apiKey.Last_used, &apiKey.Revoked)
		if err != nil {
			return nil, fmt.Errorf("error scanning API key: %v", err)
		}
		apiKeys = append(apiKeys, apiKey)
	}

	return apiKeys, nil
}

// RevokeApiKey marks an API key as revoked
func RevokeApiKey(conn *pgx.Conn, id int) error {
	result, err := conn.Exec(
		context.Background(),
		"UPDATE api_keys SET revoked = true WHERE id = $1",
		id,
	)
	if err != nil {
		return fmt.Errorf("error revoking API key: %v", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("API key not found")
	}

	return nil
}
