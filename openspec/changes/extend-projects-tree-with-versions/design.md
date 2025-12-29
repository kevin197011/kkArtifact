# Design: Extend Projects Tree View with Versions

## Overview

This document outlines the design for extending the projects tree view to include versions as a third hierarchical level, creating a complete Projects → Apps → Versions browsing experience.

## UI/UX Design

### Extended Tree Structure

```
Projects (Header)
├── 🔍 Search: [________________]  (filters projects, apps, and versions)
│
├── 📁 Project 1
│   ├── 📱 App 1.1
│   │   ├── 📦 version-abc123  → [View Manifest] [Promote]
│   │   ├── 📦 version-def456  → [View Manifest] [Promote]
│   │   └── 📦 version-ghi789  → [View Manifest] [Promote]
│   ├── 📱 App 1.2
│   │   └── 📦 version-xyz999  → [View Manifest] [Promote]
│   └── 📱 App 1.3 (collapsed)
│
├── 📁 Project 2 (collapsed)
│
└── 📁 Project 3
    └── 📱 App 3.1
        └── 📦 version-aaa111  → [View Manifest] [Promote]
```

### Component Structure

The tree will now have three levels:
1. **Project** (level 1): Folder icon, project name
2. **App** (level 2): App icon, app name, expandable
3. **Version** (level 3): Package/version icon, version hash, action buttons

### Visual Design

**Version Node**:
- Icon: Package icon (📦 or Ant Design `FileZipOutlined` / `FileOutlined`)
- Text: Version hash/identifier
- Actions: 
  - "View Manifest" button (opens manifest modal or navigates to versions page)
  - "Promote" button (optional, if promote functionality is available)
  - Click version name: Navigate to versions page for details

**Indentation**: Three levels with consistent spacing (16px or 24px per level)

### Search Functionality Enhancement

**Extended Unified Search**:
- Filters projects by name
- Filters apps by name
- **NEW**: Filters versions by hash/identifier
- When searching for versions:
  - Show parent project (even if project name doesn't match)
  - Show parent app (even if app name doesn't match)
  - Auto-expand project and app to show matching versions
  - Highlight matching version

### Data Fetching Strategy

**Three-Level Lazy Loading**:
1. On page load: Fetch all projects
2. When project expanded: Fetch apps for that project (existing)
3. **NEW**: When app expanded: Fetch versions for that app

**API Calls**:
- Initial: `GET /api/v1/projects?limit=10000&offset=0` (existing)
- On project expand: `GET /api/v1/projects/{project}/apps?limit=1000&offset=0` (existing)
- **NEW**: On app expand: `GET /api/v1/projects/{project}/apps/{app}/versions?limit=1000&offset=0`

**Caching**: Use React Query's queryKey to cache versions per app

### Expanded State Management

**Default**: 
- All projects collapsed
- Apps collapsed (when project is expanded)
- **NEW**: Versions collapsed (when app is expanded)

**When Searching**:
- Projects with matching items are auto-expanded
- **NEW**: Apps with matching versions are auto-expanded
- All matching items are visible

### Loading States

- Projects loading: Show spinner (existing)
- Apps loading: Show spinner next to project name (existing)
- **NEW**: Versions loading: Show spinner next to app name when app is expanded

### Empty States

- No projects: "No projects found. Sync storage to create projects." (existing)
- No apps in project: "No apps" under project (existing)
- **NEW**: No versions in app: "No versions" under app
- **NEW**: Search no results: "No projects, apps, or versions match your search."

## Technical Design

### Tree Data Structure Extension

```typescript
interface TreeDataNode {
  key: string
  title: React.ReactNode
  children?: TreeDataNode[]
  isLeaf?: boolean
  project?: Project
  app?: App
  version?: Version  // NEW
  isProject?: boolean
  isApp?: boolean
  isVersion?: boolean  // NEW
}
```

### State Management

```typescript
- searchTerm: string
- debouncedSearchTerm: string (300ms debounce)
- expandedKeys: string[] (project IDs, app IDs, version IDs)
- loadedProjects: Project[] (all projects)
- appsByProject: Record<projectId, App[]> (cached apps)
- versionsByApp: Record<appId, Version[]> (NEW - cached versions per app)
```

### Filtering Logic Extension

1. Filter projects by name (existing)
2. Filter apps by name (existing)
3. **NEW**: Filter versions by hash/identifier
4. Build tree structure:
   - Include projects that match OR have matching apps OR have matching versions
   - Include apps that match OR have matching versions
   - Show only matching versions under each app
   - Auto-expand all parent levels when matching items are found

### Data Transformation

The tree building logic needs to handle three levels:
```typescript
Projects (level 1)
  └── Apps (level 2)
      └── Versions (level 3)
```

When building tree nodes:
1. Process projects (level 1)
2. For each project, process its apps (level 2) - if project is expanded or has matching items
3. **NEW**: For each app, process its versions (level 3) - if app is expanded or has matching items

### Expand/Collapse Handling

**When App is Expanded**:
- Fetch versions for that app (lazy load)
- Display versions as children nodes
- Cache versions data

**When App is Collapsed**:
- Versions are hidden
- Versions data remains in cache (for quick re-expansion)

**When Searching for Versions**:
- Auto-expand parent project and app
- Load versions if not already cached
- Show only matching versions

## Performance Considerations

1. **Three-Level Lazy Loading**: Only load data when parent is expanded
2. **Caching**: Cache apps and versions data separately
3. **Debounce**: 300ms debounce on search (existing)
4. **Memoization**: Memoize filtered tree data structure (existing, extended for versions)
5. **Batch Loading**: When searching, batch-load versions for all apps (similar to current app loading on search)

## User Interaction

**Click Actions**:
- Click project: Toggle expand/collapse to show apps (existing)
- Click app: Toggle expand/collapse to show versions (**NEW**)
- Click version: Navigate to version details page (or open manifest modal)
- Click "View Manifest": Open manifest modal or navigate to versions page
- Click "Promote": Trigger promote action (if implemented)

## Migration Path

1. Extend existing tree view implementation
2. Add version fetching logic (similar to app fetching)
3. Extend tree data structure building to include versions
4. Update search filtering to include versions
5. Test with projects that have multiple apps and versions

## Accessibility

- Keyboard navigation: Support arrow keys to navigate three-level tree
- Screen reader: Proper ARIA labels for all three levels
- Focus management: Maintain focus when expanding/collapsing any level
- Search: Keyboard accessible, filters all three levels

