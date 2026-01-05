## ADDED Requirements

### Requirement: Backend Inventory Service
The backend SHALL provide an inventory service that retrieves complete inventory data (all projects, apps, and versions) in a hierarchical structure. The service SHALL provide methods to get complete inventory, project-specific inventory, and inventory summaries.

#### Scenario: Get complete inventory via service
- **WHEN** the backend calls `InventoryService.GetCompleteInventory()`
- **THEN** the service returns a complete inventory structure containing all projects
- **AND** each project includes all its apps
- **AND** each app includes all its versions
- **AND** the data is organized in a hierarchical structure
- **AND** the service handles errors appropriately (database errors, missing data)

#### Scenario: Get inventory for specific project via service
- **WHEN** the backend calls `InventoryService.GetProjectInventory(projectName)` with a valid project name
- **THEN** the service returns inventory data for that project
- **AND** the returned data includes the project, all its apps, and all versions for each app
- **AND** if the project does not exist, the service returns an appropriate error

#### Scenario: Get inventory summary via service
- **WHEN** the backend calls `InventoryService.GetInventorySummary()`
- **THEN** the service returns summary statistics
- **AND** the summary includes total count of projects, apps, and versions
- **AND** the summary is calculated efficiently

### Requirement: Inventory Service Data Structure
The inventory service SHALL return data in a structured format that represents the hierarchical relationship between projects, apps, and versions.

#### Scenario: Inventory structure hierarchy
- **WHEN** inventory data is retrieved
- **THEN** the structure contains a list of projects
- **AND** each project contains a list of apps
- **AND** each app contains a list of versions
- **AND** the structure preserves the parent-child relationships

### Requirement: Optional Admin Inventory API Endpoints
The backend MAY provide administrative API endpoints to access inventory data. These endpoints SHALL require authentication and appropriate permissions.

#### Scenario: Get complete inventory via admin API
- **WHEN** an authenticated admin user makes a GET request to `/api/v1/admin/inventory`
- **THEN** the request succeeds and returns complete inventory data in JSON format
- **AND** the response includes all projects with their apps and versions
- **AND** unauthenticated or non-admin requests are rejected with appropriate error

#### Scenario: Get project inventory via admin API
- **WHEN** an authenticated admin user makes a GET request to `/api/v1/admin/inventory/:project`
- **THEN** the request succeeds and returns inventory data for the specified project
- **AND** the response includes the project with all its apps and versions
- **AND** if the project does not exist, the request returns 404 Not Found
- **AND** unauthenticated or non-admin requests are rejected with appropriate error

#### Scenario: Get inventory summary via admin API
- **WHEN** an authenticated admin user makes a GET request to `/api/v1/admin/inventory/summary`
- **THEN** the request succeeds and returns inventory summary statistics
- **AND** the response includes total counts of projects, apps, and versions
- **AND** unauthenticated or non-admin requests are rejected with appropriate error

