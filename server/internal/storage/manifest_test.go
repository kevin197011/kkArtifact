// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package storage

import (
	"crypto/sha256"
	"strings"
	"testing"
)

func TestSerializeManifest(t *testing.T) {
	manifest := &Manifest{
		Project:   "test-project",
		App:       "test-app",
		Version:   "v1.0.0",
		GitCommit: "abc123",
		Builder:   "test-builder",
		Files: []ManifestFile{
			{
				Path:   "file1.txt",
				SHA256: "hash1",
				Size:   100,
			},
		},
	}

	data, err := SerializeManifest(manifest)
	if err != nil {
		t.Fatalf("Failed to serialize manifest: %v", err)
	}

	if len(data) == 0 {
		t.Error("Serialized manifest is empty")
	}

	// Parse back
	parsed, err := ParseManifest(data)
	if err != nil {
		t.Fatalf("Failed to parse manifest: %v", err)
	}

	if parsed.Project != manifest.Project {
		t.Errorf("Expected project %s, got %s", manifest.Project, parsed.Project)
	}

	if len(parsed.Files) != len(manifest.Files) {
		t.Errorf("Expected %d files, got %d", len(manifest.Files), len(parsed.Files))
	}
}

func TestCalculateSHA256(t *testing.T) {
	testContent := "test content"

	reader := strings.NewReader(testContent)
	hash, err := CalculateSHA256(reader)
	if err != nil {
		t.Fatalf("Failed to calculate hash: %v", err)
	}

	// Calculate expected hash
	hasher := sha256.New()
	hasher.Write([]byte(testContent))
	expectedBytes := hasher.Sum(nil)
	expectedHash := ""
	for _, b := range expectedBytes {
		expectedHash += string("0123456789abcdef"[b>>4])
		expectedHash += string("0123456789abcdef"[b&0x0f])
	}

	if len(hash) != 64 { // SHA256 hex string length
		t.Errorf("Expected hash length 64, got %d", len(hash))
	}

	// Verify it matches expected
	if hash != expectedHash {
		t.Errorf("Hash mismatch: got %s, expected %s", hash, expectedHash)
	}
}

