// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package storage

import (
	"crypto/sha256"
	"fmt"
	"io"
	"time"

	"gopkg.in/yaml.v3"
)

// Manifest represents the meta.yaml file structure
type Manifest struct {
	Project   string         `yaml:"project"`
	App       string         `yaml:"app"`
	Version   string         `yaml:"version"`
	GitCommit string         `yaml:"git_commit,omitempty"`
	BuildTime string         `yaml:"build_time"`
	Builder   string         `yaml:"builder"`
	Files     []ManifestFile `yaml:"files"`
}

// ManifestFile represents a file entry in the manifest
type ManifestFile struct {
	Path   string `yaml:"path"`
	SHA256 string `yaml:"sha256"`
	Size   int64  `yaml:"size"`
}

// CalculateSHA256 calculates SHA256 hash of the given reader
func CalculateSHA256(reader io.Reader) (string, error) {
	hash := sha256.New()
	if _, err := io.Copy(hash, reader); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}

// CalculateSHA256Parallel calculates SHA256 hash of a large file in parallel chunks
func CalculateSHA256Parallel(reader io.Reader, chunkSize int64) (string, error) {
	// For now, use sequential calculation
	// TODO: Implement parallel computation for large files
	return CalculateSHA256(reader)
}

// SerializeManifest serializes manifest to YAML bytes
func SerializeManifest(manifest *Manifest) ([]byte, error) {
	manifest.BuildTime = time.Now().Format(time.RFC3339)
	return yaml.Marshal(manifest)
}

// ParseManifest parses manifest from YAML bytes
func ParseManifest(data []byte) (*Manifest, error) {
	var manifest Manifest
	if err := yaml.Unmarshal(data, &manifest); err != nil {
		return nil, err
	}
	return &manifest, nil
}

