# Build System Specification

## ADDED Requirements

### Requirement: Agent Binary Building in Server Docker Image

The server Docker image build process SHALL compile agent binaries for all supported architectures during the Docker build, eliminating the need for pre-built binaries.

#### Scenario: Docker build creates all agent binaries
- **WHEN** building the server Docker image
- **THEN** agent binaries are compiled for linux/amd64, linux/arm64, darwin/amd64, darwin/arm64, and windows/amd64
- **AND** all binaries are placed in `server/static/agent/` directory in the final image
- **AND** `version.json` is generated with metadata for all built binaries

#### Scenario: Agent binaries match server version
- **WHEN** agent binaries are built during Docker image creation
- **THEN** version information embedded in binaries matches the server version
- **AND** build time and git commit information are correctly injected

#### Scenario: No pre-build step required
- **WHEN** building the server Docker image from source
- **THEN** no manual pre-build step is required
- **AND** agent binaries are automatically built as part of the Docker build process

### Requirement: Multi-Architecture Agent Binary Support

The server Docker image SHALL contain agent binaries for all supported platforms, enabling users to download the appropriate binary for their system.

#### Scenario: All architectures available
- **WHEN** the server Docker image is built
- **THEN** agent binaries are available for linux/amd64, linux/arm64, darwin/amd64, darwin/arm64, and windows/amd64
- **AND** each binary is correctly named with platform suffix (e.g., `kkartifact-agent-linux-amd64`, `kkartifact-agent-windows-amd64.exe`)

#### Scenario: Version metadata available
- **WHEN** agent binaries are built
- **THEN** `version.json` is generated in `server/static/agent/` directory
- **AND** version.json contains version, build_time, and list of all binaries with platform, filename, size, and URL
