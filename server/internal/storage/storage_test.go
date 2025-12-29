// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package storage

import (
	"context"
	"io"
	"strings"
	"testing"
)

func TestLocalStorage_PutGet(t *testing.T) {
	storage, err := NewLocalStorage("/tmp/test-storage")
	if err != nil {
		t.Fatalf("Failed to create local storage: %v", err)
	}

	ctx := context.Background()
	testPath := "test/file.txt"
	testContent := "test content"

	// Put file
	reader := strings.NewReader(testContent)
	err = storage.Put(ctx, testPath, reader, int64(len(testContent)))
	if err != nil {
		t.Fatalf("Failed to put file: %v", err)
	}

	// Get file
	fileReader, err := storage.Get(ctx, testPath)
	if err != nil {
		t.Fatalf("Failed to get file: %v", err)
	}
	defer fileReader.Close()

	content, err := io.ReadAll(fileReader)
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}

	if string(content) != testContent {
		t.Errorf("Expected %s, got %s", testContent, string(content))
	}
}

func TestLocalStorage_Exists(t *testing.T) {
	storage, err := NewLocalStorage("/tmp/test-storage")
	if err != nil {
		t.Fatalf("Failed to create local storage: %v", err)
	}

	ctx := context.Background()
	testPath := "test/exists.txt"

	// File should not exist initially
	exists, err := storage.Exists(ctx, testPath)
	if err != nil {
		t.Fatalf("Failed to check existence: %v", err)
	}
	if exists {
		t.Error("File should not exist")
	}

	// Create file
	reader := strings.NewReader("test")
	err = storage.Put(ctx, testPath, reader, 4)
	if err != nil {
		t.Fatalf("Failed to put file: %v", err)
	}

	// File should exist now
	exists, err = storage.Exists(ctx, testPath)
	if err != nil {
		t.Fatalf("Failed to check existence: %v", err)
	}
	if !exists {
		t.Error("File should exist")
	}
}

func TestValidatePath(t *testing.T) {
	tests := []struct {
		path    string
		wantErr bool
	}{
		{"valid/path.txt", false},
		{"../../etc/passwd", true},
		{"../parent/file.txt", true},
		{"normal/file.txt", false},
		{"file\x00withnull.txt", true},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			err := ValidatePath(tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidatePath() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

