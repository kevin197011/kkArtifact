// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package manifest

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

// Manifest represents the meta.yaml structure
type Manifest struct {
	Project   string         `yaml:"project"`
	App       string         `yaml:"app"`
	Version   string         `yaml:"version"`
	GitCommit string         `yaml:"git_commit,omitempty"`
	BuildTime string         `yaml:"build_time"`
	Builder   string         `yaml:"builder"`
	Files     []ManifestFile `yaml:"files"`
}

// ManifestFile represents a file entry
type ManifestFile struct {
	Path   string `yaml:"path"`
	SHA256 string `yaml:"sha256"`
	Size   int64  `yaml:"size"`
}

// Generate generates a manifest from a directory
func Generate(project, app, version, basePath string, ignorePatterns []string) (*Manifest, error) {
	manifest := &Manifest{
		Project:   project,
		App:       app,
		Version:   version,
		BuildTime: time.Now().Format(time.RFC3339),
		Builder:   "kkartifact-agent",
		Files:     []ManifestFile{},
	}

	err := filepath.Walk(basePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		// Get relative path
		relPath, err := filepath.Rel(basePath, path)
		if err != nil {
			return err
		}

		// Check ignore patterns
		if shouldIgnore(relPath, ignorePatterns) {
			return nil
		}

		// Calculate SHA256
		hash, size, err := calculateFileHash(path)
		if err != nil {
			return fmt.Errorf("failed to calculate hash for %s: %w", path, err)
		}

		manifest.Files = append(manifest.Files, ManifestFile{
			Path:   relPath,
			SHA256: hash,
			Size:   size,
		})

		return nil
	})

	return manifest, err
}

// Serialize serializes manifest to YAML bytes
func (m *Manifest) Serialize() ([]byte, error) {
	return yaml.Marshal(m)
}

// Parse parses manifest from YAML bytes
func Parse(data []byte) (*Manifest, error) {
	var manifest Manifest
	if err := yaml.Unmarshal(data, &manifest); err != nil {
		return nil, err
	}
	return &manifest, nil
}

func calculateFileHash(filePath string) (string, int64, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", 0, err
	}
	defer file.Close()

	hash := sha256.New()
	size, err := io.Copy(hash, file)
	if err != nil {
		return "", 0, err
	}

	return fmt.Sprintf("%x", hash.Sum(nil)), size, nil
}

// shouldIgnore checks if a path should be ignored based on glob patterns
// Supports standard glob patterns including:
// - * matches any sequence of non-separator characters
// - ** matches any sequence of characters including separators (recursive)
// - ? matches any single non-separator character
// - [abc] matches any character in the set
// - Directory patterns ending with / match all files in that directory and subdirectories
func shouldIgnore(path string, patterns []string) bool {
	// Normalize path separators to forward slashes for consistent matching
	normalizedPath := filepath.ToSlash(path)
	
	for _, pattern := range patterns {
		// Normalize pattern separators
		normalizedPattern := filepath.ToSlash(pattern)
		
		// Handle directory patterns (ending with /)
		if strings.HasSuffix(normalizedPattern, "/") {
			dirPrefix := strings.TrimSuffix(normalizedPattern, "/")
			// Match if path starts with the directory prefix
			if normalizedPath == dirPrefix || strings.HasPrefix(normalizedPath, dirPrefix+"/") {
				return true
			}
			continue
		}
		
		if matchGlob(normalizedPattern, normalizedPath) {
			return true
		}
		
		// Check if pattern matches any parent directory
		dir := filepath.ToSlash(filepath.Dir(path))
		for dir != "." && dir != "/" && dir != "" {
			if matchGlob(normalizedPattern, dir) {
				return true
			}
			parentDir := filepath.Dir(dir)
			if parentDir == dir {
				break // Reached root, stop
			}
			dir = filepath.ToSlash(parentDir)
		}
	}
	return false
}

// matchGlob implements glob pattern matching with support for ** recursive matching
func matchGlob(pattern, path string) bool {
	// Handle exact matches
	if pattern == path {
		return true
	}
	
	// Handle ** recursive matching
	if strings.Contains(pattern, "**") {
		return matchRecursive(pattern, path)
	}
	
	// Use standard filepath.Match for simple patterns
	matched, err := filepath.Match(pattern, path)
	return err == nil && matched
}

// matchRecursive handles patterns with ** for recursive directory matching
func matchRecursive(pattern, path string) bool {
	// Split pattern by ** to handle recursive matching
	parts := strings.Split(pattern, "**")
	
	if len(parts) == 1 {
		// No ** in pattern, use standard matching
		matched, err := filepath.Match(pattern, path)
		return err == nil && matched
	}
	
	// Pattern has **, need recursive matching
	// Example: "node_modules/**" should match "node_modules/file.js" and "node_modules/sub/file.js"
	// Example: "**/*.log" should match "test.log", "sub/test.log", "a/b/c/test.log"
	
	// Case 1: Pattern starts with ** (e.g., "**/*.log")
	if strings.HasPrefix(pattern, "**") {
		suffix := strings.TrimPrefix(pattern, "**")
		if suffix == "" || suffix == "/" {
			// Pattern is just "**" or "**/", matches everything
			return true
		}
		// Remove leading slash if present
		if strings.HasPrefix(suffix, "/") {
			suffix = suffix[1:]
		}
		// Check if path ends with suffix pattern
		return matchSuffix(suffix, path)
	}
	
	// Case 2: Pattern ends with ** (e.g., "node_modules/**")
	if strings.HasSuffix(pattern, "**") {
		prefix := strings.TrimSuffix(pattern, "**")
		// Remove trailing slash if present
		if strings.HasSuffix(prefix, "/") {
			prefix = prefix[:len(prefix)-1]
		}
		// Check if path starts with prefix
		return strings.HasPrefix(path, prefix+"/") || path == prefix
	}
	
	// Case 3: Pattern has ** in the middle (e.g., "src/**/*.js")
	// Split into prefix and suffix
	prefix := parts[0]
	suffix := strings.Join(parts[1:], "**")
	
	// Remove trailing/leading slashes
	if strings.HasSuffix(prefix, "/") {
		prefix = prefix[:len(prefix)-1]
	}
	if strings.HasPrefix(suffix, "/") {
		suffix = suffix[1:]
	}
	
	// Check if path starts with prefix and ends with suffix
	if !strings.HasPrefix(path, prefix+"/") && path != prefix {
		return false
	}
	
	// Extract the middle part between prefix and suffix
	middle := strings.TrimPrefix(path, prefix+"/")
	if path == prefix {
		middle = ""
	}
	
	if suffix == "" {
		return true // Pattern ends with **, matches everything after prefix
	}
	
	// Check if middle part contains suffix pattern
	return matchSuffix(suffix, middle)
}

// matchSuffix checks if path matches a suffix pattern (supports wildcards)
func matchSuffix(pattern, path string) bool {
	// Try matching from the end
	if strings.HasSuffix(path, pattern) {
		return true
	}
	
	// Try standard glob matching
	matched, err := filepath.Match(pattern, path)
	if err == nil && matched {
		return true
	}
	
	// Try matching at any position in the path
	// For patterns like "*.log", check if any suffix matches
	if strings.Contains(pattern, "*") || strings.Contains(pattern, "?") {
		parts := strings.Split(path, "/")
		for i := len(parts); i >= 0; i-- {
			suffix := strings.Join(parts[i:], "/")
			matched, err := filepath.Match(pattern, suffix)
			if err == nil && matched {
				return true
			}
		}
	}
	
	return false
}

