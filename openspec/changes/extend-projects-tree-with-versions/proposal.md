# Change: Extend Projects Tree View to Include Versions

## Why

The current projects tree view only displays two levels: Projects → Apps. Users still need to navigate to a separate page to view versions for each app. This creates unnecessary navigation steps and breaks the unified browsing experience.

By extending the tree view to include versions as a third level (Projects → Apps → Versions), users can:
- Browse the complete hierarchy in a single view
- See all versions at a glance without navigating away
- Search for versions across all projects and apps
- Have a more complete overview of the artifact structure

## What Changes

This change extends the existing tree view with:

1. **Third Level - Versions**: Add versions as child nodes under apps in the tree view
2. **Lazy Loading for Versions**: Load versions on-demand when an app is expanded (similar to how apps are loaded when projects expand)
3. **Unified Search Enhancement**: Extend search functionality to include version filtering (search by version hash/identifier)
4. **Version Display**: Show version hash/identifier in the tree node
5. **Version Actions**: Provide action buttons to view version details (manifest, promote, etc.)
6. **Auto-expand Logic**: When searching for versions, auto-expand parent project and app

**Extends**: The existing tree view implemented in `refactor-projects-tree-view` change

**Maintains**: All existing functionality (project and app browsing, search, etc.)

## Impact

- **Affected specs**: `artifact-web-ui` (Projects List View requirement - MODIFIED to include versions)
- **Affected code**: 
  - `web-ui/src/pages/Projects.tsx` - Extend tree structure to three levels
  - `web-ui/src/api/projects.ts` - Uses existing `getVersions` API method (no changes needed)
- **User-facing**: Enhanced tree view showing complete hierarchy (Projects → Apps → Versions)
- **Breaking changes**: None (additive change)

