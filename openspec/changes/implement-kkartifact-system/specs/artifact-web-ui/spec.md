## ADDED Requirements

### Requirement: Projects List View
The web UI SHALL display a list of all projects, ordered by creation time (newest first). Each project SHALL be clickable to navigate to project details.

#### Scenario: View projects list
- **WHEN** the user navigates to the projects page
- **THEN** all projects are displayed in a table or list
- **AND** projects are sorted by creation time descending
- **AND** each project shows name and creation time

### Requirement: Apps List View
The web UI SHALL display a list of all apps for a selected project, ordered by creation time (newest first). Each app SHALL be clickable to navigate to app details.

#### Scenario: View apps for project
- **WHEN** the user selects a project
- **THEN** all apps for that project are displayed
- **AND** apps are sorted by creation time descending
- **AND** each app shows name and creation time

### Requirement: Versions List View
The web UI SHALL display a list of all versions for a selected app, ordered by creation time (newest first). Each version SHALL show the hash, creation time, build time, and builder.

#### Scenario: View versions for app
- **WHEN** the user selects an app
- **THEN** all versions for that app are displayed
- **AND** versions are sorted by creation time descending
- **AND** each version shows hash, creation time, build time, builder, and promoted status

### Requirement: Version Manifest View
The web UI SHALL display the manifest (meta.yaml) content for a selected version, including file list with paths, sizes, and SHA256 checksums.

#### Scenario: View version manifest
- **WHEN** the user selects a version
- **THEN** the manifest content is displayed
- **AND** the file list shows path, size, and SHA256 for each file
- **AND** file information is formatted in a readable table

### Requirement: Pull Operation Trigger
The web UI SHALL allow users to trigger pull operations for a specific version to a deployment path. The operation SHALL be initiated via API call to the server.

#### Scenario: Trigger pull from UI
- **WHEN** the user clicks "Pull" button for a version and provides deployment path
- **THEN** a pull operation is triggered via API
- **AND** operation status is displayed to the user
- **AND** a pull event is generated

### Requirement: Version Promotion
The web UI SHALL allow users to promote versions (mark as ready for deployment) via a promote button. The operation SHALL call the promote API endpoint.

#### Scenario: Promote version from UI
- **WHEN** the user clicks "Promote" button for a version
- **THEN** a promote API request is sent
- **AND** the version is marked as promoted
- **AND** a promote event is generated
- **AND** the UI updates to show promoted status

### Requirement: Rollback Operation
The web UI SHALL allow users to trigger rollback operations to a previous version. The operation SHALL pull the specified version to the deployment path.

#### Scenario: Rollback to previous version
- **WHEN** the user selects a previous version and clicks "Rollback"
- **THEN** a pull operation is triggered for that version
- **AND** a rollback event is generated
- **AND** operation status is displayed

### Requirement: Webhook Management
The web UI SHALL provide interfaces for creating, viewing, updating, and deleting webhook configurations. Users SHALL be able to configure webhook name, event types, URL, headers, and enabled status.

#### Scenario: Create webhook from UI
- **WHEN** the user navigates to webhooks page and clicks "Create Webhook"
- **THEN** a form is displayed for entering webhook configuration
- **AND** upon submission, a webhook is created via API
- **AND** the new webhook appears in the webhooks list

#### Scenario: Edit webhook from UI
- **WHEN** the user clicks "Edit" on a webhook
- **THEN** a form is displayed with current webhook configuration
- **AND** upon submission, the webhook is updated via API
- **AND** the updated configuration is reflected in the list

#### Scenario: Delete webhook from UI
- **WHEN** the user clicks "Delete" on a webhook and confirms
- **THEN** the webhook is deleted via API
- **AND** the webhook is removed from the list

### Requirement: Token Management
The web UI SHALL provide interfaces for viewing and managing tokens (create, view permissions, revoke). Token creation SHALL allow specifying scope (Global/Project/App), permissions, and expiration.

#### Scenario: View tokens
- **WHEN** the user navigates to tokens page
- **THEN** all tokens are displayed with scope, permissions, creation time, and expiration
- **AND** token values are never displayed (security)

#### Scenario: Create token from UI
- **WHEN** the user clicks "Create Token" and fills in scope, permissions, and expiration
- **THEN** a token is created via API
- **AND** the token value is displayed once (never again)
- **AND** the token appears in the tokens list

#### Scenario: Revoke token from UI
- **WHEN** the user clicks "Revoke" on a token and confirms
- **THEN** the token is revoked via API
- **AND** the token is marked as revoked in the list

### Requirement: Audit Log View
The web UI SHALL display audit logs for operations including push, pull, promote, rollback, and webhook triggers. Logs SHALL be filterable by project, app, operation type, and time range.

#### Scenario: View audit logs
- **WHEN** the user navigates to audit logs page
- **THEN** operation logs are displayed with timestamp, operation type, project, app, version, and agent identifier
- **AND** logs are sorted by timestamp descending
- **AND** logs can be filtered by project, app, and operation type

### Requirement: API Integration
The web UI SHALL integrate with kkArtifact-server API endpoints for all operations. API calls SHALL use authentication tokens and handle errors appropriately.

#### Scenario: API call with authentication
- **WHEN** the web UI makes an API request
- **THEN** an authentication token is included in the Authorization header
- **AND** the token is stored securely (e.g., in localStorage or session)

#### Scenario: Handle API errors
- **WHEN** an API request fails (4xx or 5xx status)
- **THEN** an appropriate error message is displayed to the user
- **AND** the error message indicates the cause (unauthorized, not found, etc.)

### Requirement: Responsive Design
The web UI SHALL be responsive and work on desktop and tablet devices. The UI SHALL use modern, accessible components.

#### Scenario: Desktop view
- **WHEN** the UI is viewed on a desktop browser
- **THEN** content is displayed in a multi-column layout
- **AND** tables and lists are fully visible

#### Scenario: Tablet view
- **WHEN** the UI is viewed on a tablet device
- **THEN** content adapts to the screen size
- **AND** navigation remains accessible

