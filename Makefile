# Copyright (c) 2025 kk
#
# This software is released under the MIT License.
# https://opensource.org/licenses/MIT

.PHONY: build-server build-agent build-all test clean docker-up docker-down docker-build help

# Build server
build-server:
	cd server && go build -o ../bin/kkartifact-server ./main.go

# Build agent
build-agent:
	cd agent && go build -o ../bin/kkartifact-agent ./main.go

# Build agent for all platforms
build-agent-all:
	ruby scripts/build-agent-binaries.rb

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

