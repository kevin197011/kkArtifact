## MODIFIED Requirements

### Requirement: Projects List View (Tree View with Versions)
The web UI SHALL display projects, apps, and versions in a hierarchical tree/collapsible view that shows the complete project-app-version relationship. Users SHALL be able to expand projects to view apps, expand apps to view versions, and navigate to version details.

#### Scenario: View projects tree with versions
- **WHEN** the user navigates to the projects page
- **THEN** all projects are displayed in a tree view structure
- **AND** projects are shown as parent nodes with folder icons
- **AND** apps are shown as child nodes under their parent project when expanded
- **AND** versions are shown as child nodes under their parent app when app is expanded
- **AND** projects and apps are collapsed by default

#### Scenario: Expand app to view versions
- **WHEN** the user clicks on an app node to expand it
- **THEN** the app expands to show its versions
- **AND** versions are fetched from the API (lazy loading)
- **AND** a loading indicator is shown while versions are being fetched
- **AND** versions are displayed with version icons under the app

#### Scenario: Search projects, apps, and versions
- **WHEN** the user enters text in the search input
- **THEN** projects, apps, and versions are filtered by name/hash (case-insensitive)
- **AND** projects with matching names are shown
- **AND** apps with matching names are shown under their parent project
- **AND** versions with matching hashes are shown under their parent app
- **AND** parent projects and apps of matching versions are shown even if their names don't match
- **AND** projects and apps containing matching versions are automatically expanded

#### Scenario: Navigate to version details
- **WHEN** the user clicks on a version node or "View Manifest" button
- **THEN** the user is navigated to the versions page for that app and version
- **AND** the URL follows the pattern `/projects/{project}/apps/{app}/versions`

#### Scenario: Search with no results
- **WHEN** the user enters a search term that matches no projects, apps, or versions
- **THEN** an appropriate empty state message is displayed
- **AND** the message indicates no projects, apps, or versions match the search

#### Scenario: App with no versions
- **WHEN** an app has no versions
- **THEN** the app can still be expanded
- **AND** an empty state message is shown under the app (e.g., "No versions")

## ADDED Requirements

### Requirement: Versions in Tree View
The web UI SHALL display versions as child nodes under their parent apps in the projects tree view, with lazy loading and search support.

#### Scenario: Versions displayed in tree
- **WHEN** a user expands an app node in the projects tree
- **THEN** versions for that app are displayed as child nodes
- **AND** each version node shows the version hash/identifier
- **AND** each version node has action buttons (e.g., "View Manifest")
- **AND** versions are loaded on-demand (lazy loading)

#### Scenario: Search filters versions
- **WHEN** the user enters a search term that matches version hashes
- **THEN** matching versions are shown in the tree
- **AND** parent apps and projects are shown and expanded to display matching versions
- **AND** the search is case-insensitive

