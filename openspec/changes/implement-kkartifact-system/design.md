## Context

The kkArtifact system is designed to replace rsync + SSH deployment with a modern, secure, and scalable artifact management system. The system must support:
- Multi-project and multi-app isolation
- Hash-based immutable versioning
- HTTP-based API access
- Token-based authentication with fine-grained permissions
- Web UI for management and monitoring
- Event-driven architecture with webhook support

## Goals / Non-Goals

### Goals
- Complete replacement for rsync-based deployment workflow
- Support for large files (>1GB) with HTTP Range requests
- Concurrent file operations for performance
- Idempotent operations with retry support
- Stateless server architecture (horizontally scalable)
- Version immutability (hash-based versions cannot be modified or overwritten)
- Configuration-driven ignore rules via `.kkartifact.yml`
- No soft links in deployment (direct file operations based on deployment path configuration)
- Docker Compose-based development environment for easy local testing and development
- **Large-scale support: 2000+ apps, 2TB+ storage capacity**
- **Object storage support (S3/OSS) as primary storage backend**
- **Redis caching for performance optimization**
- **Database indexing and query optimization for large datasets**

### Non-Goals (V1)
- Git integration (repo server doesn't care about Git)
- Build/compilation (only artifact storage and distribution)
- K8s/systemd automation (manual deployment via agent)
- Grayscale/batch deployment
- Advanced analytics and reporting
- Migration tools from rsync (out of scope for V1)

## Decisions

### Technology Stack

**Server (kkArtifact-server)**:
- **Language**: Go 1.21+ (chosen for performance, concurrency, and deployment simplicity)
- **Web Framework**: Standard library `net/http` or lightweight framework (gin/fiber) if needed
- **Storage**: Object storage (S3/OSS compatible) as primary backend, local filesystem for development
- **Storage Interface**: Abstract storage interface to support multiple backends (S3, OSS, MinIO, local filesystem)
- **Database**: PostgreSQL for metadata (tokens, webhooks, audit logs)
- **Cache**: Redis for caching热点数据 (project/app lists, latest versions, manifest metadata)
- **API Format**: RESTful HTTP API with JSON responses (Gzip compression enabled)

**Agent (kkArtifact-agent)**:
- **Language**: Go 1.21+ (shared codebase with server for common utilities)
- **CLI Framework**: `cobra` for command-line interface
- **Configuration**: YAML via `.kkartifact.yml` (using `gopkg.in/yaml.v3`)
- **File Operations**: Direct file I/O with SHA256 checksums

**Web UI**:
- **Frontend**: TypeScript + React 18+ (modern, component-based)
- **Build Tool**: Vite (fast development and build)
- **State Management**: React Query (for API state management)
- **UI Framework**: Ant Design (antd) - enterprise-grade React UI library
- **Backend API**: Uses same kkArtifact-server API endpoints

**Configuration**:
- **Format**: YAML (`.kkartifact.yml`)
- **Location**: Repository root (for push) and deployment root (for pull)
- **Required**: Must exist for push/pull operations (error if missing)
- **Agent Config Fields**: server, token, concurrency, chunk_size, ignore, optional retain_versions (for reference)
- **Server Config**: Global version retention limit stored in database (separate from agent config)

### Architecture Patterns

1. **Immutable Versions**: Once a version (hash) is created, it cannot be modified. New deployments = new hash.
2. **Token Scopes**: Hierarchical permission model (Global > Project > App)
3. **Event-Driven**: All major operations emit events that can trigger webhooks
4. **Manifest-Based**: All versions include `meta.yaml` with file manifest and metadata
5. **Idempotent Operations**: Push/pull operations can be retried safely
6. **Stateless Server**: Server maintains no session state, all auth via tokens

### Data Model

**Directory Structure**:
```
/repos
└── {project}
    └── {app}
        └── {hash}/
            ├── bin/
            ├── config/
            └── meta.yaml
```

**meta.yaml Format**:
```yaml
project: project-a
app: app-api
version: a8f3c21d
git_commit: a8f3c21d
build_time: 2025-12-26T17:30:00
builder: build-01
files:
  - path: bin/app
    sha256: xxx
    size: 123456
```

**Token Model**:
- Token ID (unique identifier)
- Token string (hashed for storage)
- Project scope (nullable for global tokens)
- App scope (nullable for project/global tokens)
- Permissions (push, pull, promote, admin)
- Expiration (optional)
- Created timestamp

**Webhook Model**:
- Webhook ID
- Name
- Event types (push, pull, promote, rollback, delete)
- URL (HTTP POST endpoint)
- Headers (optional)
- Enabled/disabled status
- Project/App filter (optional)

**Global Configuration Model**:
- Version retention limit (number of versions to keep per app, global setting)
- Cleanup schedule (cron expression, default: daily at 3:00 AM)
- Configuration updated timestamp

### Security Considerations

1. **Token Storage**: Tokens stored as bcrypt/argon2 hashes (never plaintext)
2. **HTTPS**: All API communication should use HTTPS in production
3. **Input Validation**: All user inputs validated and sanitized
4. **Path Traversal**: Prevent directory traversal attacks in file paths
5. **Rate Limiting**: Implement rate limiting for API endpoints (future enhancement)

### Performance Considerations

1. **Concurrent Operations**: Use worker pools for file uploads/downloads
2. **HTTP Range Support**: Support Range requests for large file downloads
3. **Chunked Uploads**: Support chunked uploads for large files
4. **Connection Pooling**: Reuse HTTP connections in agent
5. **Lazy Loading**: Web UI loads data on-demand, pagination for large lists
6. **Object Storage**: Use S3/OSS for distributed storage, better I/O performance for large-scale
7. **Redis Caching**: Cache frequently accessed data (project/app lists, latest versions, manifest metadata)
8. **Database Indexing**: Optimize indexes on key fields (project, app, created_at) for fast queries
9. **Query Pagination**: All list APIs support pagination (default 50, max 500 per page)
10. **Response Compression**: Gzip compression for all API responses to reduce bandwidth
11. **Parallel SHA256**: Parallel hash calculation for large files to utilize multi-core CPUs
12. **Async Operations**: Large-scale cleanup tasks run asynchronously to avoid blocking

## Risks / Trade-offs

### Risks
- **Migration Complexity**: Existing rsync deployments need manual migration (mitigated: out of scope for V1)
- **Storage Growth**: Immutable versions lead to storage growth (mitigated: implement cleanup policies in future)
- **Token Management**: Token rotation and revocation complexity (mitigated: simple token model for V1)
- **Large File Performance**: Network bandwidth for large artifacts (mitigated: HTTP Range, chunked uploads, retry logic)

### Trade-offs
- **PostgreSQL vs SQLite**: PostgreSQL chosen for better scalability, concurrent access, and production readiness (essential for 2000+ apps)
- **Local Storage vs Object Storage**: Object storage (S3/OSS) chosen as primary backend for V1 to support 2TB+ capacity and distributed access, local filesystem for development only
- **Single Binary vs Microservices**: Monolithic server binary for V1 (simpler deployment, can split later when needed)
- **Global vs Per-App Retention**: Global retention limit chosen for simplicity (no per-app configuration needed), applies uniformly to all apps
- **Scheduled vs On-Demand Cleanup**: Scheduled cleanup (daily at 3 AM) chosen to avoid performance impact during operations, with agent-triggered cleanup as supplement
- **With vs Without Redis**: Redis caching included in V1 to support 2000+ apps query performance (essential for large-scale)
- **Synchronous vs Async Operations**: Large-scale operations (cleanup) run asynchronously to avoid blocking API requests

## Migration Plan

N/A - This is a greenfield implementation. Migration from rsync will be handled separately (out of scope for V1).

## Open Questions

1. Should we support multiple storage backends from the start, or start with filesystem and add S3/GCS later?
   - **Decision**: Start with filesystem, design storage interface for future extension
2. Should token expiration be mandatory or optional?
   - **Decision**: Optional for V1, recommend setting expiration in production
3. Should we support token refresh mechanism?
   - **Decision**: Out of scope for V1, tokens can be re-issued manually
4. Should webhook retries be automatic or manual?
   - **Decision**: Manual retries for V1, add automatic retry with exponential backoff in future

