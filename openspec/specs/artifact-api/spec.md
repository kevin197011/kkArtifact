# artifact-api Specification

## Purpose
TBD - created by archiving change add-g02-batch-migration-workflow. Update Purpose after archive.
## Requirements
### Requirement: Public Inventory API
The server SHALL expose a public read-only endpoint that returns the complete inventory hierarchy (projects, apps, and versions) in a single response without authentication.

#### Scenario: Get complete public inventory
- **WHEN** an unauthenticated client requests `GET /api/v1/public/inventory`
- **THEN** the response includes all projects with nested apps and versions
- **AND** no authentication header is required

### Requirement: Audit Log Operation Filter
The server SHALL support filtering audit logs by operation type via query parameter.

#### Scenario: Filter audit logs by operation
- **WHEN** an authenticated client requests `GET /api/v1/audit-logs?operation=push`
- **THEN** only audit log entries with operation `push` are returned
- **AND** the response includes paginated `data` and `total` count reflecting the filter

### Requirement: API Permission Enforcement
The server SHALL enforce fine-grained permissions on artifact and management endpoints.

#### Scenario: Push requires push permission
- **WHEN** a caller without `push` permission attempts upload endpoints
- **THEN** the request is rejected with HTTP 403

#### Scenario: Pull requires pull permission
- **WHEN** a caller without `pull` permission attempts manifest or file download
- **THEN** the request is rejected with HTTP 403

#### Scenario: Publish requires promote permission
- **WHEN** a caller without `promote` permission attempts publish or unpublish
- **THEN** the request is rejected with HTTP 403

#### Scenario: Admin-only management operations
- **WHEN** a non-admin caller attempts delete project/app/version, sync-storage, token create/delete, or webhook create/update/delete
- **THEN** the request is rejected with HTTP 403

#### Scenario: Token scope enforcement
- **WHEN** an API token is scoped to a project or app
- **THEN** artifact operations are only allowed within that scope

