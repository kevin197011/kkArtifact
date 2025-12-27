// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package server

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/kk/kkartifact-server/internal/api"
	"github.com/kk/kkartifact-server/internal/auth"
	"github.com/kk/kkartifact-server/internal/cache"
	"github.com/kk/kkartifact-server/internal/config"
	"github.com/kk/kkartifact-server/internal/database"
	"github.com/kk/kkartifact-server/internal/bootstrap"
	"github.com/kk/kkartifact-server/internal/middleware"
	"github.com/kk/kkartifact-server/internal/storage"
	
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	
	_ "github.com/kk/kkartifact-server/docs" // swagger docs
)

// Server represents the HTTP server
type Server struct {
	config  *config.Config
	router  *gin.Engine
	db      *database.DB
	storage storage.Storage
	cache   cache.Cache
}

// New creates a new server instance
func New(cfg *config.Config) (*Server, error) {
	// Run database migrations (optional, can be disabled via env var)
	if os.Getenv("SKIP_MIGRATIONS") != "true" {
		migrationsPath := os.Getenv("MIGRATIONS_PATH")
		if migrationsPath == "" {
			migrationsPath = "./migrations/migrations"
		}
		if err := database.RunMigrations(&cfg.Database, migrationsPath); err != nil {
			log.Printf("Warning: failed to run migrations: %v", err)
			// Continue anyway - migrations might have been run manually
		}
	}

	// Initialize database
	db, err := database.New(&cfg.Database)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	// Initialize storage
	storageBackend, err := storage.NewStorage(&cfg.Storage)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize storage: %w", err)
	}

	// Initialize cache
	cacheBackend, err := cache.NewRedisCache(
		cfg.Redis.Host,
		cfg.Redis.Port,
		cfg.Redis.Password,
		cfg.Redis.DB,
	)
	if err != nil {
		log.Printf("Warning: failed to initialize Redis cache: %v", err)
		// Continue without cache in development
		cacheBackend = nil
	}

	// Setup router
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())
	router.Use(middleware.CORS())
	router.Use(middleware.Gzip())
	
	// Swagger documentation route (optional, can be disabled in production)
	if os.Getenv("ENABLE_SWAGGER") != "false" {
		router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	// Initialize admin user if enabled
	if os.Getenv("SKIP_ADMIN_USER") != "true" {
		adminUsername, err := bootstrap.EnsureAdminUser(db)
		if err != nil {
			log.Printf("Warning: failed to initialize admin user: %v", err)
		} else if adminUsername != "" {
			adminPassword := os.Getenv("ADMIN_PASSWORD")
			if adminPassword == "" {
				adminPassword = "admin"
			}
			log.Printf("Admin user available: %s", adminUsername)
		}
	}

	// Initialize admin token only if ADMIN_TOKEN is set
	// If ADMIN_TOKEN is not set, skip creation
	adminToken, err := bootstrap.EnsureAdminToken(db)
	if err != nil {
		log.Printf("Warning: failed to initialize admin token: %v", err)
	} else if adminToken != "" {
		log.Printf("========================================")
		log.Printf("Admin Token Created/Found:")
		adminTokenName := os.Getenv("ADMIN_TOKEN_NAME")
		if adminTokenName == "" {
			adminTokenName = "admin-initial-token"
		}
		log.Printf("  Token: %s", adminToken)
		log.Printf("  Name: %s", adminTokenName)
		log.Printf("  Permissions: pull, push, promote, admin")
		log.Printf("========================================")
	}

	// Initialize handler
	authenticator := auth.NewTokenAuthenticator(db)
	handler := api.NewHandler(db, storageBackend, authenticator)
	handler.RegisterRoutes(router)

	server := &Server{
		config:  cfg,
		router:  router,
		db:      db,
		storage: storageBackend,
		cache:   cacheBackend,
	}

	// TODO: Initialize scheduler and cleanup task
	// scheduler := scheduler.New()
	// artifactManager := storage.NewArtifactManager(storageBackend)
	// cleanupTask := scheduler.NewCleanupTask(db, artifactManager)
	// scheduler.AddTask(cleanupTask)
	// go scheduler.Start(context.Background())

	return server, nil
}

// Start starts the HTTP server
func (s *Server) Start() error {
	addr := fmt.Sprintf("%s:%s", s.config.Server.Host, s.config.Server.Port)
	log.Printf("Starting server on %s", addr)
	return s.router.Run(addr)
}

// Close closes the server and cleans up resources
func (s *Server) Close() error {
	if s.db != nil {
		return s.db.Close()
	}
	return nil
}

