// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package bootstrap

import (
	"fmt"
	"log"
	"os"

	"github.com/kk/kkartifact-server/internal/auth"
	"github.com/kk/kkartifact-server/internal/database"
)

// EnsureAdminToken ensures an admin token exists in the database
// Only creates token if ADMIN_TOKEN is set in environment
// If ADMIN_TOKEN is not set, skips creation
// Returns the token string (plain text) for display
func EnsureAdminToken(db *database.DB) (string, error) {
	// Use os.Getenv directly to check if ADMIN_TOKEN is set (not using getEnv with default)
	adminTokenEnv := os.Getenv("ADMIN_TOKEN")
	adminTokenName := getEnv("ADMIN_TOKEN_NAME", "admin-initial-token")

	// If ADMIN_TOKEN is not set, skip creation
	if adminTokenEnv == "" {
		log.Printf("Skipping admin token creation (ADMIN_TOKEN not set)")
		return "", nil
	}

	// Check if admin token already exists
	tokenRepo := database.NewTokenRepository(db)
	
	// Try to find existing admin token by name
	tokens, err := tokenRepo.List()
	if err == nil {
		for _, token := range tokens {
			if token.Name.Valid && token.Name.String == adminTokenName {
				// Admin token exists
				log.Printf("Admin token already exists with name: %s (using provided ADMIN_TOKEN)", adminTokenName)
				return adminTokenEnv, nil
			}
		}
	}

	// Use provided token
	tokenPlain := adminTokenEnv
	log.Printf("Creating admin token using provided ADMIN_TOKEN from environment")

	// Hash the token
	tokenHash, err := auth.HashToken(tokenPlain)
	if err != nil {
		return "", fmt.Errorf("failed to hash admin token: %w", err)
	}

	// Create admin token with all permissions
	permissions := []string{"pull", "push", "promote", "admin"}
	_, err = tokenRepo.Create(
		tokenHash,
		adminTokenName,
		nil, // Global scope (no project_id)
		nil, // Global scope (no app_id)
		permissions,
		nil, // No expiration
	)

	if err != nil {
		return "", fmt.Errorf("failed to create admin token: %w", err)
	}

	return tokenPlain, nil
}

// getEnv gets environment variable with default value
func getEnv(key, defaultValue string) string {
	value := defaultValue
	if v := getEnvFunc(key); v != "" {
		value = v
	}
	return value
}

// getEnvFunc is a function variable for getting environment variables
// This allows testing by overriding the function
var getEnvFunc = func(key string) string {
	// Default implementation uses os.Getenv
	return osGetenv(key)
}

// osGetenv is a wrapper that can be mocked in tests
var osGetenv = os.Getenv

