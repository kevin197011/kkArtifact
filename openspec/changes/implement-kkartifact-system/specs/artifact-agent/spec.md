## ADDED Requirements

### Requirement: Configuration File Requirement
The agent SHALL require a `.kkartifact.yml` configuration file to exist before executing push or pull operations. If the file is missing, the agent SHALL exit with an error.

#### Scenario: Push without config file
- **WHEN** push command is executed without `.kkartifact.yml` in the current directory
- **THEN** the agent exits with an error message indicating the config file is required

#### Scenario: Pull without config file
- **WHEN** pull command is executed without `.kkartifact.yml` in the deployment directory
- **THEN** the agent exits with an error message indicating the config file is required

### Requirement: Ignore Rules Processing
The agent SHALL apply ignore rules from `.kkartifact.yml` during push and pull operations. Ignore rules SHALL support glob patterns, directory prefixes, and file patterns.

#### Scenario: Ignore files during push
- **WHEN** `.kkartifact.yml` contains `ignore: ["logs/", "*.log"]`
- **THEN** files matching these patterns are excluded from push operations

#### Scenario: Ignore files during pull
- **WHEN** `.kkartifact.yml` contains `ignore: ["tmp/", "*.tmp"]`
- **THEN** files matching these patterns are excluded from pull operations

### Requirement: Push Operation
The agent SHALL support pushing artifacts from a local build directory to the server. The agent SHALL accept a version parameter, scan the local directory, generate a manifest, and upload files to the server.

#### Scenario: Push artifact with version
- **WHEN** `kkArtifact-agent push --project myproject --app myapp --version 1.0.0 --path /local/build/path`
- **THEN** the agent scans files in `/local/build/path`
- **AND** applies ignore rules from `.kkartifact.yml`
- **AND** calculates SHA256 for each file
- **AND** generates a manifest with metadata
- **AND** uploads all files to the server
- **AND** finalizes the upload to create the version

#### Scenario: Push with concurrent uploads
- **WHEN** pushing multiple files
- **THEN** files are uploaded concurrently using a configurable worker pool (default 8 workers)

#### Scenario: Push with retry on failure
- **WHEN** a file upload fails due to network error
- **THEN** the agent retries the upload with exponential backoff
- **AND** the operation is idempotent (can be safely retried)

### Requirement: Pull Operation
The agent SHALL support pulling artifacts from the server to a local deployment directory. The agent SHALL download only missing or changed files based on manifest comparison, and SHALL delete files that no longer exist in the new version.

#### Scenario: Pull artifact version
- **WHEN** `kkArtifact-agent pull --project myproject --app myapp --version 1.0.0 --deploy-path /opt/apps/myproject/myapp`
- **THEN** the agent fetches the manifest for the specified version
- **AND** compares local files with manifest
- **AND** downloads only missing or changed files
- **AND** applies ignore rules from `.kkartifact.yml`
- **AND** files are stored directly in `/opt/apps/myproject/myapp/` (no soft links)

#### Scenario: Pull with file cleanup
- **WHEN** pulling a new version that no longer includes a file that existed in the previous version
- **THEN** the agent deletes the file from the deployment directory
- **AND** only files present in the new version manifest remain in the deployment directory

#### Scenario: Pull with concurrent downloads
- **WHEN** pulling multiple files
- **THEN** files are downloaded concurrently using a configurable worker pool (default 8 workers)

#### Scenario: Pull with HTTP Range support
- **WHEN** downloading a large file that was partially downloaded previously
- **THEN** the agent uses HTTP Range requests to resume the download

### Requirement: Manifest Generation
The agent SHALL generate a manifest (`meta.yaml`) during push operations that includes project, app, version, git commit (if available), build time, builder, and file list with SHA256 checksums and sizes.

#### Scenario: Generate manifest during push
- **WHEN** pushing files, the agent generates a manifest
- **THEN** the manifest includes project, app, version from command-line parameters
- **AND** build_time is set to current timestamp
- **AND** builder is set from environment variable or hostname
- **AND** git_commit is extracted from git repository if available
- **AND** file list includes all files with their paths, SHA256, and sizes

### Requirement: Manifest Diff
The agent SHALL compare local files with a remote manifest to determine which files need to be downloaded during pull operations.

#### Scenario: Diff manifest for pull
- **WHEN** pulling a version, the agent fetches the remote manifest
- **THEN** the agent compares each file in the manifest with local files
- **AND** files are identified as missing, changed (different SHA256), or unchanged
- **AND** only missing or changed files are downloaded

### Requirement: Local Directory Structure
The agent SHALL handle different directory structures for push (build directory) and pull (deployment directory). Push source files are in `/local/build/path/`, and pull destination files are directly in `/opt/apps/{project}/{app}/`.

#### Scenario: Push from build directory
- **WHEN** pushing from `/local/build/path/`
- **THEN** all files in that directory (after ignore rules) are included in the artifact
- **AND** relative paths are preserved (e.g., `bin/app` stays as `bin/app`)

#### Scenario: Pull to deployment directory
- **WHEN** pulling to `/opt/apps/myproject/myapp`
- **THEN** files are placed directly in that directory
- **AND** no soft links are created
- **AND** the directory structure from the artifact is preserved

### Requirement: Agent Configuration
The agent SHALL support configuration via `.kkartifact.yml` including server URL, token, concurrency level, chunk size, ignore rules, and optional version retention limit.

#### Scenario: Load configuration
- **WHEN** the agent starts, it loads `.kkartifact.yml`
- **THEN** server URL, token, concurrency, chunk_size, ignore rules, and optional retain_versions are read
- **AND** command-line parameters override config file values

### Requirement: App-Level Version Retention Control
The agent SHALL check version count for the current app after push operations. If the number of versions exceeds the global retention limit configured on the server, the agent SHALL trigger deletion of the oldest versions for that specific app only, without affecting other apps.

#### Scenario: Check version count after push
- **WHEN** a push operation completes successfully for project `myproject` and app `myapp`
- **THEN** the agent checks the number of versions for `myproject/myapp`
- **AND** compares it with the global retention limit from server configuration

#### Scenario: Delete old versions for current app
- **WHEN** after push, the app has 8 versions and the global retention limit is 5
- **THEN** the agent triggers deletion of the 3 oldest versions for that app only
- **AND** other apps are not affected
- **AND** only versions for the current app (myproject/myapp) are considered

#### Scenario: Version retention is app-specific
- **WHEN** version retention cleanup is triggered for app `myapp1`
- **THEN** only versions belonging to `myapp1` are considered for deletion
- **AND** versions for `myapp2` in the same project are unaffected
- **AND** versions for apps in other projects are unaffected

