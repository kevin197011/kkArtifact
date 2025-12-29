// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package storage

import (
	"context"
	"fmt"
	"io"
	"path/filepath"
	"strings"
)

// ArtifactManager manages artifact versions
type ArtifactManager struct {
	storage Storage
}

// NewArtifactManager creates a new artifact manager
func NewArtifactManager(storage Storage) *ArtifactManager {
	return &ArtifactManager{
		storage: storage,
	}
}

// StoreVersion stores an artifact version
func (am *ArtifactManager) StoreVersion(ctx context.Context, project, app, version string, manifest *Manifest, files map[string]io.Reader) error {
	versionPath := am.versionPath(project, app, version)
	
	// Check if version already exists (immutability)
	exists, err := am.storage.Exists(ctx, filepath.Join(versionPath, "meta.yaml"))
	if err != nil {
		return fmt.Errorf("failed to check version existence: %w", err)
	}
	if exists {
		return fmt.Errorf("version %s already exists and cannot be overwritten", version)
	}

	// Store files
	for filePath, reader := range files {
		fullPath := filepath.Join(versionPath, filePath)
		
		// Validate path
		if err := ValidatePath(fullPath); err != nil {
			return fmt.Errorf("invalid path %s: %w", filePath, err)
		}
		
		// Get file size if possible
		var size int64
		if seeker, ok := reader.(io.Seeker); ok {
			if currentPos, err := seeker.Seek(0, io.SeekCurrent); err == nil {
				if endPos, err := seeker.Seek(0, io.SeekEnd); err == nil {
					size = endPos - currentPos
					seeker.Seek(currentPos, io.SeekStart)
				}
			}
		}
		
		if err := am.storage.Put(ctx, fullPath, reader, size); err != nil {
			return fmt.Errorf("failed to store file %s: %w", filePath, err)
		}
	}

	// Store manifest
	manifestBytes, err := SerializeManifest(manifest)
	if err != nil {
		return fmt.Errorf("failed to serialize manifest: %w", err)
	}
	
	manifestPath := filepath.Join(versionPath, "meta.yaml")
	manifestReader := io.NopCloser(io.Reader(io.Reader(strings.NewReader(string(manifestBytes)))))
	if err := am.storage.Put(ctx, manifestPath, manifestReader, int64(len(manifestBytes))); err != nil {
		return fmt.Errorf("failed to store manifest: %w", err)
	}

	return nil
}

// GetManifest retrieves manifest for a version
func (am *ArtifactManager) GetManifest(ctx context.Context, project, app, version string) (*Manifest, error) {
	manifestPath := filepath.Join(am.versionPath(project, app, version), "meta.yaml")
	
	reader, err := am.storage.Get(ctx, manifestPath)
	if err != nil {
		return nil, fmt.Errorf("failed to get manifest: %w", err)
	}
	defer reader.Close()

	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to read manifest: %w", err)
	}

	return ParseManifest(data)
}

// DeleteVersion deletes an artifact version
func (am *ArtifactManager) DeleteVersion(ctx context.Context, project, app, version string) error {
	versionPath := am.versionPath(project, app, version)
	return am.storage.Delete(ctx, versionPath)
}

// ListVersions lists all versions for an app
func (am *ArtifactManager) ListVersions(ctx context.Context, project, app string) ([]string, error) {
	appPath := am.appPath(project, app)
	
	entries, err := am.storage.List(ctx, appPath)
	if err != nil {
		return nil, fmt.Errorf("failed to list versions: %w", err)
	}

	var versions []string
	for _, entry := range entries {
		// Extract version from path (entries are like "version/meta.yaml" or "version/")
		baseName := filepath.Base(entry)
		ext := filepath.Ext(entry)
		if baseName == "meta.yaml" || ext == "" {
			version := filepath.Base(filepath.Dir(entry))
			if version != "." && version != appPath {
				versions = append(versions, version)
			}
		}
	}

	return versions, nil
}

// versionPath returns the storage path for a version
func (am *ArtifactManager) versionPath(project, app, version string) string {
	return filepath.Join(am.appPath(project, app), version)
}

// appPath returns the storage path for an app
func (am *ArtifactManager) appPath(project, app string) string {
	return filepath.Join(project, app)
}

// DeleteApp deletes an app and all its versions from storage
func (am *ArtifactManager) DeleteApp(ctx context.Context, project, app string) error {
	appPath := am.appPath(project, app)
	return am.storage.Delete(ctx, appPath)
}

// DeleteProject deletes a project and all its apps and versions from storage
func (am *ArtifactManager) DeleteProject(ctx context.Context, project string) error {
	return am.storage.Delete(ctx, project)
}

