# Design: Projects Tree View

## Overview

This document outlines the design decisions for refactoring the projects page from a table layout to a tree/collapsible view that displays the hierarchical structure of Projects → Apps.

## UI/UX Design

### Tree Structure

```
Projects (Header)
├── 🔍 Search: [________________]  (filters projects and apps)
│
├── 📁 Project 1
│   ├── 📱 App 1.1  → [View Versions]
│   ├── 📱 App 1.2  → [View Versions]
│   └── 📱 App 1.3  → [View Versions]
│
├── 📁 Project 2 (collapsed)
│
└── 📁 Project 3
    └── 📱 App 3.1  → [View Versions]
```

### Component Selection

**Option 1: Ant Design Tree Component**
- **Pros**: Built-in tree functionality, supports lazy loading, expand/collapse, search filtering
- **Cons**: May need customization for action buttons (View Versions)

**Option 2: Ant Design Collapse Component**
- **Pros**: Simple collapse/expand, easy to add custom content (apps list with buttons)
- **Cons**: Less tree-like, more like accordion

**Option 3: Custom Tree Implementation**
- **Pros**: Full control over UI and behavior
- **Cons**: More code to maintain, need to implement expand/collapse logic

**Recommended**: Use Ant Design `Tree` component with custom `title` render function to include action buttons.

### Search Functionality

**Unified Search**:
- Single search input at the top
- Filters both projects and apps by name (case-insensitive)
- Real-time filtering with debounce (300ms)
- Filtered tree shows:
  - Projects that match search term (always expanded)
  - Apps that match search term (shown under their parent project)
  - Parent projects of matching apps (even if project name doesn't match)

**Search Behavior**:
```
User searches "api":
- Project "myproject" with App "api-server" → Shows project and highlights app
- Project "api-gateway" → Shows project (name matches)
- Project "backend" with App "backend-api" → Shows project and highlights app
```

### Data Fetching Strategy

**Phase 1: Load All Projects, Lazy Load Apps**
1. On page load: Fetch all projects
2. When project expanded: Fetch apps for that project
3. Cache apps data per project (don't refetch if already loaded)

**API Calls**:
- Initial: `GET /api/v1/projects?limit=10000&offset=0`
- On expand: `GET /api/v1/projects/{project}/apps?limit=1000&offset=0`

**Caching**: Use React Query's queryKey to cache apps per project

### Expanded State

**Default**: All projects collapsed (better initial load performance)

**User Preference** (Optional):
- Remember expanded projects in component state
- Optionally persist to localStorage (future enhancement)

### Visual Design

**Project Node**:
- Icon: Folder icon (📁 or Ant Design `FolderOutlined`)
- Text: Project name
- Click: Toggle expand/collapse
- Optional badge: Show app count (e.g., "3 apps")

**App Node**:
- Icon: App icon (📱 or Ant Design `AppstoreOutlined`)
- Text: App name
- Click: Navigate to versions page
- Action button: "View Versions" (link button)

**Spacing**: Indentation to show hierarchy (16px or 24px per level)

### Loading States

- Projects loading: Show skeleton or spinner
- Apps loading (when expanding): Show loading indicator next to project name
- Empty states:
  - No projects: "No projects found. Sync storage to create projects."
  - No apps in project: Show "No apps" message under project
  - Search no results: "No projects or apps match your search."

### Responsive Design

- Desktop: Full tree view with all features
- Tablet: Tree view, may need horizontal scrolling if project names are long
- Mobile: Consider stacked list instead of tree (future enhancement)

## Technical Design

### Component Structure

```typescript
ProjectsTreeView
├── SearchInput (debounced)
├── Tree
│   ├── ProjectNode (multiple)
│   │   ├── ProjectTitle (clickable, expand/collapse)
│   │   └── AppNodes (children, loaded on expand)
│   │       └── AppNode (clickable, navigate to versions)
└── Loading/Empty states
```

### State Management

```typescript
- searchTerm: string
- debouncedSearchTerm: string (300ms debounce)
- expandedKeys: string[] (project IDs or keys)
- loadedProjects: Project[] (all projects)
- appsByProject: Record<projectId, App[]> (cached apps)
```

### Filtering Logic

1. Filter projects by name (client-side)
2. Filter apps by name (client-side)
3. Build tree structure:
   - Include projects that match OR have matching apps
   - For projects with matching apps: always expand
   - Show only matching apps under each project

### Data Transformation

```typescript
// Transform flat data to tree structure
interface TreeDataNode {
  key: string
  title: React.ReactNode
  children?: TreeDataNode[]
  isLeaf?: boolean
  project?: Project
  app?: App
}
```

## Performance Considerations

1. **Lazy Loading**: Only load apps when project is expanded
2. **Debounce**: 300ms debounce on search input
3. **Memoization**: Memoize filtered tree data structure
4. **Virtual Scrolling**: For large datasets, consider virtual scrolling (future)
5. **Pagination**: Not needed for tree view (all visible items loaded)

## Migration Path

1. Implement tree view on `/projects` page
2. Keep existing `/projects/:project/apps` route for direct navigation
3. Users can use either interface based on preference
4. Future: Consider deprecating separate apps page if tree view is preferred

## Accessibility

- Keyboard navigation: Support arrow keys to navigate tree
- Screen reader: Proper ARIA labels for tree nodes
- Focus management: Maintain focus when expanding/collapsing
- Search: Keyboard accessible search input

