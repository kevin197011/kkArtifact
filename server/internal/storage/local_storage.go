// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package storage

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// LocalStorage implements Storage interface using local filesystem
type LocalStorage struct {
	basePath string
}

// NewLocalStorage creates a new local filesystem storage instance
func NewLocalStorage(basePath string) (*LocalStorage, error) {
	if err := os.MkdirAll(basePath, 0755); err != nil {
		return nil, err
	}
	return &LocalStorage{
		basePath: basePath,
	}, nil
}

// buildPath constructs the full path including base path
func (l *LocalStorage) buildPath(path string) string {
	cleanPath := strings.TrimPrefix(path, "/")
	return filepath.Join(l.basePath, cleanPath)
}

// Put stores a file at the given path
func (l *LocalStorage) Put(ctx context.Context, path string, reader io.Reader, size int64) error {
	fullPath := l.buildPath(path)
	
	// Create directory if it doesn't exist
	dir := filepath.Dir(fullPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	file, err := os.Create(fullPath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, reader)
	return err
}

// Get retrieves a file from the given path
func (l *LocalStorage) Get(ctx context.Context, path string) (io.ReadCloser, error) {
	fullPath := l.buildPath(path)
	return os.Open(fullPath)
}

// Delete removes a file or directory at the given path
func (l *LocalStorage) Delete(ctx context.Context, path string) error {
	fullPath := l.buildPath(path)
	return os.RemoveAll(fullPath)
}

// Exists checks if a path exists
func (l *LocalStorage) Exists(ctx context.Context, path string) (bool, error) {
	fullPath := l.buildPath(path)
	_, err := os.Stat(fullPath)
	if os.IsNotExist(err) {
		return false, nil
	}
	return err == nil, err
}

// List lists files in a directory
func (l *LocalStorage) List(ctx context.Context, prefix string) ([]string, error) {
	fullPath := l.buildPath(prefix)
	
	var paths []string
	err := filepath.Walk(fullPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		
		// Get relative path from base
		relPath, err := filepath.Rel(l.basePath, path)
		if err != nil {
			return err
		}
		
		if relPath != "." {
			paths = append(paths, relPath)
		}
		return nil
	})
	
	return paths, err
}

// Stat returns metadata about a file or directory
func (l *LocalStorage) Stat(ctx context.Context, path string) (*FileInfo, error) {
	fullPath := l.buildPath(path)
	info, err := os.Stat(fullPath)
	if err != nil {
		return nil, err
	}

	return &FileInfo{
		Path:    path,
		Size:    info.Size(),
		ModTime: info.ModTime().Unix(),
		IsDir:   info.IsDir(),
	}, nil
}

