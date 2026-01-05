## Implementation Tasks

### Phase 0: Backend API (Public Endpoints)

- [x] 0.1 Add public API endpoints for listing projects (`/api/v1/public/projects`)
- [x] 0.2 Add public API endpoints for listing apps (`/api/v1/public/projects/:project/apps`)
- [x] 0.3 Add public API endpoints for listing versions (`/api/v1/public/projects/:project/apps/:app/versions`)
- [x] 0.4 Ensure public endpoints are read-only (GET only, no mutations)
- [x] 0.5 Update API client to support calls without authentication token
- [x] 0.6 Test public endpoints return correct data

### Phase 1: Data Fetching and State Management

- [x] 1.1 Add React Query hooks to fetch all projects, apps, and versions for the inventory
- [x] 1.2 Implement efficient data fetching strategy (lazy load or batch load)
- [x] 1.3 Structure data to support hierarchical display (Projects → Apps → Versions)
- [x] 1.4 Add caching strategy for inventory data to avoid unnecessary API calls

### Phase 2: UI Components

- [x] 2.1 Create new public root page component (`InventoryPage.tsx` or `HomePage.tsx`)
- [x] 2.2 Create simple public layout (header only, no sidebar)
- [x] 2.3 Create inventory list component structure
- [x] 2.4 Implement search/filter input component
- [x] 2.5 Design layout for inventory list (tree view)
- [x] 2.6 Add loading states for inventory data
- [x] 2.7 Add empty states when no data matches filters
- [x] 2.8 Add login link/navigation in header (optional)

### Phase 3: Filtering and Search

- [x] 3.1 Implement client-side filtering logic for project names
- [x] 3.2 Implement client-side filtering logic for app names
- [x] 3.3 Implement client-side filtering logic for version hashes
- [x] 3.4 Add debounce to search input to optimize performance
- [x] 3.5 Ensure filtered results maintain hierarchical structure

### Phase 4: Navigation and Integration

- [x] 4.1 Add route for root path (/) without ProtectedRoute wrapper
- [x] 4.2 Configure route to render new public inventory page
- [x] 4.3 Add navigation links from inventory items to detailed pages (may require login)
- [x] 4.4 Handle navigation when user is not authenticated (redirect to login or show message)
- [x] 4.5 Test responsive design on different screen sizes
- [x] 4.6 Add proper spacing and styling for public page
- [x] 4.7 Ensure existing Dashboard and authenticated pages continue to work

### Phase 5: Testing and Polish

- [x] 5.1 Test with different data sizes (few items, many items)
- [x] 5.2 Test filtering with various search terms
- [x] 5.3 Verify navigation links work correctly
- [x] 5.4 Test performance with large datasets
- [x] 5.5 Ensure accessibility (keyboard navigation, screen readers)

