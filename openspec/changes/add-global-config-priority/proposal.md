# Change: Add Global Config Priority and Command-Line Override

## Why

Currently, the agent loads configuration from global config (`/etc/kkartifact/kkartifact.yml`) and local config (`.kkartifact.yml`), with local config overriding global config. However, there are limitations:

1. The global config path is hardcoded and doesn't match the user's expected path (`/etc/kkArtifact/config.yml`)
2. Command-line parameters cannot override configuration file values
3. The `ignore` field is completely replaced rather than merged, making it difficult to combine global and command-line ignore patterns

This change improves configuration flexibility by:
- Supporting the expected global config path (`/etc/kkArtifact/config.yml`)
- Allowing command-line parameters to override all configuration file values
- Merging `ignore` patterns from global config, local config, and command-line (with command-line taking precedence for duplicates)

## What Changes

- **MODIFIED**: Configuration loading priority: `/etc/kkArtifact/config.yml` → `.kkartifact.yml` → command-line flags
- **ADDED**: Command-line `--ignore` flag support for push/pull commands
- **MODIFIED**: `ignore` field merging logic: combine global + local + command-line, with command-line overriding duplicates
- **MODIFIED**: All configuration fields can be overridden by command-line flags (server_url, token, project, app, concurrency, etc.)

**BREAKING**: None - this is backward compatible. Existing config files continue to work.

## Impact

- **Affected specs**: `artifact-agent` (configuration loading and command-line interface)
- **Affected code**: 
  - `agent/internal/config/config.go` - Configuration loading and merging logic
  - `agent/internal/cli/push.go` - Add `--ignore` flag and override logic
  - `agent/internal/cli/pull.go` - Add `--ignore` flag and override logic
- **User impact**: Users can now use command-line flags to override config values and merge ignore patterns
