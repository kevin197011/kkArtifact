# artifact-agent Specification

## Purpose
TBD - created by archiving change add-g02-batch-migration-workflow. Update Purpose after archive.
## Requirements
### Requirement: Push Path Inference
The agent SHALL infer app and version from the last two path segments of `--path` when `--app` and `--version` are omitted.

#### Scenario: Push with path inference
- **WHEN** `kkartifact-agent push --project g02 --path /data/vcs/G02/tidb/G02_agent_api/cd884eb...`
- **THEN** app is inferred as `G02_agent_api` and version as `cd884eb...`
- **AND** files under that directory are pushed to the server

### Requirement: Push-tree Batch Command
The agent SHALL provide a `push-tree` command that traverses `root/{app}/{version}/` and pushes each version directory.

#### Scenario: Batch push directory tree
- **WHEN** `kkartifact-agent push-tree /data/vcs/G02/tidb --project g02`
- **THEN** each `{app}/{version}` subdirectory is pushed sequentially
- **AND** `--skip-existing` skips versions already on the server
- **AND** `--dry-run` prints planned pushes without uploading

### Requirement: Push Publish Flag
The agent SHALL support a `--publish` flag on `push` and `push-tree` to mark the version as published after a successful push.

#### Scenario: Auto-publish after push
- **WHEN** push completes with `--publish` and the token has `promote` permission
- **THEN** the agent calls the publish API for that version
- **AND** the version becomes available via `pull --version latest`

### Requirement: Push Skip Unchanged
The agent SHALL skip uploading when the remote version manifest matches the local manifest, and skip individual files whose hashes match the remote manifest.

#### Scenario: Skip fully unchanged version
- **WHEN** the remote version already exists with an identical manifest
- **THEN** the agent skips upload and reports the version as up to date

#### Scenario: Skip unchanged files within version
- **WHEN** only some files differ from the remote manifest
- **THEN** the agent uploads only changed files

