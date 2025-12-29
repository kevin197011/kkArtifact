// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package cache

import "fmt"

const (
	keyPrefix = "kkartifact:"

	// Cache key patterns
	keyProjectsList      = keyPrefix + "projects:list"
	keyAppsList          = keyPrefix + "apps:list:project:%d"
	keyVersionsList      = keyPrefix + "versions:list:app:%d"
	keyLatestVersion     = keyPrefix + "version:latest:app:%d"
	keyManifest          = keyPrefix + "manifest:%s:%s:%s"
	keyConfig            = keyPrefix + "config:%s"
)

// ProjectsListKey returns cache key for projects list
func ProjectsListKey() string {
	return keyProjectsList
}

// AppsListKey returns cache key for apps list
func AppsListKey(projectID int) string {
	return fmt.Sprintf(keyAppsList, projectID)
}

// VersionsListKey returns cache key for versions list
func VersionsListKey(appID int) string {
	return fmt.Sprintf(keyVersionsList, appID)
}

// LatestVersionKey returns cache key for latest version
func LatestVersionKey(appID int) string {
	return fmt.Sprintf(keyLatestVersion, appID)
}

// ManifestKey returns cache key for manifest
func ManifestKey(project, app, version string) string {
	return fmt.Sprintf(keyManifest, project, app, version)
}

// ConfigKey returns cache key for config
func ConfigKey(key string) string {
	return fmt.Sprintf(keyConfig, key)
}

