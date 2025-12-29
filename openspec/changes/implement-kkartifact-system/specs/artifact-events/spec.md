## ADDED Requirements

### Requirement: Event Generation
The system SHALL generate events for major operations including push completion, pull completion, promote, rollback, and version deletion.

#### Scenario: Generate push event
- **WHEN** an artifact upload is completed successfully
- **THEN** a push event is generated with project, app, version hash, timestamp, and agent identifier
- **AND** the event is made available to webhook handlers

#### Scenario: Generate promote event
- **WHEN** a version is marked as promoted
- **THEN** a promote event is generated with project, app, version hash, and timestamp
- **AND** the event is made available to webhook handlers

#### Scenario: Generate rollback event
- **WHEN** a rollback operation is performed
- **THEN** a rollback event is generated with project, app, from version, to version, and timestamp
- **AND** the event is made available to webhook handlers

#### Scenario: Generate version delete event
- **WHEN** a version is deleted due to retention limit cleanup
- **THEN** a delete event is generated with project, app, version hash, deletion reason, and timestamp
- **AND** the event is made available to webhook handlers

### Requirement: Webhook Configuration
The system SHALL support configuring webhooks that trigger HTTP POST requests to external URLs when events occur. Webhooks SHALL have names, event types, URLs, optional headers, and enabled/disabled status.

#### Scenario: Create webhook
- **WHEN** a POST request is made to `/webhooks` with name, event types, URL, and optional headers
- **THEN** a webhook configuration is created and stored
- **AND** the webhook is enabled by default

#### Scenario: List webhooks
- **WHEN** a GET request is made to `/webhooks`
- **THEN** all configured webhooks are returned
- **AND** webhook configurations include ID, name, event types, URL, status, and creation time

### Requirement: Webhook Triggering
The system SHALL trigger configured webhooks when matching events occur. Webhooks SHALL be triggered asynchronously and SHALL not block the primary operation.

#### Scenario: Trigger webhook on push
- **WHEN** a push event occurs and a webhook is configured for push events
- **THEN** an HTTP POST request is sent to the webhook URL
- **AND** the request body contains event data (project, app, version, timestamp, etc.)
- **AND** the operation completes without waiting for webhook response

#### Scenario: Trigger webhook with headers
- **WHEN** a webhook has custom headers configured
- **THEN** those headers are included in the HTTP POST request
- **AND** standard headers (Content-Type: application/json) are also included

#### Scenario: Filter webhook by project/app
- **WHEN** a webhook is configured with project/app filters
- **THEN** the webhook is only triggered for events matching those filters
- **AND** events for other projects/apps do not trigger the webhook

### Requirement: Webhook Update
The system SHALL support updating webhook configurations including event types, URL, headers, and enabled/disabled status.

#### Scenario: Update webhook
- **WHEN** a PUT request is made to `/webhooks/{id}` with updated configuration
- **THEN** the webhook configuration is updated
- **AND** subsequent events use the updated configuration

### Requirement: Webhook Deletion
The system SHALL support deleting webhook configurations. Deleted webhooks SHALL stop receiving events immediately.

#### Scenario: Delete webhook
- **WHEN** a DELETE request is made to `/webhooks/{id}`
- **THEN** the webhook configuration is removed
- **AND** no further events are sent to that webhook

### Requirement: Webhook Event Payload
Webhook HTTP POST requests SHALL include a JSON payload with event type, timestamp, project, app, version hash, and operation-specific data.

#### Scenario: Push event payload
- **WHEN** a push event triggers a webhook
- **THEN** the POST body contains JSON with event type "push", timestamp, project, app, version hash, builder, and build_time

#### Scenario: Promote event payload
- **WHEN** a promote event triggers a webhook
- **THEN** the POST body contains JSON with event type "promote", timestamp, project, app, and version hash

### Requirement: Webhook Failure Handling
The system SHALL handle webhook delivery failures gracefully. Failed webhooks SHALL not affect the primary operation, and failures SHALL be logged for auditing.

#### Scenario: Webhook timeout
- **WHEN** a webhook HTTP POST request times out
- **THEN** the failure is logged
- **AND** the primary operation completes successfully
- **AND** no retry is attempted (future enhancement)

#### Scenario: Webhook HTTP error
- **WHEN** a webhook returns a non-2xx HTTP status code
- **THEN** the failure is logged with the status code
- **AND** the primary operation completes successfully

### Requirement: Webhook Audit Log
The system SHALL maintain a log of webhook triggers including timestamp, webhook ID, event type, HTTP status code, and response time for auditing purposes.

#### Scenario: Log webhook trigger
- **WHEN** a webhook is triggered
- **THEN** an audit log entry is created with webhook ID, event type, timestamp, URL, status code, and response time
- **AND** the log entry is stored for future querying

