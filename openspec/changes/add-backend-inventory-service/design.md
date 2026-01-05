# Design: Backend Inventory Service

## Overview

Create a service layer that provides methods to retrieve complete inventory data (projects, apps, versions) in a hierarchical structure. This service will abstract the complexity of fetching and organizing inventory data from multiple repositories.

## Design Decisions

### Service Location and Structure

**Decision: Create `InventoryService` in `server/internal/services/`**

Rationale:
- Follows separation of concerns (service layer between API handlers and repositories)
- Makes the service reusable across different parts of the backend
- Clear naming convention for services
- Easy to extend with additional inventory-related operations

### Data Structure

**Decision: Use nested structures to represent hierarchical inventory**

```go
type Inventory struct {
    Projects []ProjectInventory `json:"projects"`
}

type ProjectInventory struct {
    Project Project            `json:"project"`
    Apps    []AppInventory     `json:"apps"`
}

type AppInventory struct {
    App      App       `json:"app"`
    Versions []Version `json:"versions"`
}
```

This structure naturally represents the hierarchy and is easy to serialize to JSON if needed for API responses.

### Service Methods

**Decision: Provide multiple methods for different use cases**

1. `GetCompleteInventory()` - Returns all projects with all apps and versions
2. `GetProjectInventory(projectName string)` - Returns inventory for a specific project
3. `GetInventorySummary()` - Returns summary statistics (counts)

This provides flexibility for different use cases without over-fetching data.

### API Endpoints (Optional)

**Decision: Add optional API endpoints for administrative access**

- `GET /api/v1/admin/inventory` - Get complete inventory (requires admin auth)
- `GET /api/v1/admin/inventory/:project` - Get inventory for a specific project (requires admin auth)

These endpoints are optional and can be used for administrative tools or monitoring.

### Error Handling

**Decision: Return errors from service methods, let callers decide how to handle**

Service methods return `(data, error)` tuples, allowing flexibility in error handling at the caller level.

### Performance Considerations

For initial implementation:
- Load all data in memory (acceptable for moderate-sized inventories)
- Use existing repository methods (already optimized with SQL queries)
- No caching initially (can be added later if needed)

For large inventories, consider:
- Pagination support
- Lazy loading
- Caching layer

## Technical Implementation

### Service Interface

```go
type InventoryService interface {
    GetCompleteInventory() (*Inventory, error)
    GetProjectInventory(projectName string) (*ProjectInventory, error)
    GetInventorySummary() (*InventorySummary, error)
}

type Inventory struct {
    Projects []ProjectInventory
}

type ProjectInventory struct {
    Project Project
    Apps    []AppInventory
}

type AppInventory struct {
    App      App
    Versions []Version
}

type InventorySummary struct {
    TotalProjects int
    TotalApps     int
    TotalVersions int
}
```

### Implementation Details

1. **GetCompleteInventory()**:
   - Fetch all projects using `ProjectRepository.List()`
   - For each project, fetch apps using `AppRepository.ListByProject()`
   - For each app, fetch versions using `VersionRepository.ListByApp()`
   - Build hierarchical structure
   - Return complete inventory

2. **GetProjectInventory(projectName)**:
   - Fetch project by name
   - Fetch apps for the project
   - Fetch versions for each app
   - Return project inventory

3. **GetInventorySummary()**:
   - Use efficient SQL queries to get counts (or aggregate in memory)
   - Return summary statistics

### Integration Points

1. **Handler Integration**:
   - Add `InventoryService` to `Handler` struct
   - Initialize service in `NewHandler()` using existing repositories
   - Use service methods in API handlers (if endpoints are added)

2. **Repository Usage**:
   - Service uses existing repositories (ProjectRepository, AppRepository, VersionRepository)
   - No changes needed to repositories
   - Service is a client of repositories

## Alternative Designs Considered

### 1. Add methods directly to Handler
- **Rejected**: Mixes concerns, less reusable

### 2. Create a single SQL query with JOINs
- **Rejected**: More complex, harder to maintain, loses repository abstraction

### 3. Add methods to repositories
- **Rejected**: Repositories should be simple data access, not business logic

### 4. Use existing API endpoints internally
- **Rejected**: Adds unnecessary HTTP overhead, less efficient

## Future Enhancements (Out of Scope)

- Caching layer for frequently accessed inventory
- Pagination support for large inventories
- Filtering/searching at service level
- Real-time inventory updates via events
- Inventory change tracking/diff capabilities

