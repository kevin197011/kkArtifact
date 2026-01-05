# Proposal: Add Public Inventory List to Root Page

## Summary

Create a public (non-authenticated) root page (/) that displays a comprehensive inventory list showing all projects, apps, and versions in the system with filtering and search capabilities. This provides anyone with quick access to browse all artifacts without requiring login.

## Motivation

Currently, accessing artifact information requires authentication and navigation to separate pages. This proposal adds a public inventory list on the root page, making it easier for users (including non-authenticated users) to browse and search all artifacts. This is useful for sharing artifact information or providing a quick overview without login barriers.

## Goals

1. Create a public root page (/) that displays a unified inventory list showing all projects, apps, and versions
2. Support filtering and search across all three entity types
3. Provide quick navigation to detailed views (which may still require authentication)
4. Allow access without authentication (public page)
5. Ensure the page works independently from the authenticated Dashboard

## Non-Goals

- This does not replace the existing Dashboard, Projects, Apps, or Versions pages
- This does not require authentication (it's a public page)
- Detailed views (versions page, etc.) may still require authentication
- This may require new public API endpoints or modification of existing endpoints to support public access

## Scope

### In Scope
- Create a new public root page component (InventoryPage or HomePage)
- Add route for root path (/) that doesn't require authentication
- Implement client-side filtering/search functionality
- Support hierarchical display (Projects → Apps → Versions)
- Add search input for filtering by project name, app name, or version hash
- Display basic information for each entity (name, creation time)
- Add navigation links to detailed views (may redirect to login if authentication required)
- Either use existing API endpoints without auth token, or create public API endpoints for listing projects/apps/versions
- Create a simple public layout (without sidebar/menu, or minimal header with login link)

### Out of Scope
- Server-side search/filtering endpoints
- Real-time updates (current polling/cache refresh is sufficient)
- Advanced filtering (e.g., by date range, status)
- Bulk operations on the inventory list
- Export functionality
- Making detailed views (versions, etc.) public (they remain protected)

## Impact

### Affected Components
- `web-ui/src/App.tsx` - Add public route for root path
- `web-ui/src/pages/InventoryPage.tsx` - New public inventory page component (or HomePage.tsx)
- `server/internal/api/handlers.go` - May need to add public endpoints for projects/apps/versions listing (OR modify existing endpoints to allow optional authentication)
- API client configuration - May need to support API calls without authentication token

### Dependencies
- Projects, Apps, and Versions API endpoints (need to support public access)
- React Query for data fetching and caching
- Ant Design components (Table, Input, Tree, etc.)
- Public page layout (simple layout without authenticated features)

## Success Criteria

1. Users can access the root page (/) without authentication
2. Users can see all projects, apps, and versions in a single view on the root page
3. Users can filter/search by project name, app name, or version hash
4. Search results update in real-time as user types
5. Users can navigate to detailed views from the inventory list (may require login)
6. The page loads efficiently without authentication
7. The page has a clean, simple layout appropriate for public access
8. Existing Dashboard and other authenticated pages continue to work as before

