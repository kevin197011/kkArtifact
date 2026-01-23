// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package config

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"unicode"

	"gopkg.in/yaml.v3"
)

// Config represents the agent configuration
type Config struct {
	ServerURL      string   `yaml:"server_url"`
	Token          string   `yaml:"token"`
	Project        string   `yaml:"project"`
	App            string   `yaml:"app"`
	Ignore         []string `yaml:"ignore,omitempty"`
	RetainVersions *int     `yaml:"retain_versions,omitempty"`
	Concurrency    int      `yaml:"concurrency"` // Number of concurrent uploads/downloads (default: 50)
}

// GetGlobalConfigPath returns the path to the global configuration file
// Unix/Linux/macOS: Tries /etc/kkArtifact/config.yml first (with capital A), then falls back to /etc/kkartifact/kkartifact.yml
// Windows: Uses C:\ProgramData\kkArtifact\config.yml
func GetGlobalConfigPath() (string, error) {
	if runtime.GOOS == "windows" {
		// Windows: Use ProgramData directory
		programData := os.Getenv("ProgramData")
		if programData == "" {
			programData = "C:\\ProgramData"
		}
		configPath := filepath.Join(programData, "kkArtifact", "config.yml")
		return configPath, nil
	}
	
	// Unix-like systems: Try primary path first: /etc/kkArtifact/config.yml (with capital A)
	primaryPath := "/etc/kkArtifact/config.yml"
	if _, err := os.Stat(primaryPath); err == nil {
		return primaryPath, nil
	}
	
	// Fallback to legacy path: /etc/kkartifact/kkartifact.yml (lowercase, backward compatibility)
	fallbackPath := "/etc/kkartifact/kkartifact.yml"
	return fallbackPath, nil
}

// loadConfigFile loads a single configuration file
func loadConfigFile(configPath string) (*Config, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Clean token: remove all whitespace (including newlines, tabs, etc.)
	// YAML parser should handle inline comments correctly, but we need to be extra careful
	config.Token = cleanTokenValue(config.Token)
	config.ServerURL = strings.TrimSpace(config.ServerURL)

	return &config, nil
}

// cleanTokenValue cleans a token value by removing whitespace and handling edge cases
// This function is defensive and handles cases where YAML parser might include comment text
func cleanTokenValue(token string) string {
	if token == "" {
		return ""
	}
	
	// Remove UTF-8 BOM if present
	if len(token) >= 3 && token[0] == 0xEF && token[1] == 0xBB && token[2] == 0xBF {
		token = token[3:]
	}
	
	// Remove all types of whitespace and invisible characters
	// Use a more comprehensive approach to remove all Unicode whitespace
	var cleaned strings.Builder
	for _, r := range token {
		// Skip all Unicode whitespace characters and zero-width characters
		if !unicode.IsSpace(r) && r != '\u200B' && r != '\u200C' && r != '\u200D' && r != '\uFEFF' {
			cleaned.WriteRune(r)
		}
	}
	token = cleaned.String()
	
	// Remove any inline comments that might have been parsed as part of the value
	// YAML parser should handle this correctly, but we'll be extra safe
	// Look for comment patterns: space followed by # followed by comment text
	// Only remove if we're confident it's a comment (not part of the token value)
	if idx := strings.Index(token, " #"); idx >= 0 {
		beforeHash := token[:idx]
		afterHash := strings.TrimSpace(token[idx+2:])
		// Check if the part after # looks like a comment (common comment phrases)
		// Base64 URL-encoded tokens shouldn't contain spaces, so " #" is likely a comment separator
		commentPhrases := []string{
			"Uncomment",
			"set your token",
			"YOUR_TOKEN",
			"Uncomment and",
			"token here",
		}
		for _, phrase := range commentPhrases {
			if strings.HasPrefix(afterHash, phrase) {
				// This looks like a comment, remove it
				token = strings.TrimSpace(beforeHash)
				break
			}
		}
	}
	
	// Also check for # at the start (shouldn't happen with proper YAML, but be safe)
	// If token starts with #, it's likely a mis-parsed comment line
	if strings.HasPrefix(token, "#") {
		return ""
	}
	
	// Final trim to remove any remaining leading/trailing whitespace
	return strings.TrimSpace(token)
}

// ValidateTokenFormat validates that a token matches the expected base64 URL encoding pattern
// Base64 URL encoding uses: A-Z, a-z, 0-9, - (hyphen), _ (underscore)
// Padding characters (=) are allowed only at the end (0, 1, or 2 padding characters)
func ValidateTokenFormat(token string) error {
	if token == "" {
		return fmt.Errorf("token is empty")
	}
	
	// Base64 URL encoding pattern: ^[A-Za-z0-9_-]+={0,2}$
	// Allows padding characters (=) only at the end (0, 1, or 2 padding characters)
	// Tokens are typically 32-44 characters (32 bytes = 43 chars, 44 with padding)
	base64URLPattern := regexp.MustCompile(`^[A-Za-z0-9_-]+={0,2}$`)
	if !base64URLPattern.MatchString(token) {
		return fmt.Errorf("token format is invalid: tokens must be base64 URL encoded strings (A-Z, a-z, 0-9, -, _), with optional = padding at the end")
	}
	
	// Check length (base64 encoding of 32 bytes = 43 chars, with padding = 44)
	// Allow some flexibility: 20-100 characters
	if len(token) < 20 {
		return fmt.Errorf("token is too short: expected at least 20 characters, got %d", len(token))
	}
	if len(token) > 100 {
		return fmt.Errorf("token is too long: expected at most 100 characters, got %d", len(token))
	}
	
	return nil
}

// MaskToken returns a masked version of a token for display in error messages
// Shows first 5 and last 5 characters: "YVza5...b_rc="
func MaskToken(token string) string {
	if token == "" {
		return "<empty>"
	}
	if len(token) <= 10 {
		return strings.Repeat("*", len(token))
	}
	return token[:5] + "..." + token[len(token)-5:]
}

// mergeIgnorePatterns merges ignore patterns from multiple sources
// Priority: global → local → command-line
// Removes duplicates, keeping the last occurrence (command-line patterns appear last)
func mergeIgnorePatterns(global, local, commandLine []string) []string {
	// Combine all patterns in priority order
	allPatterns := make([]string, 0)
	
	// Add global patterns first
	if global != nil {
		allPatterns = append(allPatterns, global...)
	}
	
	// Add local patterns
	if local != nil {
		allPatterns = append(allPatterns, local...)
	}
	
	// Add command-line patterns last
	if commandLine != nil {
		allPatterns = append(allPatterns, commandLine...)
	}
	
	// Remove duplicates, keeping the last occurrence
	seen := make(map[string]bool)
	result := make([]string, 0)
	
	// Iterate in reverse to keep the last occurrence of each pattern
	for i := len(allPatterns) - 1; i >= 0; i-- {
		pattern := allPatterns[i]
		if !seen[pattern] {
			seen[pattern] = true
			result = append([]string{pattern}, result...)
		}
	}
	
	return result
}

// mergeConfigs merges global config with local config
// Local config values override global config values
// This is the legacy function for backward compatibility
func mergeConfigs(global, local *Config) *Config {
	return mergeConfigsWithOverrides(global, local, nil)
}

// Overrides represents command-line overrides for configuration
type Overrides struct {
	ServerURL   string
	Token       string
	Project     string
	App         string
	Ignore      []string
	Concurrency int // 0 means not set
}

// mergeConfigsWithOverrides merges global config, local config, and command-line overrides
// Priority: global config → local config → command-line overrides
func mergeConfigsWithOverrides(global, local *Config, overrides *Overrides) *Config {
	result := &Config{}

	// Start with global config
	if global != nil {
		result.ServerURL = strings.TrimSpace(global.ServerURL)
		result.Token = cleanTokenValue(global.Token)
		result.Project = global.Project
		result.App = global.App
		result.Ignore = global.Ignore
		result.RetainVersions = global.RetainVersions
		result.Concurrency = global.Concurrency
	}

	// Override with local config (if present)
	// Only override if local config has non-empty values (empty strings don't override)
	if local != nil {
		if local.ServerURL != "" {
			result.ServerURL = strings.TrimSpace(local.ServerURL)
		}
		if local.Token != "" {
			// Only override if local token is non-empty (preserve global token if local is empty)
			result.Token = cleanTokenValue(local.Token)
		}
		if local.Project != "" {
			result.Project = local.Project
		}
		if local.App != "" {
			result.App = local.App
		}
		if local.RetainVersions != nil {
			result.RetainVersions = local.RetainVersions
		}
		if local.Concurrency > 0 {
			result.Concurrency = local.Concurrency
		}
		// Note: ignore patterns are merged separately below
	}

	// Apply command-line overrides (highest priority)
	if overrides != nil {
		if overrides.ServerURL != "" {
			result.ServerURL = strings.TrimSpace(overrides.ServerURL)
		}
		if overrides.Token != "" {
			result.Token = cleanTokenValue(overrides.Token)
		}
		if overrides.Project != "" {
			result.Project = overrides.Project
		}
		if overrides.App != "" {
			result.App = overrides.App
		}
		if overrides.Concurrency > 0 {
			result.Concurrency = overrides.Concurrency
		}
	}

	// Merge ignore patterns: global → local → command-line
	var globalIgnore, localIgnore, cmdIgnore []string
	if global != nil {
		globalIgnore = global.Ignore
	}
	if local != nil && local.Ignore != nil {
		localIgnore = local.Ignore
	}
	if overrides != nil && overrides.Ignore != nil {
		cmdIgnore = overrides.Ignore
	}
	
	// If command-line has ignore patterns, merge all sources
	if cmdIgnore != nil {
		result.Ignore = mergeIgnorePatterns(globalIgnore, localIgnore, cmdIgnore)
	} else if localIgnore != nil {
		// If local config has ignore (even if empty), use it (allows clearing global ignore)
		result.Ignore = mergeIgnorePatterns(globalIgnore, localIgnore, nil)
	} else {
		// Use global ignore only
		result.Ignore = globalIgnore
	}

	return result
}

// Load loads configuration with priority: global config → local config → command-line overrides
// If configPath is empty or ".kkartifact.yml", it will try to load from current directory
// Global config is loaded from /etc/kkArtifact/config.yml (or /etc/kkartifact/kkartifact.yml as fallback)
// If overrides is nil, Load behaves the same as before (backward compatible)
func Load(configPath string, overrides *Overrides) (*Config, error) {
	var globalConfig *Config
	var localConfig *Config
	var err error

	// Try to load global config (ignore errors if it doesn't exist)
	globalConfigPath, err := GetGlobalConfigPath()
	if err == nil {
		if runtime.GOOS == "windows" {
			// Windows: Use the path directly
			if cfg, err := loadConfigFile(globalConfigPath); err == nil {
				globalConfig = cfg
			}
			// Ignore error if global config doesn't exist
		} else {
			// Unix-like: Try primary path first
			primaryPath := "/etc/kkArtifact/config.yml"
			if cfg, err := loadConfigFile(primaryPath); err == nil {
				globalConfig = cfg
			} else {
				// Fallback to legacy path
				if cfg, err := loadConfigFile(globalConfigPath); err == nil {
					globalConfig = cfg
				}
				// Ignore error if global config doesn't exist
			}
		}
	}

	// Load local config (required if no global config or overrides)
	// If configPath is empty or ".kkartifact.yml", use current directory
	if configPath == "" || configPath == ".kkartifact.yml" {
		wd, err := os.Getwd()
		if err != nil {
			return nil, fmt.Errorf("failed to get current directory: %w", err)
		}
		configPath = filepath.Join(wd, ".kkartifact.yml")
	}

	localConfig, err = loadConfigFile(configPath)
	if err != nil {
		// If local config doesn't exist and global config exists, use global config
		if globalConfig != nil {
			// Apply overrides if provided
			mergedConfig := mergeConfigsWithOverrides(globalConfig, nil, overrides)
			
			// Validate merged config
			if mergedConfig.ServerURL == "" {
				return nil, fmt.Errorf("server_url is required in global config")
			}
			if mergedConfig.Token == "" {
				return nil, fmt.Errorf("token is required in global config")
			}
			// Set default concurrency if not specified
			if mergedConfig.Concurrency <= 0 {
				mergedConfig.Concurrency = 50
			}
			return mergedConfig, nil
		}
		// If overrides provide required fields, we can proceed without config files
		if overrides != nil && overrides.ServerURL != "" && overrides.Token != "" {
			mergedConfig := mergeConfigsWithOverrides(nil, nil, overrides)
			if mergedConfig.Concurrency <= 0 {
				mergedConfig.Concurrency = 50
			}
			return mergedConfig, nil
		}
		return nil, fmt.Errorf("failed to load config file %s: %w", configPath, err)
	}

	// Merge configs: global → local → command-line overrides
	mergedConfig := mergeConfigsWithOverrides(globalConfig, localConfig, overrides)

	// Validate required fields
	if mergedConfig.ServerURL == "" {
		return nil, fmt.Errorf("server_url is required")
	}
	if mergedConfig.Token == "" {
		return nil, fmt.Errorf("token is required")
	}

	// Set default concurrency if not specified or invalid
	if mergedConfig.Concurrency <= 0 {
		mergedConfig.Concurrency = 50 // Default to 50 concurrent operations
	}

	return mergedConfig, nil
}

