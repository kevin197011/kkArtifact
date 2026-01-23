## MODIFIED Requirements

### Requirement: Audit Logs List API
The API SHALL provide an endpoint to list audit logs with pagination support. The endpoint SHALL support filtering by project ID and app ID, and SHALL return both the paginated results and the total count of matching records.

#### Scenario: List audit logs with pagination
- **WHEN** a client requests audit logs with limit and offset parameters
- **THEN** the API returns a JSON response containing:
  - `data`: An array of audit log entries (limited by `limit` parameter)
  - `total`: The total count of audit logs matching the filter criteria (if any)
- **AND** the response includes pagination metadata for accurate frontend display

#### Scenario: List audit logs with project filter
- **WHEN** a client requests audit logs filtered by project_id
- **THEN** the API returns only audit logs for that project
- **AND** the `total` field reflects the count of logs for that project only

#### Scenario: List audit logs with app filter
- **WHEN** a client requests audit logs filtered by app_id
- **THEN** the API returns only audit logs for that app
- **AND** the `total` field reflects the count of logs for that app only

#### Scenario: List audit logs with both project and app filters
- **WHEN** a client requests audit logs filtered by both project_id and app_id
- **THEN** the API returns only audit logs matching both filters
- **AND** the `total` field reflects the count of logs matching both criteria

#### Scenario: Count query performance
- **WHEN** the audit logs table contains 10,000+ records
- **THEN** the count query completes in <500ms
- **AND** database indexes are utilized for efficient counting
