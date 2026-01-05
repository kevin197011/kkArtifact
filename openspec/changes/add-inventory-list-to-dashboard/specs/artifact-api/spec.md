## ADDED Requirements

### Requirement: Public Read-Only API Endpoints
The API SHALL provide public (unauthenticated) read-only endpoints for listing projects, apps, and versions. These endpoints SHALL not require authentication and SHALL return the same data structure as their protected counterparts.

#### Scenario: List projects via public endpoint
- **WHEN** a client makes a GET request to `/api/v1/public/projects` without an Authorization header
- **THEN** the request succeeds (HTTP 200)
- **AND** the response contains a JSON array of projects
- **AND** each project includes id, name, and created_at fields
- **AND** the response format matches the protected `/api/v1/projects` endpoint

#### Scenario: List apps via public endpoint
- **WHEN** a client makes a GET request to `/api/v1/public/projects/{project}/apps` without an Authorization header
- **THEN** the request succeeds (HTTP 200)
- **AND** the response contains a JSON array of apps for the specified project
- **AND** each app includes id, project_id, name, and created_at fields
- **AND** the response format matches the protected `/api/v1/projects/{project}/apps` endpoint

#### Scenario: List versions via public endpoint
- **WHEN** a client makes a GET request to `/api/v1/public/projects/{project}/apps/{app}/versions` without an Authorization header
- **THEN** the request succeeds (HTTP 200)
- **AND** the response contains a JSON array of versions for the specified app
- **AND** each version includes id, app_id, version (hash), and created_at fields
- **AND** the response format matches the protected `/api/v1/projects/{project}/apps/{app}/versions` endpoint

#### Scenario: Public endpoints reject non-GET requests
- **WHEN** a client makes a POST, PUT, DELETE, or PATCH request to any `/api/v1/public/*` endpoint
- **THEN** the request fails with HTTP 405 (Method Not Allowed)
- **AND** public endpoints only accept GET requests

#### Scenario: Public endpoints handle invalid project/app names
- **WHEN** a client requests apps or versions for a non-existent project or app
- **THEN** the request returns HTTP 404 (Not Found)
- **AND** the error response indicates the resource was not found

### Requirement: Public API Security
Public API endpoints SHALL be read-only and SHALL not expose sensitive information or allow any mutations to the system.

#### Scenario: Public endpoints are read-only
- **WHEN** a client attempts to modify data through public endpoints
- **THEN** the request is rejected (HTTP 405 Method Not Allowed)
- **AND** no data is modified

#### Scenario: Public endpoints do not expose tokens or credentials
- **WHEN** a client requests data via public endpoints
- **THEN** the response does not include authentication tokens
- **AND** the response does not include user credentials or sensitive metadata
- **AND** only basic artifact information (name, creation time, hash) is included

