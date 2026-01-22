# Change: Build Agent Binaries in Server Docker Image

## Why

Currently, agent binaries are pre-built and stored in `server/static/agent/` directory before building the server Docker image. This requires a separate build step (`scripts/build-agent-binaries.rb`) to be run manually or in CI before Docker builds, which:

1. **Adds complexity**: Requires developers to remember to build agent binaries before building Docker images
2. **Breaks automation**: Docker builds fail if agent binaries are missing or outdated
3. **Increases maintenance**: Two separate build processes need to be kept in sync
4. **Reduces portability**: Pre-built binaries must be committed to the repository or managed separately

By integrating agent binary building into the server Dockerfile, we ensure that:
- Agent binaries are always built for all supported architectures during Docker image creation
- No manual pre-build steps are required
- The Docker build process is self-contained and reproducible
- Agent binaries are guaranteed to match the server version

## What Changes

- **MODIFIED**: `server/Dockerfile` - Add multi-architecture agent binary building in the builder stage
- **MODIFIED**: Build process - Agent binaries are now built during Docker image creation instead of pre-built
- **ADDED**: Build script integration - The Dockerfile will execute agent builds for all platforms (linux/amd64, linux/arm64, darwin/amd64, darwin/arm64, windows/amd64)
- **ADDED**: Version info generation - Generate `version.json` during Docker build

## Impact

- **Affected specs**: Build system (new capability)
- **Affected code**: 
  - `server/Dockerfile` - Multi-stage build with agent compilation
  - Build process documentation
- **Breaking changes**: None - this is an internal build process change
- **Migration**: Existing pre-built binaries in `server/static/agent/` can be removed from version control (added to `.gitignore`)

## Benefits

1. **Simplified workflow**: Single `docker build` command builds everything
2. **Consistency**: Agent binaries always match server version
3. **Reproducibility**: Docker builds are fully self-contained
4. **CI/CD friendly**: No separate build steps required in CI pipelines
