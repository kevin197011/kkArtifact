// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/kk/kkartifact-agent/internal/client"
	"github.com/kk/kkartifact-agent/internal/config"
	"github.com/kk/kkartifact-agent/internal/manifest"
)

// pushParams holds resolved push target and runtime configuration
type pushParams struct {
	project   string
	app       string
	version   string
	path      string
	cfg       *config.Config
	publish   bool
	apiClient *client.Client // optional: reuse across push-tree iterations
}

// pushTreeTarget is a single app/version directory discovered under a push-tree root
type pushTreeTarget struct {
	app     string
	version string
	path    string
}

// inferPushTargetFromPath infers app and version from a path like .../G02_agent_api/cd884eb...
func inferPushTargetFromPath(absPath string) (app, version string, ok bool) {
	clean := filepath.Clean(absPath)
	version = filepath.Base(clean)
	app = filepath.Base(filepath.Dir(clean))
	if version == "" || app == "" || app == "." || version == "." {
		return "", "", false
	}
	return app, version, true
}

// resolvePushParams loads config and resolves project/app/version from flags and path
func resolvePushParams(
	project, app, version, path, configPath string,
	serverURL, token string,
	concurrency int,
	ignore []string,
) (*pushParams, error) {
	cfg, err := loadAgentConfig(configPath, serverURL, token, concurrency, ignore)
	if err != nil {
		return nil, err
	}

	if project == "" {
		project = cfg.Project
	}
	if app == "" {
		app = cfg.App
	}

	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve path: %w", err)
	}

	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("path does not exist: %s", absPath)
	}

	if app == "" || version == "" {
		inferredApp, inferredVersion, ok := inferPushTargetFromPath(absPath)
		if !ok {
			return nil, fmt.Errorf("cannot infer app/version from path %s; use --app and --version flags", absPath)
		}
		if app == "" {
			app = inferredApp
		}
		if version == "" {
			version = inferredVersion
		}
	}

	if project == "" || app == "" || version == "" {
		return nil, fmt.Errorf("project, app, and version are required (via flags, config file, or path inference)")
	}

	return &pushParams{
		project: project,
		app:     app,
		version: version,
		path:    absPath,
		cfg:     cfg,
	}, nil
}

func parseIgnorePatterns(ignore []string) []string {
	patterns := make([]string, 0)
	for _, ignoreFlag := range ignore {
		for _, pattern := range strings.Split(ignoreFlag, ",") {
			trimmed := strings.TrimSpace(pattern)
			if trimmed != "" {
				patterns = append(patterns, trimmed)
			}
		}
	}
	return patterns
}

// loadAgentConfig loads and validates agent config with CLI overrides
func loadAgentConfig(configPath, serverURL, token string, concurrency int, ignore []string) (*config.Config, error) {
	ignorePatterns := parseIgnorePatterns(ignore)

	overrides := &config.Overrides{
		ServerURL:   serverURL,
		Token:       token,
		Concurrency: concurrency,
		Ignore:      ignorePatterns,
	}
	if len(ignorePatterns) == 0 {
		overrides.Ignore = nil
	}
	if concurrency == 0 {
		overrides.Concurrency = 0
	}

	cfg, err := config.Load(configPath, overrides)
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	if cfg.Token == "" {
		return nil, fmt.Errorf("token is required but not found in config. Please check:\n  - Global config: /etc/kkArtifact/config.yml\n  - Local config: %s\n  - Or use --token flag", configPath)
	}

	if err := config.ValidateTokenFormat(cfg.Token); err != nil {
		return nil, fmt.Errorf("token validation failed: %w\nToken preview: %s\nConfig file: %s", err, config.MaskToken(cfg.Token), configPath)
	}

	return cfg, nil
}

// loadPushConfig loads agent config and resolves project name
func loadPushConfig(
	project, configPath, serverURL, token string,
	concurrency int,
	ignore []string,
) (*config.Config, string, error) {
	cfg, err := loadAgentConfig(configPath, serverURL, token, concurrency, ignore)
	if err != nil {
		return nil, "", err
	}

	if project == "" {
		project = cfg.Project
	}
	if project == "" {
		return nil, "", fmt.Errorf("project is required (use --project or set project in config)")
	}

	return cfg, project, nil
}

// versionExistsOnServer checks whether a version is already present
func versionExistsOnServer(apiClient *client.Client, project, app, version string) bool {
	_, err := apiClient.GetManifest(project, app, version)
	return err == nil
}

// filterFilesToUpload returns files that are new or changed compared to remote manifest
func filterFilesToUpload(local *manifest.Manifest, remote *client.Manifest) ([]manifest.ManifestFile, bool) {
	remoteHashes := remote.FileHashes()
	if len(remoteHashes) == 0 {
		return local.Files, false
	}

	unchanged := true
	toUpload := make([]manifest.ManifestFile, 0, len(local.Files))
	for _, f := range local.Files {
		rh, ok := remoteHashes[f.Path]
		if ok && strings.EqualFold(rh, f.SHA256) {
			continue
		}
		unchanged = false
		toUpload = append(toUpload, f)
	}

	if unchanged && len(local.Files) == len(remoteHashes) {
		return nil, true
	}
	return toUpload, false
}
