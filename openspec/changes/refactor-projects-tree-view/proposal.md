# Change: Refactor Projects Page to Tree View with Search

## Why

The current projects page uses a flat table layout that requires navigation through multiple pages to browse projects and their apps:
- Users must click "View Apps" to navigate to a separate page for each project
- Cannot see the hierarchical structure (Projects → Apps) in one view
- Search only filters projects, not apps
- No visual representation of the project-app relationship
- Less efficient for users who need to browse multiple projects and apps

This change refactors the projects page to use a tree/collapsible view that displays the hierarchical structure (Projects → Apps) in a single page, similar to a file browser. This provides better visibility and navigation efficiency.

## What Changes

This change refactors the projects page with:

1. **Tree/Collapsible View**: Replace table layout with a tree view showing Projects as parent nodes and Apps as child nodes
2. **Unified Search**: Single search input that filters both projects and apps by name (case-insensitive)
3. **Lazy Loading**: Apps are loaded on-demand when a project is expanded (performance optimization)
4. **Visual Hierarchy**: Clear visual indication of project-app relationships with icons and indentation
5. **Click Actions**: 
   - Click project name: Expand/collapse to show apps
   - Click app name: Navigate to app versions page
   - Maintain "View Apps" button for backward compatibility (optional)
6. **Search Highlighting**: Highlight matching text in project and app names (optional enhancement)
7. **Expanded State Management**: Remember which projects are expanded (optional: persist to localStorage)

**Replaces**: The current table-based projects page layout

**Complements**: The existing `/projects/:project/apps` route can still exist for direct navigation, but the main projects page provides a unified tree view

## Impact

- **Affected specs**: `artifact-web-ui` (Projects List View requirement)
- **Affected code**: 
  - `web-ui/src/pages/Projects.tsx` - Complete refactor to tree view
  - `web-ui/src/api/projects.ts` - No changes (uses existing APIs)
  - Routing: No changes (tree view is on `/projects`, existing routes remain)
- **User-facing**: Major UI/UX improvement - tree view replaces table view
- **Breaking changes**: Visual layout change (table → tree), but functionality remains similar

