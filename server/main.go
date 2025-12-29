// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

// @title           kkArtifact API
// @version         1.0
// @description     kkArtifact is a modern artifact management and synchronization system designed to replace rsync + SSH deployment workflows.
// @description     The system provides multi-project and multi-app artifact storage with hash-based versioning, HTTP-based API for artifact upload, download, and management.
// @description     Token-based authentication with fine-grained permissions (Global/Project/App scopes).

// @contact.name   kk
// @contact.url    https://opensource.org/licenses/MIT
// @license.name   MIT
// @license.url    https://opensource.org/licenses/MIT

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token. Example: "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."

package main

import (
	"log"
	"os"

	"github.com/kk/kkartifact-server/internal/config"
	"github.com/kk/kkartifact-server/internal/server"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	srv, err := server.New(cfg)
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
		os.Exit(1)
	}

	if err := srv.Start(); err != nil {
		log.Fatalf("Server failed to start: %v", err)
		os.Exit(1)
	}
}

