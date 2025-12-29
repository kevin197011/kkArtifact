## ADDED Requirements

### Requirement: Artifact Version Storage
The system SHALL store artifacts in an immutable, hash-based version structure. Each version SHALL be uniquely identified by a hash and SHALL be stored at `/repos/{project}/{app}/{hash}/`.

#### Scenario: Store new artifact version
- **WHEN** an artifact is pushed with project `myproject`, app `myapp`, and hash `a8f3c21d`
- **THEN** files are stored at `/repos/myproject/myapp/a8f3c21d/`
- **AND** the version cannot be modified or overwritten once created

#### Scenario: Prevent version overwrite
- **WHEN** attempting to push an artifact with an existing hash
- **THEN** the operation fails with an error indicating the version already exists

### Requirement: Manifest File Generation
The system SHALL automatically generate a `meta.yaml` manifest file for each artifact version containing project, app, version, git commit, build time, builder, and file list with SHA256 checksums and sizes.

#### Scenario: Generate manifest on push
- **WHEN** an artifact is pushed with version `a8f3c21d` and files `bin/app` (size 123456, SHA256 `abc123...`)
- **THEN** a `meta.yaml` file is created at `/repos/{project}/{app}/a8f3c21d/meta.yaml`
- **AND** the manifest contains project, app, version, build_time, builder, and all files with their SHA256 and size

### Requirement: File Directory Structure
The system SHALL preserve the directory structure of uploaded files within each version directory. Files SHALL be organized in subdirectories (e.g., `bin/`, `config/`) as provided during upload.

#### Scenario: Preserve directory structure
- **WHEN** files `bin/app` and `config/app.yml` are uploaded
- **THEN** they are stored at `{hash}/bin/app` and `{hash}/config/app.yml` respectively

### Requirement: File Integrity Verification
The system SHALL verify file integrity using SHA256 checksums. Files SHALL be verified on upload and their checksums SHALL be stored in the manifest.

#### Scenario: Verify file checksum on upload
- **WHEN** a file is uploaded with a provided SHA256 checksum
- **THEN** the system verifies the uploaded file matches the provided checksum
- **AND** if checksums do not match, the upload fails

### Requirement: Version Query
The system SHALL support querying versions by project and app, returning versions ordered by creation time (newest first).

#### Scenario: List versions for app
- **WHEN** querying versions for project `myproject` and app `myapp`
- **THEN** all version hashes for that app are returned
- **AND** versions are ordered by creation time descending (newest first)

### Requirement: Global Version Retention Configuration
The system SHALL support a global configuration for the maximum number of versions to retain per app. This configuration SHALL apply to all apps and projects, eliminating the need to configure retention per app.

#### Scenario: Configure global version retention
- **WHEN** the global version retention is set to 5
- **THEN** each app retains a maximum of 5 versions (newest)
- **AND** older versions are eligible for deletion

### Requirement: Scheduled Version Cleanup
The system SHALL execute a scheduled task daily at 3:00 AM to clean up versions that exceed the global retention limit. The cleanup SHALL remove the oldest versions first while preserving the configured number of newest versions.

#### Scenario: Scheduled cleanup execution
- **WHEN** it is 3:00 AM and an app has 10 versions with retention limit of 5
- **THEN** the 5 oldest versions are deleted
- **AND** the 5 newest versions are retained
- **AND** cleanup runs per app independently

#### Scenario: Cleanup preserves newest versions
- **WHEN** scheduled cleanup runs for an app with versions ordered by creation time
- **THEN** the newest N versions (where N is the retention limit) are preserved
- **AND** only versions older than the Nth newest version are deleted

### Requirement: Version Deletion
The system SHALL support deleting version directories and their associated files when versions exceed the retention limit. Deletion SHALL remove all files in the version directory including the manifest.

#### Scenario: Delete old version
- **WHEN** a version exceeds the retention limit and is marked for deletion
- **THEN** the entire version directory `/repos/{project}/{app}/{hash}/` is deleted
- **AND** all files including `meta.yaml` are removed
- **AND** the deletion is logged in the audit log

