// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package config

import (
	"os"
	"strconv"
)

// Config holds all configuration for the server
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Redis    RedisConfig
	Storage  StorageConfig
	Log      LogConfig
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Host string
	Port string
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Host         string
	Port         string
	Name         string
	User         string
	Password     string
	SSLMode      string
	MaxOpenConns int
	MaxIdleConns int
}

// RedisConfig holds Redis configuration
type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

// StorageConfig holds storage configuration
type StorageConfig struct {
	Type       string
	S3Endpoint string
	S3AccessKey string
	S3SecretKey string
	S3Bucket   string
	S3Region   string
	S3UseSSL   bool
	LocalPath  string
	BasePath   string
}

// LogConfig holds logging configuration
type LogConfig struct {
	Level  string
	Format string
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	return &Config{
		Server: ServerConfig{
			Host: getEnv("SERVER_HOST", "0.0.0.0"),
			Port: getEnv("SERVER_PORT", "8080"),
		},
		Database: DatabaseConfig{
			Host:         getEnv("DB_HOST", "localhost"),
			Port:         getEnv("DB_PORT", "5432"),
			Name:         getEnv("DB_NAME", "kkartifact"),
			User:         getEnv("DB_USER", "kkartifact"),
			Password:     getEnv("DB_PASSWORD", ""),
			SSLMode:      getEnv("DB_SSLMODE", "disable"),
			MaxOpenConns: getEnvAsInt("DB_MAX_OPEN_CONNS", 50), // Increased to 50 for high concurrency uploads (PostgreSQL max_connections=200)
			MaxIdleConns: getEnvAsInt("DB_MAX_IDLE_CONNS", 10),  // Increased to 10 for better connection reuse
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnv("REDIS_PORT", "6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnvAsInt("REDIS_DB", 0),
		},
		Storage: StorageConfig{
			Type:        getEnv("STORAGE_TYPE", "s3"),
			S3Endpoint:  getEnv("STORAGE_S3_ENDPOINT", ""),
			S3AccessKey: getEnv("STORAGE_S3_ACCESS_KEY", ""),
			S3SecretKey: getEnv("STORAGE_S3_SECRET_KEY", ""),
			S3Bucket:    getEnv("STORAGE_S3_BUCKET", "kkartifact"),
			S3Region:    getEnv("STORAGE_S3_REGION", "us-east-1"),
			S3UseSSL:    getEnvAsBool("STORAGE_S3_USE_SSL", false),
			LocalPath:   getEnv("STORAGE_LOCAL_PATH", "/repos"),
			BasePath:    getEnv("STORAGE_BASE_PATH", "/repos"),
		},
		Log: LogConfig{
			Level:  getEnv("LOG_LEVEL", "info"),
			Format: getEnv("LOG_FORMAT", "json"),
		},
	}, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	intValue, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	return intValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value == "true" || value == "1"
}
