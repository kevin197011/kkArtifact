## MODIFIED Requirements

### Requirement: Projects List View (Tree View)
The web UI SHALL display projects and apps in a hierarchical tree/collapsible view that shows the project-app relationship. Users SHALL be able to expand projects to view their apps and navigate to app versions.

#### Scenario: View projects tree
- **WHEN** the user navigates to the projects page
- **THEN** all projects are displayed in a tree view structure
- **AND** projects are shown as parent nodes with folder icons
- **AND** apps are shown as child nodes under their parent project when expanded
- **AND** projects are collapsed by default

#### Scenario: Expand project to view apps
- **WHEN** the user clicks on a project node
- **THEN** the project expands to show its apps
- **AND** apps are fetched from the API (lazy loading)
- **AND** a loading indicator is shown while apps are being fetched
- **AND** apps are displayed with app icons under the project

#### Scenario: Search projects and apps
- **WHEN** the user enters text in the search input
- **THEN** both projects and apps are filtered by name (case-insensitive)
- **AND** projects with matching names are shown
- **AND** apps with matching names are shown under their parent project
- **AND** parent projects of matching apps are shown even if project name doesn't match
- **AND** projects containing matching apps are automatically expanded

#### Scenario: Navigate to app versions
- **WHEN** the user clicks on an app node or "View Versions" button
- **THEN** the user is navigated to the versions page for that app
- **AND** the URL follows the pattern `/projects/{project}/apps/{app}/versions`

#### Scenario: Search with no results
- **WHEN** the user enters a search term that matches no projects or apps
- **THEN** an appropriate empty state message is displayed
- **AND** the message indicates no projects or apps match the search

#### Scenario: Project with no apps
- **WHEN** a project has no apps
- **THEN** the project can still be expanded
- **AND** an empty state message is shown under the project (e.g., "No apps")

## ADDED Requirements

### Requirement: Unified Search for Projects and Apps
The web UI SHALL provide a single search input that filters both projects and apps by name simultaneously.

#### Scenario: Search filters both projects and apps
- **WHEN** the user enters a search term in the unified search input
- **THEN** the tree view is filtered to show:
  - Projects whose name contains the search term
  - Apps whose name contains the search term (shown under their parent project)
  - Parent projects of matching apps (even if project name doesn't match)
- **AND** the search is case-insensitive
- **AND** the filtering updates in real-time with debounce

#### Scenario: Clear search
- **WHEN** the user clears the search input
- **THEN** all projects and apps are displayed again
- **AND** the tree returns to default state (projects collapsed)

