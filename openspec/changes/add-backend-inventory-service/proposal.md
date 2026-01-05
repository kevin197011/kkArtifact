# Proposal: Add Backend Inventory Service

## Summary

Add a backend inventory service that provides methods to retrieve complete inventory data (all projects, apps, and versions) in a structured format. This service will be used by backend internal operations and can also be exposed via API endpoints for administrative purposes.

## Motivation

Currently, the backend has separate repository methods to list projects, apps, and versions, but there's no unified way to get the complete inventory structure. This is needed for:
- Backend internal operations that need to work with the full inventory (e.g., cleanup tasks, reporting, statistics)
- Administrative API endpoints that need to provide complete inventory overview
- Future features that require hierarchical inventory data access

## Goals

1. Create an inventory service that provides methods to retrieve complete inventory data
2. Support retrieving inventory data in a hierarchical structure (Projects → Apps → Versions)
3. Provide both internal service methods and optional API endpoints
4. Ensure efficient data retrieval with proper error handling
5. Make the service reusable for various backend operations

## Non-Goals

- This does not replace existing repository methods (they remain for specific queries)
- This does not change the database schema
- This does not add caching (can be added later if needed)
- This does not provide real-time updates (uses current database state)

## Scope

### In Scope
- Create an `InventoryService` in `server/internal/services/` (or similar location)
- Add methods to retrieve complete inventory data (all projects with their apps and versions)
- Add methods to retrieve inventory for a specific project
- Provide structured data models for inventory representation
- Optionally expose inventory data via API endpoints (for admin/management use)
- Integrate the service into existing Handler structure

### Out of Scope
- Caching layer (can be added later)
- Real-time inventory updates
- Filtering/searching at the service level (can use repositories directly)
- Performance optimization for very large inventories (can be optimized later)

## Impact

### Affected Components
- `server/internal/services/inventory_service.go` - New service file
- `server/internal/api/inventory_handlers.go` - New API handlers (optional)
- `server/internal/api/handlers.go` - Register new routes (if API endpoints are added)

### Dependencies
- Existing database repositories (ProjectRepository, AppRepository, VersionRepository)
- Existing database models

## Success Criteria

1. Backend can retrieve complete inventory data via service methods
2. Inventory data is returned in a hierarchical structure
3. Service methods are efficient and handle errors properly
4. Service can be easily integrated into existing backend code
5. Optional API endpoints work correctly (if implemented)

