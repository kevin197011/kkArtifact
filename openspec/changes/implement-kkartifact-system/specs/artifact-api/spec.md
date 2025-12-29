## ADDED Requirements

### Requirement: Manifest Retrieval API
The system SHALL provide a GET endpoint `/manifest/{project}/{app}/{hash}` that returns the `meta.yaml` manifest for the specified version.

#### Scenario: Retrieve manifest
- **WHEN** a GET request is made to `/manifest/myproject/myapp/a8f3c21d`
- **THEN** the `meta.yaml` content is returned as JSON or YAML
- **AND** the response includes all file metadata with paths, SHA256, and sizes

### Requirement: File Download API
The system SHALL provide a GET endpoint `/file/{project}/{app}/{hash}?path={filepath}` that returns the file content for the specified path.

#### Scenario: Download file
- **WHEN** a GET request is made to `/file/myproject/myapp/a8f3c21d?path=bin/app`
- **THEN** the file content is returned with appropriate Content-Type headers
- **AND** the file matches the SHA256 stored in the manifest

#### Scenario: Download with HTTP Range
- **WHEN** a GET request includes a Range header for partial content
- **THEN** the server responds with 206 Partial Content
- **AND** only the requested byte range is returned

### Requirement: File Upload API
The system SHALL provide a PUT endpoint `/file/{project}/{app}/{hash}?path={filepath}` that accepts file content for upload.

#### Scenario: Upload single file
- **WHEN** a PUT request is made to `/file/myproject/myapp/a8f3c21d?path=bin/app` with file content
- **THEN** the file is stored at the specified path within the version directory
- **AND** the file SHA256 is calculated and stored in the manifest

### Requirement: Batch Upload Init API
The system SHALL provide a POST endpoint `/upload/init` that initializes a new upload session and returns an upload token or session ID.

#### Scenario: Initialize batch upload
- **WHEN** a POST request is made to `/upload/init` with project, app, version, and file list
- **THEN** an upload session is created
- **AND** a session identifier is returned for subsequent file uploads

### Requirement: Batch Upload Finish API
The system SHALL provide a POST endpoint `/upload/finish` that finalizes an upload session, generates the manifest, and marks the version as ready.

#### Scenario: Finalize batch upload
- **WHEN** a POST request is made to `/upload/finish` with the session identifier
- **THEN** the manifest is generated from all uploaded files
- **AND** the version is marked as complete
- **AND** a push event is triggered

### Requirement: Version Promotion API
The system SHALL provide a POST endpoint `/promote` that marks a version hash as promoted (ready for deployment).

#### Scenario: Promote version
- **WHEN** a POST request is made to `/promote` with project, app, and hash
- **THEN** the version is marked as promoted
- **AND** a promote event is triggered

### Requirement: Projects List API
The system SHALL provide a GET endpoint `/projects` that returns all projects, ordered by creation time (newest first).

#### Scenario: List all projects
- **WHEN** a GET request is made to `/projects`
- **THEN** a list of all project names is returned
- **AND** projects are ordered by creation time descending

### Requirement: Apps List API
The system SHALL provide a GET endpoint `/projects/{project}/apps` that returns all apps for a project, ordered by creation time (newest first).

#### Scenario: List apps for project
- **WHEN** a GET request is made to `/projects/myproject/apps`
- **THEN** a list of all app names for that project is returned
- **AND** apps are ordered by creation time descending

### Requirement: Versions List API
The system SHALL provide a GET endpoint `/projects/{project}/apps/{app}/versions` that returns all version hashes for an app, ordered by creation time (newest first).

#### Scenario: List versions for app
- **WHEN** a GET request is made to `/projects/myproject/apps/myapp/versions`
- **THEN** a list of all version hashes for that app is returned
- **AND** versions are ordered by creation time descending

### Requirement: API Authentication
All API endpoints SHALL require authentication via Bearer token in the Authorization header. Requests without valid tokens SHALL be rejected with 401 Unauthorized.

#### Scenario: Authenticated request
- **WHEN** a request includes a valid Bearer token in the Authorization header
- **THEN** the request is processed if the token has required permissions

#### Scenario: Unauthenticated request
- **WHEN** a request does not include an Authorization header or includes an invalid token
- **THEN** the request is rejected with 401 Unauthorized

### Requirement: API Error Responses
The system SHALL return appropriate HTTP status codes and error messages in JSON format for all error conditions.

#### Scenario: Not found error
- **WHEN** a request references a non-existent project, app, or version
- **THEN** the response is 404 Not Found with an error message

#### Scenario: Permission denied
- **WHEN** a request uses a token that lacks required permissions
- **THEN** the response is 403 Forbidden with an error message

### Requirement: Global Configuration API
The system SHALL provide API endpoints for managing global configuration including version retention limits. These endpoints SHALL require admin permissions.

#### Scenario: Get global configuration
- **WHEN** a GET request is made to `/config` with admin token
- **THEN** the global configuration is returned including version retention limit

#### Scenario: Update global version retention
- **WHEN** a PUT request is made to `/config` with admin token and new retention limit
- **THEN** the global version retention limit is updated
- **AND** the new limit applies to all apps for future cleanup operations

