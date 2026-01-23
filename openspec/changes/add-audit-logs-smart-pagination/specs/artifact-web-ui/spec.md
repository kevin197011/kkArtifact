## MODIFIED Requirements

### Requirement: Audit Logs Display
The web UI SHALL display audit logs in a paginated table with accurate total counts and page size selection options.

#### Scenario: Display audit logs with accurate pagination
- **WHEN** the user navigates to the audit logs page
- **THEN** the page displays audit logs in a table
- **AND** the pagination component shows the accurate total count (e.g., "共 1,234 条审计日志")
- **AND** the user can navigate between pages using page numbers or next/previous buttons

#### Scenario: Change page size
- **WHEN** the user selects a different page size (10, 20, 50, or 100 items per page)
- **THEN** the table reloads with the new page size
- **AND** the pagination updates to reflect the new page size
- **AND** the total count remains accurate

#### Scenario: Navigate to specific page
- **WHEN** the user clicks on a page number or uses next/previous buttons
- **THEN** the table loads the corresponding page of audit logs
- **AND** the pagination highlights the current page
- **AND** the total count display remains accurate

#### Scenario: Filter audit logs
- **WHEN** the user applies filters (project, app, operation type, time range) - if implemented
- **THEN** the table displays only matching audit logs
- **AND** the pagination shows the accurate total count of filtered results
- **AND** the page resets to page 1 when filters change
