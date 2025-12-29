## 1. Component Refactoring

- [x] 1.1 Remove table-based layout from Projects.tsx
- [x] 1.2 Import Ant Design Tree component
- [x] 1.3 Create tree data structure transformation function (projects + apps → tree nodes)
- [x] 1.4 Implement project node rendering (with folder icon, expand/collapse)
- [x] 1.5 Implement app node rendering (with app icon, action button)
- [x] 1.6 Add expand/collapse state management
- [x] 1.7 Add click handlers for navigation

## 2. Search Functionality

- [x] 2.1 Add search input at top of tree view
- [x] 2.2 Implement debounce for search (300ms)
- [x] 2.3 Implement filtering logic for projects
- [x] 2.4 Implement filtering logic for apps
- [x] 2.5 Handle search that matches apps (show parent project, expand it)
- [x] 2.6 Clear search functionality

## 3. Data Fetching

- [x] 3.1 Fetch all projects on page load
- [x] 3.2 Implement lazy loading for apps (fetch when project expanded)
- [x] 3.3 Add React Query queries for apps per project
- [x] 3.4 Cache apps data to avoid refetching
- [x] 3.5 Handle loading states (skeleton/spinner for apps)
- [x] 3.6 Handle error states

## 4. UI/UX Enhancements

- [x] 4.1 Add appropriate icons (FolderOutlined for projects, AppstoreOutlined for apps)
- [x] 4.2 Style tree with proper indentation
- [x] 4.3 Add action buttons (View Versions) to app nodes
- [ ] 4.4 Add app count badge to project nodes (optional)
- [x] 4.5 Implement empty states (no projects, no apps, no search results)
- [x] 4.6 Add loading indicators for apps being fetched
- [x] 4.7 Ensure responsive design works

## 5. State Management

- [x] 5.1 Manage expanded keys state
- [x] 5.2 Handle default expanded state (all collapsed)
- [x] 5.3 Auto-expand projects with matching apps when searching
- [ ] 5.4 Clear expanded state when search is cleared (optional)

## 6. Integration and Testing

- [ ] 6.1 Ensure navigation to versions page works correctly
- [ ] 6.2 Test with multiple projects and apps
- [ ] 6.3 Test search functionality (project search, app search, combined)
- [ ] 6.4 Test expand/collapse behavior
- [ ] 6.5 Test lazy loading (verify apps are fetched on expand)
- [ ] 6.6 Test with empty states
- [ ] 6.7 Verify responsive design
- [ ] 6.8 Test keyboard navigation (if implemented)

## 7. Optional Enhancements

- [ ] 7.1 Add search highlighting (highlight matching text in names)
- [ ] 7.2 Persist expanded state to localStorage
- [ ] 7.3 Add "Expand All" / "Collapse All" buttons
- [ ] 7.4 Add project/app count statistics at top

## Notes

- Focus on core tree view functionality first
- Search highlighting and localStorage persistence can be added later
- Keep existing routes intact for backward compatibility
- Performance: Consider virtual scrolling if dataset is very large (>1000 projects)

