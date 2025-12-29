# Design: Public Inventory List on Root Page

## Overview

Create a public (non-authenticated) root page (/) that displays a comprehensive inventory list showing all projects, apps, and versions in a searchable, filterable format. This page is accessible without login.

## UI Design Decisions

### Layout Option: Tree View vs Table vs Cards

**Decision: Tree View (Collapsible Hierarchical List)**

Rationale:
- Naturally represents the hierarchy: Projects → Apps → Versions
- Allows expanding/collapsing to manage screen space
- Familiar pattern for users navigating file systems or project structures
- Supports filtering at any level while maintaining hierarchy
- Consistent with common artifact management UIs

Alternative considered: Flat table with columns for Project/App/Version
- Rejected: Less intuitive for hierarchical data, harder to see relationships

### Search/Filter Placement

**Decision: Search input above the inventory list**

- Single search box that filters across all three entity types
- Real-time filtering as user types (with debounce)
- Clear icon and placeholder text indicating searchable fields

### Data Loading Strategy

**Decision: Progressive/On-Demand Loading**

- Load all projects on Dashboard load (already done)
- Load apps for each project on-demand when project is expanded (or pre-load for inventory)
- Load versions for each app on-demand when app is expanded

For the inventory view specifically:
- Option A: Load all data upfront (projects, apps, versions) - simpler but may be slow
- Option B: Lazy load apps/versions as user expands nodes - better performance
- Option C: Pre-load first N items, load more on scroll/expand - balanced approach

**Decision: Option A for initial implementation (can optimize later)**
- Dashboard is expected to show full inventory overview
- Data should be cached by React Query, so subsequent loads are fast
- If performance becomes an issue, we can switch to lazy loading

### Visual Hierarchy

```
Root Page (/)
└── Public Inventory List
    ├── Header (with login link if needed)
    ├── Search Input
    └── Tree View
        ├── Project 1
        │   ├── App 1
        │   │   ├── Version 1
        │   │   └── Version 2
        │   └── App 2
        └── Project 2
```

Note: This is separate from the Dashboard page. The Dashboard (at /dashboard) remains authenticated and unchanged.

### Filtering Logic

When user types in search box:
1. Filter projects by name (case-insensitive)
2. Filter apps by name (case-insensitive) - show parent project even if project name doesn't match
3. Filter versions by hash/version string (case-insensitive) - show parent project and app even if they don't match
4. Auto-expand nodes that contain matching items
5. Show "No results" message if nothing matches

## Technical Implementation

### Authentication & API Access

**Decision: Create public API endpoints or make existing endpoints support optional auth**

Options:
- Option A: Add new public endpoints (e.g., `/api/v1/public/projects`, `/api/v1/public/projects/:project/apps`)
- Option B: Modify existing endpoints to allow optional authentication (if no token provided, allow read-only access)
- Option C: Use existing endpoints without auth token (requires backend changes to support)

**Decision: Option A - Add new public read-only endpoints**
- Cleaner separation of concerns
- Explicit security model (public vs protected)
- Easier to audit and maintain
- Can have different rate limiting if needed

Alternative considered: Option B (optional auth)
- Rejected: Less secure, harder to maintain, mixes concerns

### Page Layout

**Decision: Simple public layout without sidebar**

- Minimal header with logo and optional "Login" link
- No sidebar navigation (since user is not authenticated)
- Clean, focused design for the inventory list
- Footer optional (minimal or none)

### Component Structure

```
RootPage (new component at /)
  └── PublicInventoryList (new component)
      ├── PublicHeader (optional, with login link)
      ├── InventorySearch (search input)
      └── InventoryTree (tree view component)
          └── InventoryNode (recursive component)
```

### Data Structure

```typescript
interface InventoryItem {
  type: 'project' | 'app' | 'version'
  id: number
  name: string
  createdAt: string
  projectName?: string // for app and version
  appName?: string // for version
  versionHash?: string // for version
  parentId?: number // for hierarchical relationships
}
```

### Component Structure

```
Dashboard
  └── InventoryList (new component)
      ├── InventorySearch (search input)
      └── InventoryTree (tree view component)
          └── InventoryNode (recursive component)
```

### State Management

- Use React Query for data fetching
- Use local state for:
  - Search/filter term
  - Expanded keys (which nodes are expanded)
  - Debounced search term

### Performance Considerations

1. Debounce search input (300ms delay)
2. Memoize filtered data with `useMemo`
3. Use React Query's caching to avoid refetching
4. Consider virtualization if tree becomes very large (future optimization)

## Alternative Designs Considered

### Flat Table with Filters
- Pros: Simple, familiar UI pattern
- Cons: Doesn't show hierarchy well, harder to navigate

### Tabbed Interface (Projects / Apps / Versions)
- Pros: Clear separation, can use existing table components
- Cons: Requires switching tabs, doesn't show relationships

### Accordion/Collapsible Cards
- Pros: Modern UI, good for mobile
- Cons: Harder to see all data at once, less space-efficient

## Future Enhancements (Out of Scope)

- Server-side search/filtering
- Advanced filters (date range, status)
- Bulk operations
- Export to CSV/JSON
- Bookmark/favorite items
- Recently viewed items

