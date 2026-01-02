# Change: Implement kkArtifact System

## Why

The current deployment system relies on rsync + SSH, which has several limitations:
- No project-level permission isolation
- No unified version metadata
- Difficult deployment auditing
- Not suitable for cloud-native / agent-based deployment
- Coarse security model (account-based rather than project-based)
- Lack of visual management interface

This change implements a modern artifact management and synchronization system that replaces rsync + SSH with HTTP + Token-based authentication, supporting multi-project/multi-app/hash-based version management, web UI, and event system with webhooks.

## What Changes

This is a **greenfield implementation** of the kkArtifact system, including:

- **kkArtifact-server**: HTTP API server for artifact storage, version management, and web UI backend
- **kkArtifact-agent**: CLI tool for push/pull operations with manifest generation
- **Token-based authentication**: Multi-level permission model (Global/Project/App scopes)
- **Event system**: Event-driven architecture with webhook support
- **Web management UI**: Frontend interface for project/app/version management
- **Configuration**: `.kkartifact.yml` configuration file format
- **Storage model**: Immutable versioned artifacts with hash-based versioning
- **Object storage support**: S3/OSS compatible storage backend for large-scale (2000+ apps, 2TB+ capacity)
- **Redis caching**: Cache layer for performance optimization (project/app lists, latest versions, manifest metadata)
- **Version retention management**: Global configuration for version retention limit with scheduled cleanup (daily at 3:00 AM)
- **App-level version cleanup**: Agent-triggered cleanup after push operations (per app, independent)
- **Performance optimizations**: Database indexing, API pagination, response compression, parallel SHA256 computation
- **Docker Compose development environment**: One-command startup for all services (server, PostgreSQL, Redis, web-ui) with hot-reload support

**BREAKING**: This is a new system, no existing systems are affected.

## Impact

- **Affected specs**: All capabilities are new (artifact-storage, artifact-api, artifact-agent, artifact-auth, artifact-events, artifact-web-ui)
- **Affected code**: New codebase implementation
- **Infrastructure**: Requires HTTP server deployment and storage backend
- **Migration**: Existing rsync-based deployments will need migration scripts (out of scope for V1)

