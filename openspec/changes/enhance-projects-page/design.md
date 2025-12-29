# Design: Enhanced Projects Page

## Overview

This document outlines the design decisions for enhancing the projects page with search and filtering capabilities.

## UI/UX Design

### Search Input
- **Placement**: Above the projects table, aligned to the right or full width
- **Component**: Ant Design `Input.Search` component for better UX (includes search icon)
- **Placeholder**: "Search projects by name..."
- **Behavior**: Real-time filtering with debounce (300ms) to avoid excessive re-renders
- **Case sensitivity**: Case-insensitive matching

### Sorting
- **Options**: 
  - Name (A-Z, Z-A)
  - Creation Date (Newest first, Oldest first)
- **UI**: Dropdown or toggle buttons in table column headers
- **Default**: Creation Date (Newest first) - matches current behavior

### Table Enhancements
- **Styling**: Improve spacing, add hover effects
- **Responsive**: Ensure table is scrollable on smaller screens
- **Pagination**: Fix to show correct total count from API response

### Empty States
- **No projects**: "No projects found. Sync storage to create projects."
- **No search results**: "No projects match your search. Try a different term."

## Technical Design

### Client-Side Filtering (Phase 1)

**Pros:**
- Simple implementation
- No API changes required
- Instant feedback

**Cons:**
- Only filters already loaded data
- Doesn't scale well with large datasets (1000+ projects)

**Implementation:**
```typescript
const [searchTerm, setSearchTerm] = useState('')
const filteredData = useMemo(() => {
  if (!searchTerm) return data
  return data?.filter(project => 
    project.name.toLowerCase().includes(searchTerm.toLowerCase())
  ) || []
}, [data, searchTerm])
```

### Server-Side Search (Phase 2 - Future)

**Pros:**
- Works with large datasets
- More efficient
- Can add advanced filters later

**Cons:**
- Requires API changes
- Slightly more complex implementation

**API Design:**
```
GET /api/v1/projects?limit=50&offset=0&search=term
```

**Database Query:**
```sql
SELECT id, name, created_at FROM projects 
WHERE name ILIKE '%term%'
ORDER BY created_at DESC 
LIMIT $1 OFFSET $2
```

## Data Flow

### Current Flow
```
User → ProjectsPage → projectsApi.list() → API → Database → Response → Table
```

### Enhanced Flow (Client-side search)
```
User → ProjectsPage → projectsApi.list() → API → Database → Response → 
  Filter by searchTerm → Table
```

### Enhanced Flow (Server-side search - future)
```
User → ProjectsPage → projectsApi.list(searchTerm) → API (with search param) → 
  Database (with WHERE clause) → Filtered Response → Table
```

## Performance Considerations

1. **Debounce**: 300ms debounce on search input to reduce re-renders
2. **Memoization**: Use `useMemo` for filtered data
3. **Pagination**: Keep existing pagination, filtering happens on current page initially
4. **Future optimization**: Implement server-side search for datasets > 100 projects

## Accessibility

- Search input should be keyboard accessible
- Table should be navigable via keyboard
- Screen reader announcements for filtered results
- Clear labeling of search and sort controls

