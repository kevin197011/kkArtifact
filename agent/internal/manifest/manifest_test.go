// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package manifest

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGenerate(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir := t.TempDir()

	// Create test files
	testFile1 := filepath.Join(tmpDir, "file1.txt")
	os.WriteFile(testFile1, []byte("content1"), 0644)

	subDir := filepath.Join(tmpDir, "subdir")
	os.MkdirAll(subDir, 0755)
	testFile2 := filepath.Join(subDir, "file2.txt")
	os.WriteFile(testFile2, []byte("content2"), 0644)

	// Generate manifest
	manifest, err := Generate("test-project", "test-app", "v1.0.0", tmpDir, []string{})
	if err != nil {
		t.Fatalf("Failed to generate manifest: %v", err)
	}

	if manifest.Project != "test-project" {
		t.Errorf("Expected project 'test-project', got '%s'", manifest.Project)
	}

	if manifest.App != "test-app" {
		t.Errorf("Expected app 'test-app', got '%s'", manifest.App)
	}

	if len(manifest.Files) < 2 {
		t.Errorf("Expected at least 2 files, got %d", len(manifest.Files))
	}

	// Check that files have valid hashes
	for _, file := range manifest.Files {
		if file.SHA256 == "" {
			t.Error("File hash is empty")
		}
		if file.Size == 0 {
			t.Error("File size is 0")
		}
	}
}

func TestShouldIgnore(t *testing.T) {
	tests := []struct {
		path     string
		patterns []string
		want     bool
	}{
		// Basic patterns
		{"test.log", []string{"*.log"}, true},
		{"test.txt", []string{"*.log"}, false},
		{"test.log", []string{"test.*"}, true},
		
		// Directory patterns with **
		{"node_modules/file.js", []string{"node_modules/**"}, true},
		{"node_modules/subdir/file.js", []string{"node_modules/**"}, true},
		{"node_modules/a/b/c/file.js", []string{"node_modules/**"}, true},
		{"src/file.js", []string{"node_modules/**"}, false},
		
		// Recursive file patterns
		{"test.log", []string{"**/*.log"}, true},
		{"subdir/test.log", []string{"**/*.log"}, true},
		{"a/b/c/test.log", []string{"**/*.log"}, true},
		{"test.txt", []string{"**/*.log"}, false},
		
		// Directory prefix patterns
		{"logs/app.log", []string{"logs/"}, true},
		{"logs/subdir/app.log", []string{"logs/"}, true},
		{"tmp/file.tmp", []string{"tmp/"}, true},
		{"other/file.tmp", []string{"tmp/"}, false},
		
		// Multiple patterns
		{"test.log", []string{"*.log", "*.tmp"}, true},
		{"test.tmp", []string{"*.log", "*.tmp"}, true},
		{"test.txt", []string{"*.log", "*.tmp"}, false},
		
		// Complex patterns
		{"src/main.js", []string{"src/**/*.js"}, true},
		{"src/utils/helper.js", []string{"src/**/*.js"}, true},
		{"src/a/b/c/file.js", []string{"src/**/*.js"}, true},
		{"other/main.js", []string{"src/**/*.js"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			got := shouldIgnore(tt.path, tt.patterns)
			if got != tt.want {
				t.Errorf("shouldIgnore(%q, %v) = %v, want %v", tt.path, tt.patterns, got, tt.want)
			}
		})
	}
}

