## ADDED Requirements

### Requirement: Public Root Page
The root page (/) SHALL be publicly accessible without authentication and SHALL display an inventory list of all projects, apps, and versions in the system.

#### Scenario: Access root page without authentication
- **WHEN** a user (authenticated or not) navigates to the root path (/)
- **THEN** a public inventory page is displayed
- **AND** no authentication is required to view the page
- **AND** the page displays all projects, apps, and versions in a hierarchical list
- **AND** the page has a simple layout without sidebar navigation

### Requirement: Inventory List Display
The root page SHALL display a hierarchical inventory list of all projects, apps, and versions in the system. The list SHALL use a tree view format where projects contain apps, and apps contain versions.

#### Scenario: View full inventory on root page
- **WHEN** the user navigates to the root path (/)
- **THEN** an inventory list is displayed
- **AND** all projects are shown as top-level nodes
- **AND** each project node can be expanded to show its apps
- **AND** each app node can be expanded to show its versions
- **AND** each item displays its name and creation time
- **AND** the tree structure reflects the hierarchical relationship (Project → App → Version)

### Requirement: Inventory Search and Filter
The root page inventory list SHALL support real-time search and filtering by project name, app name, or version hash/identifier. The search SHALL filter results across all three entity types simultaneously.

#### Scenario: Search inventory by project name
- **WHEN** the user types a project name in the inventory search box on the root page
- **THEN** projects matching the search term are displayed
- **AND** matching projects are automatically expanded to show their apps and versions
- **AND** non-matching projects are hidden
- **AND** apps and versions within matching projects are shown even if they don't match the search term

#### Scenario: Search inventory by app name
- **WHEN** the user types an app name in the inventory search box on the root page
- **THEN** apps matching the search term are displayed
- **AND** parent projects containing matching apps are automatically expanded and displayed
- **AND** all versions within matching apps are shown
- **AND** projects and apps that don't contain matches are hidden

#### Scenario: Search inventory by version hash
- **WHEN** the user types a version hash or identifier in the inventory search box on the root page
- **THEN** versions matching the search term are displayed
- **AND** parent projects and apps containing matching versions are automatically expanded and displayed
- **AND** projects, apps, and versions that don't contain matches are hidden

#### Scenario: Clear search filter
- **WHEN** the user clears the search input on the root page
- **THEN** the full inventory list is displayed again
- **AND** all projects, apps, and versions are visible
- **AND** the tree structure returns to its default expanded/collapsed state

### Requirement: Inventory Navigation
The inventory list SHALL provide navigation links to detailed views for projects, apps, and versions. Clicking on an item SHALL navigate to the appropriate detail page. If authentication is required for the destination page, the user SHALL be redirected to login.

#### Scenario: Navigate to project apps from inventory (authenticated user)
- **WHEN** an authenticated user clicks on a project name in the inventory list
- **THEN** the user is navigated to the Apps page for that project (`/projects/{project}/apps`)
- **AND** the Apps page displays all apps for the selected project

#### Scenario: Navigate to project apps from inventory (unauthenticated user)
- **WHEN** an unauthenticated user clicks on a project name in the inventory list
- **THEN** the user is redirected to the login page
- **AND** after successful login, the user is navigated to the Apps page for the selected project

#### Scenario: Navigate to app versions from inventory
- **WHEN** the user clicks on an app name in the inventory list
- **THEN** if authenticated, the user is navigated to the Versions page for that app (`/projects/{project}/apps/{app}/versions`)
- **OR** if not authenticated, the user is redirected to login and then to the Versions page
- **AND** the Versions page displays all versions for the selected app

#### Scenario: Navigate to version details from inventory
- **WHEN** the user clicks on a version in the inventory list
- **THEN** if authenticated, the user is navigated to the Versions page for that app with the version selected
- **OR** if not authenticated, the user is redirected to login and then to the Versions page
- **AND** the Versions page displays the selected version's details

### Requirement: Inventory Loading States
The inventory list SHALL display appropriate loading states while data is being fetched from public API endpoints.

#### Scenario: View inventory while loading
- **WHEN** the root page is loading inventory data
- **THEN** a loading indicator is displayed in the inventory list area
- **AND** once data is loaded, the inventory list is displayed
- **AND** the page remains accessible (does not redirect or require authentication)

### Requirement: Inventory Empty States
The inventory list SHALL display appropriate empty states when there is no data or when search filters result in no matches.

#### Scenario: View empty inventory
- **WHEN** the system has no projects, apps, or versions
- **THEN** the inventory list displays an empty state message
- **AND** the message indicates that no artifacts are available
- **AND** the page remains accessible without authentication

#### Scenario: View empty search results
- **WHEN** the user searches for a term that matches no projects, apps, or versions
- **THEN** the inventory list displays an empty state message
- **AND** the message indicates that no items match the search criteria
- **AND** a suggestion to try a different search term is provided

### Requirement: Public API Endpoints
The system SHALL provide public (unauthenticated) API endpoints for reading project, app, and version lists.

#### Scenario: Fetch projects via public API
- **WHEN** a client makes a GET request to `/api/v1/public/projects` without authentication
- **THEN** the request succeeds and returns a list of all projects
- **AND** the response includes project names and creation times

#### Scenario: Fetch apps via public API
- **WHEN** a client makes a GET request to `/api/v1/public/projects/{project}/apps` without authentication
- **THEN** the request succeeds and returns a list of all apps for that project
- **AND** the response includes app names and creation times

#### Scenario: Fetch versions via public API
- **WHEN** a client makes a GET request to `/api/v1/public/projects/{project}/apps/{app}/versions` without authentication
- **THEN** the request succeeds and returns a list of all versions for that app
- **AND** the response includes version hashes and creation times

