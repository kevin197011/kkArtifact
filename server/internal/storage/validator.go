// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package storage

import (
	"path/filepath"
	"strings"
)

// ValidatePath validates that a path is safe and prevents directory traversal
func ValidatePath(path string) error {
	// Clean the path to resolve any .. or . components
	cleaned := filepath.Clean(path)
	
	// Ensure the cleaned path doesn't start with .. which would indicate traversal
	if strings.HasPrefix(cleaned, "..") {
		return ErrPathTraversal
	}
	
	// Ensure path doesn't contain null bytes
	if strings.Contains(cleaned, "\x00") {
		return ErrInvalidPath
	}
	
	return nil
}

var (
	ErrPathTraversal = &StorageError{Message: "path contains directory traversal", Code: "PATH_TRAVERSAL"}
	ErrInvalidPath   = &StorageError{Message: "invalid path", Code: "INVALID_PATH"}
)

// StorageError represents a storage operation error
type StorageError struct {
	Message string
	Code    string
}

func (e *StorageError) Error() string {
	return e.Message
}

