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

func shouldIgnore(path string, patterns []string) bool {
	// Simple glob pattern matching
	// TODO: Implement proper glob matching
	for _, pattern := range patterns {
		matched, err := filepath.Match(pattern, path)
		if err == nil && matched {
			return true
		}
		// Check if pattern matches any parent directory
		dir := filepath.Dir(path)
		for dir != "." && dir != "/" {
			matched, err := filepath.Match(pattern, dir)
			if err == nil && matched {
				return true
			}
			dir = filepath.Dir(dir)
		}
	}
	return false
}

