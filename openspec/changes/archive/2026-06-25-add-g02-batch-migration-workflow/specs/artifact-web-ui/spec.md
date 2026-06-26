## ADDED Requirements

### Requirement: Inventory Page Single Fetch
The public inventory page SHALL load complete inventory data via a single API request instead of per-project fan-out.

#### Scenario: Load inventory in one request
- **WHEN** the user opens the public inventory page
- **THEN** the page fetches `GET /api/v1/public/inventory` once
- **AND** renders the project → app → version tree from that response

### Requirement: Audit Log Operation Filter UI
The audit logs page SHALL provide a filter control for operation type.

#### Scenario: Filter audit logs in UI
- **WHEN** the user selects an operation type in the audit logs filter
- **THEN** the page requests audit logs with the corresponding `operation` query parameter
- **AND** only matching entries are displayed
