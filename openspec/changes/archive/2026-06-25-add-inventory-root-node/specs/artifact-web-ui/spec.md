## ADDED Requirements

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
