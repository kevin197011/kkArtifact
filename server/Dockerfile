# Copyright (c) 2025 kk
#
# This software is released under the MIT License.
# https://opensource.org/licenses/MIT

# Build stage
FROM golang:1.24-alpine AS builder

# Install git for version detection and jq for JSON generation
RUN apk --no-cache add git jq

WORKDIR /build

# Copy server go mod files
COPY server/go.mod server/go.sum ./
RUN go mod download

# Copy agent go mod files
COPY agent/go.mod agent/go.sum ./agent/
RUN cd agent && go mod download

# Copy source code
COPY server/ ./server/
COPY agent/ ./agent/

# Build the server application
RUN cd server && CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ../kkartifact-server ./main.go

# Verify server binary exists
RUN ls -la kkartifact-server

# Build agent binaries for all platforms
# Create output directory
RUN mkdir -p server/static/agent server/static/scripts

# Get version info (try git tag, fallback to timestamp)
RUN VERSION=$(git describe --tags --exact-match 2>/dev/null || git describe --tags 2>/dev/null || echo "v$(date +%Y%m%d%H%M%S)") && \
    BUILD_TIME=$(date -u +%Y-%m-%dT%H:%M:%SZ) && \
    GIT_COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown") && \
    echo "Building agent binaries with version: $VERSION" && \
    \
    # Build linux/amd64 \
    cd agent && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath \
      -ldflags="-s -w -X github.com/kk/kkartifact-agent/internal/cli.Version=$VERSION -X github.com/kk/kkartifact-agent/internal/cli.BuildTime=$BUILD_TIME -X github.com/kk/kkartifact-agent/internal/cli.GitCommit=$GIT_COMMIT" \
      -o ../server/static/agent/kkartifact-agent-linux-amd64 ./main.go && \
    \
    # Build linux/arm64 \
    CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -trimpath \
      -ldflags="-s -w -X github.com/kk/kkartifact-agent/internal/cli.Version=$VERSION -X github.com/kk/kkartifact-agent/internal/cli.BuildTime=$BUILD_TIME -X github.com/kk/kkartifact-agent/internal/cli.GitCommit=$GIT_COMMIT" \
      -o ../server/static/agent/kkartifact-agent-linux-arm64 ./main.go && \
    \
    # Build darwin/amd64 \
    CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -trimpath \
      -ldflags="-s -w -X github.com/kk/kkartifact-agent/internal/cli.Version=$VERSION -X github.com/kk/kkartifact-agent/internal/cli.BuildTime=$BUILD_TIME -X github.com/kk/kkartifact-agent/internal/cli.GitCommit=$GIT_COMMIT" \
      -o ../server/static/agent/kkartifact-agent-darwin-amd64 ./main.go && \
    \
    # Build darwin/arm64 \
    CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -trimpath \
      -ldflags="-s -w -X github.com/kk/kkartifact-agent/internal/cli.Version=$VERSION -X github.com/kk/kkartifact-agent/internal/cli.BuildTime=$BUILD_TIME -X github.com/kk/kkartifact-agent/internal/cli.GitCommit=$GIT_COMMIT" \
      -o ../server/static/agent/kkartifact-agent-darwin-arm64 ./main.go && \
    \
    # Build windows/amd64 \
    CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -trimpath \
      -ldflags="-s -w -X github.com/kk/kkartifact-agent/internal/cli.Version=$VERSION -X github.com/kk/kkartifact-agent/internal/cli.BuildTime=$BUILD_TIME -X github.com/kk/kkartifact-agent/internal/cli.GitCommit=$GIT_COMMIT" \
      -o ../server/static/agent/kkartifact-agent-windows-amd64.exe ./main.go && \
    \
    cd .. && \
    echo "Agent binaries built successfully"

# Generate version.json
RUN VERSION=$(git describe --tags --exact-match 2>/dev/null || git describe --tags 2>/dev/null || echo "v$(date +%Y%m%d%H%M%S)") && \
    BUILD_TIME=$(date -u +%Y-%m-%dT%H:%M:%SZ) && \
    SIZE_LINUX_AMD64=$(stat -c%s server/static/agent/kkartifact-agent-linux-amd64 2>/dev/null || echo 0) && \
    SIZE_LINUX_ARM64=$(stat -c%s server/static/agent/kkartifact-agent-linux-arm64 2>/dev/null || echo 0) && \
    SIZE_DARWIN_AMD64=$(stat -c%s server/static/agent/kkartifact-agent-darwin-amd64 2>/dev/null || echo 0) && \
    SIZE_DARWIN_ARM64=$(stat -c%s server/static/agent/kkartifact-agent-darwin-arm64 2>/dev/null || echo 0) && \
    SIZE_WINDOWS_AMD64=$(stat -c%s server/static/agent/kkartifact-agent-windows-amd64.exe 2>/dev/null || echo 0) && \
    jq -n --arg version "$VERSION" --arg build_time "$BUILD_TIME" --argjson size_linux_amd64 "$SIZE_LINUX_AMD64" --argjson size_linux_arm64 "$SIZE_LINUX_ARM64" --argjson size_darwin_amd64 "$SIZE_DARWIN_AMD64" --argjson size_darwin_arm64 "$SIZE_DARWIN_ARM64" --argjson size_windows_amd64 "$SIZE_WINDOWS_AMD64" '{version: $version, build_time: $build_time, binaries: [{platform: "linux/amd64", filename: "kkartifact-agent-linux-amd64", size: $size_linux_amd64, url: "/api/v1/downloads/agent/kkartifact-agent-linux-amd64"}, {platform: "linux/arm64", filename: "kkartifact-agent-linux-arm64", size: $size_linux_arm64, url: "/api/v1/downloads/agent/kkartifact-agent-linux-arm64"}, {platform: "darwin/amd64", filename: "kkartifact-agent-darwin-amd64", size: $size_darwin_amd64, url: "/api/v1/downloads/agent/kkartifact-agent-darwin-amd64"}, {platform: "darwin/arm64", filename: "kkartifact-agent-darwin-arm64", size: $size_darwin_arm64, url: "/api/v1/downloads/agent/kkartifact-agent-darwin-arm64"}, {platform: "windows/amd64", filename: "kkartifact-agent-windows-amd64.exe", size: $size_windows_amd64, url: "/api/v1/downloads/agent/kkartifact-agent-windows-amd64.exe"}]}' > server/static/agent/version.json

# Copy install scripts
COPY scripts/install-agent.sh scripts/install-agent.ps1 server/static/scripts/

# Verify agent binaries exist
RUN ls -la server/static/agent/ && \
    test -f server/static/agent/kkartifact-agent-linux-amd64 && \
    test -f server/static/agent/kkartifact-agent-linux-arm64 && \
    test -f server/static/agent/kkartifact-agent-darwin-amd64 && \
    test -f server/static/agent/kkartifact-agent-darwin-arm64 && \
    test -f server/static/agent/kkartifact-agent-windows-amd64.exe && \
    test -f server/static/agent/version.json && \
    test -f server/static/scripts/install-agent.sh && \
    test -f server/static/scripts/install-agent.ps1 && \
    echo "All agent binaries, scripts, and version.json created successfully"

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata wget

WORKDIR /app

# Copy binary from builder
COPY --from=builder /build/kkartifact-server .
COPY --from=builder /build/server/migrations ./migrations
# Copy static files (agent binaries and scripts)
COPY --from=builder /build/server/static ./static

# Verify binary exists in final image
RUN ls -la kkartifact-server && chmod +x kkartifact-server

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/api/v1/health || exit 1

# Run the application
CMD ["./kkartifact-server"]
