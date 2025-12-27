// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

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

