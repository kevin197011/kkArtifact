# artifact-web-ui Specification

## Purpose
TBD - created by archiving change add-inventory-root-node. Update Purpose after archive.
## Requirements
### Requirement: Inventory Tree Root Node
The public inventory page (root page `/`) SHALL display the project tree under a single root node so that users can collapse or expand all projects at once.

#### Scenario: Root node present when data exists
- **WHEN** the user visits the inventory page and there is at least one project
- **THEN** the tree displays one root node (e.g. "全部项目" or "项目") at the top level
- **AND** all project nodes are children of that root node
- **AND** the hierarchy is root → projects → apps → versions

#### Scenario: Collapse all via root node
- **WHEN** the user collapses the root node
- **THEN** the entire tree (all projects and their children) is collapsed
- **AND** only the root node remains visible in expanded form until the user expands it again

#### Scenario: Expand all via root node
- **WHEN** the user expands the root node
- **THEN** all project nodes become visible as children of the root
- **AND** project-level expand/collapse state is preserved as before (e.g. previously expanded projects stay expanded)

#### Scenario: Search with root node
- **WHEN** the user enters a search term and results are shown
- **THEN** the root node remains present and is expanded so that matching projects/apps/versions are visible under it
- **AND** search and filter behavior is unchanged from the current implementation

#### Scenario: No data or empty result
- **WHEN** there are no projects or no search matches
- **THEN** the existing empty state is shown (no misleading single root node with no children, or root node is omitted when there are no projects)
- **AND** the user experience is consistent with the current empty/loading behavior

### Requirement: Audit Log Operation Type Labels
The web UI SHALL display human-readable labels for audit log operation types. The operation type `version_delete` SHALL be displayed as 「删除」 on the audit logs page and on the dashboard where audit entries are shown.

#### Scenario: version_delete displayed as 删除 on audit logs page
- **WHEN** the user navigates to the audit logs page
- **AND** the list contains an audit log entry with operation `version_delete`
- **THEN** the operation column for that entry displays 「删除」 (not the raw value `version_delete`)

#### Scenario: version_delete displayed as 删除 on dashboard
- **WHEN** the user views the dashboard recent activity (audit log entries)
- **AND** an entry has operation `version_delete`
- **THEN** the operation for that entry is displayed as 「删除」 (not the raw value `version_delete`)

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

