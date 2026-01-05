// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package storage

import (
	"context"
	"io"
)

// Storage defines the interface for artifact storage backends
type Storage interface {
	// Put stores a file at the given path
	Put(ctx context.Context, path string, reader io.Reader, size int64) error

	// Get retrieves a file from the given path
	Get(ctx context.Context, path string) (io.ReadCloser, error)

	// Delete removes a file or directory at the given path
	Delete(ctx context.Context, path string) error

	// Exists checks if a path exists
	Exists(ctx context.Context, path string) (bool, error)

	// List lists files in a directory
	List(ctx context.Context, prefix string) ([]string, error)

	// Stat returns metadata about a file or directory
	Stat(ctx context.Context, path string) (*FileInfo, error)
}

// FileInfo contains metadata about a stored file
type FileInfo struct {
	Path    string
	Size    int64
	ModTime int64
	IsDir   bool
}

