# Project Context

## Purpose

kkArtifact is a modern artifact management and synchronization system designed to replace rsync + SSH deployment workflows. The system provides:

- Multi-project and multi-app artifact storage with hash-based versioning
- HTTP-based API for artifact upload, download, and management
- Token-based authentication with fine-grained permissions (Global/Project/App scopes)
- CLI agent for push/pull operations with manifest generation
- Web UI for visual management and monitoring
- Event-driven architecture with webhook support for integration

## Tech Stack

### Server (kkArtifact-server)
- **Language**: Go 1.21+
- **Web Framework**: Standard library `net/http` or lightweight framework (gin/fiber) if needed
- **Storage**: Local filesystem (extendable to NFS/object storage)
- **Database**: PostgreSQL for metadata (tokens, webhooks, audit logs)
- **API Format**: RESTful HTTP API with JSON responses

### Agent (kkArtifact-agent)
- **Language**: Go 1.21+ (shared codebase with server)
- **CLI Framework**: `cobra` for command-line interface
- **Configuration**: YAML via `.kkartifact.yml` (using `gopkg.in/yaml.v3`)

### Web UI
- **Frontend**: TypeScript + React 18+
- **Build Tool**: Vite
- **State Management**: React Query
- **UI Framework**: Ant Design (antd) - enterprise-grade React UI library
- **Backend API**: Uses kkArtifact-server API endpoints

### Scripts and Tooling
- **Script Language**: Ruby 3.1+ (per project conventions)
- **Task Management**: Rakefile for project tasks

## Project Conventions

### Code Style
- Follow Go standard formatting (`gofmt`)
- Use `golangci-lint` for Go code quality
- Use ESLint and Prettier for TypeScript/React code
- All new files must include MIT license header with copyright (c) 2025 kk

### Architecture Patterns
- **Immutable Versions**: Hash-based versions cannot be modified once created
- **Stateless Server**: Server maintains no session state, all auth via tokens
- **Event-Driven**: Major operations emit events that trigger webhooks
- **Manifest-Based**: All versions include `meta.yaml` with file metadata
- **Idempotent Operations**: Push/pull operations can be safely retried

### Configuration
- Configuration file: `.kkartifact.yml` (YAML format)
- Configuration is **required** for push/pull operations (agent exits if missing)
- Ignore rules support glob patterns, directory prefixes, and file patterns
- Optional `retain_versions` in config file (server uses global setting, agent can reference it)
- Server maintains global version retention limit (applies to all apps)

### Testing Strategy
- Unit tests for core business logic
- Integration tests for API endpoints
- Integration tests for agent operations
- End-to-end tests for critical workflows
- Code coverage target: ≥80% for core logic, 100% for critical paths

### Git Workflow
- Follow Conventional Commits format
- Use meaningful commit messages with type and scope
- Main branch is `main`
- Feature branches for development

## Domain Context

### Key Concepts

**Project**: Top-level namespace for organizing artifacts (e.g., `project-a`)

**App**: Application within a project (e.g., `app-api`, `app-worker`)

**Version/Hash**: Immutable artifact version identified by a hash (e.g., `a8f3c21d`)

**Manifest (meta.yaml)**: Metadata file containing project, app, version, file list with SHA256 checksums

**Token Scope**: Permission boundary (Global = all projects/apps, Project = all apps in a project, App = single app)

**Event**: System event (push, pull, promote, rollback, delete) that can trigger webhooks

### Directory Structure

**Server Storage**:
```
/repos
└── {project}
    └── {app}
        └── {hash}/
            ├── bin/
            ├── config/
            └── meta.yaml
```

**Agent Push Source**: `/local/build/path/` (build directory)

**Agent Pull Destination**: `/opt/apps/{project}/{app}/` (deployment directory, no soft links)

## Important Constraints

1. **Version Immutability**: Versions (hashes) cannot be modified or overwritten once created (but can be deleted for cleanup)
2. **No Soft Links**: Deployment uses direct file operations, not symbolic links
3. **Configuration Required**: `.kkartifact.yml` must exist for push/pull operations
4. **File Cleanup**: Pull operations delete files that no longer exist in the new version
5. **Large File Support**: System must handle files >1GB with HTTP Range support
6. **Concurrent Operations**: Use worker pools for file uploads/downloads (default 8 workers)
7. **Global Version Retention**: Server maintains a global retention limit (e.g., keep latest 5 versions per app)
8. **Scheduled Cleanup**: Server runs daily cleanup task at 3:00 AM to remove old versions
9. **App-Specific Cleanup**: Version cleanup applies per app independently, does not affect other apps

## External Dependencies

- **Storage Backend**: Object storage (S3/OSS compatible) as primary backend, local filesystem for development
- **Database**: PostgreSQL for metadata storage (essential for 2000+ apps)
- **Cache**: Redis for caching热点数据 (project/app lists, latest versions, manifest metadata)
- **Authentication**: Token-based (no external auth service required)
- **Webhooks**: External HTTP endpoints (Slack, CI/CD, internal services)

## Large-Scale Support

### Performance Targets
- **Module Scale**: Support 2000+ App modules
- **Storage Capacity**: Support 2TB+ storage capacity
- **Concurrent Access**: Support 200+ concurrent operations
- **Query Performance**:
  - Project list query: <500ms (2000+ projects)
  - App list query: <500ms (100+ apps per project)
  - Version list query: <1s (100+ versions per app)
- **Storage Performance**:
  - File upload: 100+ MB/s throughput
  - File download: 100+ MB/s throughput

### Optimization Strategies
- **Object Storage**: S3/OSS for distributed storage and high I/O performance
- **Redis Caching**: Cache frequently accessed data to reduce database load
- **Database Indexing**: Optimized indexes on key fields for fast queries
- **API Pagination**: All list APIs support pagination (default 50, max 500 per page)
- **Response Compression**: Gzip compression for all API responses
- **Parallel Processing**: Parallel SHA256 computation and async operations for large-scale tasks

## Development Environment

### Docker Compose Setup
The project provides a Docker Compose configuration for local development and testing:
- **Services**: kkArtifact-server, PostgreSQL, Redis, Web UI
- **Configuration**: Environment variables via `.env` file
- **Data Persistence**: Named volumes for database, cache, and artifact storage
- **Hot Reload**: Development mode with code hot-reload support
- **Quick Start**: `docker compose up -d` to start all services

### Prerequisites
- Docker 20.10+
- Docker Compose 2.0+
- (Optional) Local Go 1.21+ for agent development
- (Optional) Node.js 18+ for Web UI development
