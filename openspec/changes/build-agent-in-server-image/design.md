# Design: Build Agent Binaries in Server Docker Image

## Context

The kkArtifact system consists of:
- **Server** (`server/`): Go application providing HTTP API
- **Agent** (`agent/`): Go CLI tool for push/pull operations
- **Web UI** (`web-ui/`): React frontend

Currently, agent binaries are built separately using `scripts/build-agent-binaries.rb` and stored in `server/static/agent/` before Docker builds. The server Dockerfile simply copies these pre-built binaries.

## Goals

1. Integrate agent binary building into server Dockerfile
2. Build agent binaries for all supported architectures during Docker image creation
3. Generate `version.json` with binary metadata during build
4. Maintain backward compatibility with existing download API

## Non-Goals

- Building agent binaries for architectures not supported by the base image
- Cross-compilation for Windows from Linux (handled separately if needed)
- Modifying agent build logic (reuse existing build flags and version injection)

## Decisions

### Decision 1: Multi-Stage Build with Agent Compilation

**What**: Use Docker multi-stage build to compile agent binaries in the builder stage.

**Why**: 
- Keeps build dependencies (Go toolchain) out of final image
- Allows building multiple architectures in parallel
- Maintains clean separation between build and runtime

**Alternatives considered**:
- Single-stage build: Increases final image size unnecessarily
- Separate Dockerfile for agent: Adds complexity, requires coordination

### Decision 2: Build All Architectures in Builder Stage

**What**: Build agent binaries for all 5 platforms (linux/amd64, linux/arm64, darwin/amd64, darwin/arm64, windows/amd64) in the Dockerfile builder stage.

**Why**:
- Ensures all binaries are available in the final image
- Matches current behavior of `build-agent-binaries.rb`
- Supports all deployment scenarios

**Alternatives considered**:
- Build only linux/amd64: Would break darwin/arm64 and windows support
- Conditional builds based on TARGETPLATFORM: Would require buildx and add complexity

### Decision 3: Reuse Existing Build Logic

**What**: Use the same build flags and version injection logic as `scripts/build-agent-binaries.rb`.

**Why**:
- Consistency with existing build process
- Proven build configuration
- Version information properly embedded

**Implementation**:
- Use `go build` with `-trimpath` and `-ldflags` for version injection
- Set `GOOS` and `GOARCH` environment variables for cross-compilation
- Use `CGO_ENABLED=0` for static binaries

### Decision 4: Generate version.json in Dockerfile

**What**: Generate `server/static/agent/version.json` during Docker build using a simple script or inline commands.

**Why**:
- Ensures version.json matches built binaries
- Provides metadata for download API
- Maintains compatibility with existing API endpoints

**Implementation**:
- Use `jq` or inline JSON generation in Dockerfile
- Extract version from git tags or use build-time variables
- Include all built binaries in version.json

## Risks / Trade-offs

### Risk 1: Build Time Increase
**Risk**: Building 5 architectures increases Docker build time.

**Mitigation**: 
- Multi-stage builds allow layer caching
- Only rebuild agent binaries when agent code changes
- Consider parallel builds if build time becomes an issue

### Risk 2: Build Context Size
**Risk**: Including agent source code increases Docker build context.

**Mitigation**:
- Agent code is already in repository
- Use `.dockerignore` to exclude unnecessary files
- Build context size increase is minimal (~few MB)

### Risk 3: Cross-Compilation Complexity
**Risk**: Building darwin/windows binaries from linux base image may have issues.

**Mitigation**:
- Go's built-in cross-compilation is well-tested
- Test all architectures in CI
- Fallback to pre-built binaries if needed (can be added later)

## Implementation Plan

1. **Modify Dockerfile builder stage**:
   - Add agent source code copy
   - Build agent binaries for all platforms
   - Generate version.json

2. **Update .dockerignore** (if needed):
   - Ensure agent source is included
   - Exclude unnecessary files

3. **Test build process**:
   - Verify all binaries are created
   - Verify version.json is correct
   - Test download API endpoints

4. **Update documentation**:
   - Remove manual build step instructions
   - Update CI/CD workflows if needed

## Open Questions

1. Should we support building only specific architectures via build args? (Not needed initially)
2. Should we add build caching for agent binaries? (Can be optimized later)
3. Should we validate binary checksums? (Can be added later if needed)
