## MODIFIED Requirements

### Requirement: Projects List View (Enhanced)
The web UI SHALL display a list of all projects with search and filtering capabilities. Each project SHALL be clickable to navigate to project details.

#### Scenario: View projects list with search
- **WHEN** the user navigates to the projects page
- **THEN** all projects are displayed in a table
- **AND** a search input field is available above the table
- **AND** projects are sorted by creation time descending by default
- **AND** each project shows name and creation time
- **AND** pagination displays the correct total count of projects

#### Scenario: Search projects by name
- **WHEN** the user enters text in the search input field
- **THEN** the projects table is filtered in real-time to show only projects whose name contains the search term (case-insensitive)
- **AND** the filtered results update as the user types (with debounce)
- **AND** if no projects match the search term, an appropriate empty state message is displayed

#### Scenario: Sort projects
- **WHEN** the user selects a sort option (name or creation date)
- **THEN** the projects table is reordered according to the selected sort criteria
- **AND** the sort direction can be toggled (ascending/descending)
- **AND** the sort state is visually indicated

#### Scenario: Pagination with search
- **WHEN** the user has filtered projects by search term and navigates to a different page
- **THEN** the search filter is maintained across page navigation
- **AND** pagination controls reflect the correct total count of filtered results

#### Scenario: Clear search
- **WHEN** the user clears the search input field
- **THEN** all projects are displayed again
- **AND** the original sort order is maintained

