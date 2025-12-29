## 1. Tree Structure Extension

- [x] 1.1 Extend TreeDataNode interface to include version field and isVersion flag
- [x] 1.2 Update tree data building logic to handle three levels (projects → apps → versions)
- [x] 1.3 Add version node rendering (with version icon, hash display, action buttons)
- [x] 1.4 Update expand/collapse logic to handle app expansion (for loading versions)
- [x] 1.5 Add version-level click handlers for navigation

## 2. Data Fetching for Versions

- [x] 2.1 Implement lazy loading for versions (fetch when app expanded)
- [x] 2.2 Add React Query queries for versions per app
- [x] 2.3 Cache versions data to avoid refetching (similar to apps caching)
- [x] 2.4 Handle loading states for versions (spinner when loading)
- [x] 2.5 Handle error states for version fetching

## 3. Search Functionality Enhancement

- [x] 3.1 Extend filtering logic to include versions (by hash/identifier)
- [x] 3.2 Update search to load versions when searching for versions
- [x] 3.3 Handle search that matches versions (show parent project and app, expand them)
- [x] 3.4 Update empty state messages to include versions

## 4. UI/UX Enhancements

- [x] 4.1 Add appropriate icon for versions (FileZipOutlined or FileOutlined)
- [x] 4.2 Style tree with three-level indentation
- [x] 4.3 Add action buttons to version nodes (View Manifest, Promote if available)
- [x] 4.4 Implement empty states (no versions in app)
- [x] 4.5 Add loading indicators for versions being fetched
- [x] 4.6 Ensure proper visual hierarchy (three levels clearly distinguishable)

## 5. State Management

- [x] 5.1 Update expanded keys state to include version keys
- [x] 5.2 Handle default expanded state (all collapsed, including versions)
- [x] 5.3 Auto-expand projects and apps when searching for versions
- [x] 5.4 Manage version data caching (versionsByApp record)

## 6. Integration and Testing

- [x] 6.1 Ensure navigation to version details works correctly
- [x] 6.2 Test with projects that have multiple apps and versions
- [x] 6.3 Test search functionality (project search, app search, version search, combined)
- [x] 6.4 Test expand/collapse behavior for all three levels
- [x] 6.5 Test lazy loading (verify versions are fetched on app expand)
- [x] 6.6 Test with empty states (no versions)
- [x] 6.7 Verify three-level tree displays correctly
- [x] 6.8 Test performance with large number of versions

## 7. Optional Enhancements

- [ ] 7.1 Add version metadata display (created_at, file count, etc.) in tooltip or node
- [ ] 7.2 Add "Expand All" / "Collapse All" buttons that work for all three levels
- [ ] 7.3 Add version count badge to app nodes
- [ ] 7.4 Implement version search highlighting

## Notes

- Build on existing tree view implementation
- Follow the same patterns used for apps lazy loading
- Maintain backward compatibility (existing functionality should still work)
- Performance: Consider pagination or virtual scrolling for apps with many versions (>100)

