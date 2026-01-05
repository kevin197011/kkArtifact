// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package bootstrap

import (
	"log"
	"github.com/kk/kkartifact-server/internal/database"
)

// EnsureAdminUser ensures an admin user exists in the database
// Returns the username
func EnsureAdminUser(db *database.DB) (string, error) {
	adminUsername := getEnv("ADMIN_USERNAME", "admin")
	adminPassword := getEnv("ADMIN_PASSWORD", "admin")
	skipAdminUser := getEnv("SKIP_ADMIN_USER", "false")

	if skipAdminUser == "true" {
		log.Printf("Skipping admin user creation (SKIP_ADMIN_USER=true)")
		return "", nil
	}

	userRepo := database.NewUserRepository(db)

	// Check if admin user already exists
	_, err := userRepo.GetByUsername(adminUsername)
	if err == nil {
		// User already exists
		log.Printf("Admin user already exists: %s", adminUsername)
		return adminUsername, nil
	}

	// Create admin user
	passwordHash, err := database.HashPassword(adminPassword)
	if err != nil {
		return "", err
	}

	_, err = userRepo.Create(adminUsername, passwordHash)
	if err != nil {
		return "", err
	}

	log.Printf("========================================")
	log.Printf("Admin User Created:")
	log.Printf("  Username: %s", adminUsername)
	log.Printf("  Password: %s", adminPassword)
	log.Printf("========================================")

	return adminUsername, nil
}

