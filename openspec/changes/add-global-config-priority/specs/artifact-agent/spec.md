## MODIFIED Requirements

### Requirement: Agent Configuration

The agent SHALL support configuration via `.kkartifact.yml` including server URL, token, concurrency level, chunk size, ignore rules, and optional version retention limit. The agent SHALL load configuration with the following priority order: global config (`/etc/kkArtifact/config.yml` or `/etc/kkartifact/kkartifact.yml`) → local config (`.kkartifact.yml`) → command-line flags. Command-line flags SHALL override configuration file values.

#### Scenario: Load configuration with global config priority
- **WHEN** the agent starts, it loads configuration
- **THEN** it first attempts to load `/etc/kkArtifact/config.yml` (with capital A)
- **AND** if that file doesn't exist, it falls back to `/etc/kkartifact/kkartifact.yml` (lowercase)
- **AND** then loads local config `.kkartifact.yml` from current directory (if exists)
- **AND** command-line flags override any configuration file values
- **AND** server URL, token, concurrency, chunk_size, ignore rules, and optional retain_versions are read

#### Scenario: Command-line override server URL
- **WHEN** global config has `server_url: https://packages.example.com/` and command-line has `--server-url https://other.example.com/`
- **THEN** the agent uses `https://other.example.com/` from command-line

#### Scenario: Command-line override token
- **WHEN** local config has `token: abc123` and command-line has `--token xyz789`
- **THEN** the agent uses `xyz789` from command-line

#### Scenario: Command-line override concurrency
- **WHEN** global config has `concurrency: 50` and command-line has `--concurrency 100`
- **THEN** the agent uses `100` from command-line

## ADDED Requirements

### Requirement: Command-Line Ignore Flag

The agent SHALL support an `--ignore` flag for push and pull commands that accepts ignore patterns. The `--ignore` flag SHALL accept comma-separated values or multiple flags. Ignore patterns from command-line SHALL be merged with ignore patterns from global and local config files, with command-line patterns taking precedence for duplicates.

#### Scenario: Ignore flag with comma-separated values
- **WHEN** `kkartifact-agent push --project myproject --app myapp --version v1.0.0 --ignore "logs/, tmp/, *.log"`
- **THEN** the agent merges these ignore patterns with config file ignore patterns
- **AND** uses the combined ignore list during manifest generation

#### Scenario: Ignore flag with multiple flags
- **WHEN** `kkartifact-agent push --project myproject --app myapp --version v1.0.0 --ignore logs/ --ignore tmp/ --ignore "*.log"`
- **THEN** the agent combines all `--ignore` flags into a single ignore list
- **AND** merges with config file ignore patterns

#### Scenario: Ignore pattern merging with global config
- **WHEN** global config has `ignore: ["logs/", "tmp/"]` and command-line has `--ignore "test.rb, *.log"`
- **THEN** the final ignore list contains `["logs/", "tmp/", "test.rb", "*.log"]`
- **AND** all patterns are applied during push/pull operations

#### Scenario: Ignore pattern override for duplicates
- **WHEN** global config has `ignore: ["logs/"]` and command-line has `--ignore "logs/, tmp/"`
- **THEN** the final ignore list contains `["logs/", "tmp/"]` (duplicate `logs/` removed, command-line version kept)
- **AND** the order is preserved: global patterns first, then command-line patterns

#### Scenario: Ignore pattern merging with local config
- **WHEN** global config has `ignore: ["logs/"]`, local config has `ignore: ["tmp/"]`, and command-line has `--ignore "*.log"`
- **THEN** the final ignore list contains `["logs/", "tmp/", "*.log"]`
- **AND** all patterns from all sources are combined

### Requirement: Command-Line Configuration Overrides

The agent SHALL support command-line flags to override all configuration file values, including `--server-url`, `--token`, `--concurrency`, `--ignore`, and other configuration fields. Command-line flags SHALL take precedence over both global and local configuration files.

#### Scenario: Override server URL via command-line
- **WHEN** config file has `server_url: https://packages.example.com/` and command-line has `--server-url https://other.example.com/`
- **THEN** the agent uses `https://other.example.com/` for API requests

#### Scenario: Override token via command-line
- **WHEN** config file has `token: abc123` and command-line has `--token xyz789`
- **THEN** the agent uses `xyz789` for authentication

#### Scenario: Override concurrency via command-line
- **WHEN** config file has `concurrency: 50` and command-line has `--concurrency 100`
- **THEN** the agent uses 100 concurrent workers for uploads/downloads

#### Scenario: Multiple overrides in single command
- **WHEN** `kkartifact-agent push --project myproject --app myapp --version v1.0.0 --server-url https://other.example.com/ --token xyz789 --concurrency 100 --ignore "logs/"`
- **THEN** all specified flags override corresponding config file values
- **AND** the agent uses the overridden values for the operation
