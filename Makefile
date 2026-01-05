# Copyright (c) 2025 kk
#
# This software is released under the MIT License.
# https://opensource.org/licenses/MIT

.PHONY: build-server build-agent build-all test clean docker-up docker-down docker-build help

# Build server
build-server:
	@VERSION=$$(git describe --tags --exact-match 2>/dev/null || git describe --tags 2>/dev/null || echo "dev"); \
	BUILD_TIME=$$(date -u +%Y-%m-%dT%H:%M:%SZ); \
	GIT_COMMIT=$$(git rev-parse --short HEAD 2>/dev/null || echo "unknown"); \
	cd server && go build -ldflags "-X github.com/kk/kkartifact-server/internal/version.Version=$$VERSION -X github.com/kk/kkartifact-server/internal/version.BuildTime=$$BUILD_TIME -X github.com/kk/kkartifact-server/internal/version.GitCommit=$$GIT_COMMIT" -o ../bin/kkartifact-server ./main.go

# Build agent
build-agent:
	@VERSION=$$(git describe --tags --exact-match 2>/dev/null || git describe --tags 2>/dev/null || echo "dev"); \
	BUILD_TIME=$$(date -u +%Y-%m-%dT%H:%M:%SZ); \
	GIT_COMMIT=$$(git rev-parse --short HEAD 2>/dev/null || echo "unknown"); \
	cd agent && go build -ldflags "-X github.com/kk/kkartifact-agent/internal/cli.Version=$$VERSION -X github.com/kk/kkartifact-agent/internal/cli.BuildTime=$$BUILD_TIME -X github.com/kk/kkartifact-agent/internal/cli.GitCommit=$$GIT_COMMIT" -o ../bin/kkartifact-agent ./main.go

# Build agent for all platforms
build-agent-all:
	ruby scripts/build-agent-binaries.rb

# Update agent version.json only (without rebuilding binaries)
update-agent-version:
	ruby scripts/update-agent-version.rb

# Build all
build-all: build-server build-agent

# Test
test:
	cd server && go test ./...
	cd agent && go test ./...

# Clean
clean:
	rm -rf bin/
	cd server && go clean
	cd agent && go clean

# Docker commands
docker-up:
	docker compose up -d

docker-down:
	docker compose down

docker-build:
	docker compose build

docker-logs:
	docker compose logs -f

# Help
help:
	@echo "Available targets:"
	@echo "  build-server    - Build kkArtifact server"
	@echo "  build-agent     - Build kkArtifact agent"
	@echo "  build-all       - Build all components"
	@echo "  test            - Run tests"
	@echo "  clean           - Clean build artifacts"
	@echo "  docker-up       - Start Docker Compose services"
	@echo "  docker-down     - Stop Docker Compose services"
	@echo "  docker-build    - Build Docker images"
	@echo "  docker-logs     - View Docker Compose logs"

