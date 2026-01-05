// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package storage

import (
	"fmt"

	"github.com/kk/kkartifact-server/internal/config"
)

// NewStorage creates a storage instance based on configuration
func NewStorage(cfg *config.StorageConfig) (Storage, error) {
	switch cfg.Type {
	case "s3":
		return NewS3Storage(
			cfg.S3Endpoint,
			cfg.S3AccessKey,
			cfg.S3SecretKey,
			cfg.S3Bucket,
			cfg.BasePath,
			cfg.S3UseSSL,
		)
	case "local", "filesystem":
		return NewLocalStorage(cfg.LocalPath)
	default:
		return nil, fmt.Errorf("unsupported storage type: %s", cfg.Type)
	}
}

