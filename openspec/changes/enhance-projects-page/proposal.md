# Change: Enhance Projects Page with Search and Filter

## Why

The current projects page at `/projects` has several usability issues:
- Lacks visual clarity: Simple table layout doesn't provide quick insights
- No search functionality: Users cannot quickly find specific projects by name
- No filtering options: Cannot filter projects by date range or other criteria
- Limited information display: Only shows name and creation date
- No aggregation stats: Doesn't show useful metrics like app count per project
- Pagination doesn't show total count: Current implementation shows `total: data?.length` which is incorrect for server-side pagination

This change improves the projects page to be more intuitive, searchable, and informative, making it easier for users to manage and find projects.

## What Changes

This change enhances the projects page with:

1. **Search functionality**: Add a search input to filter projects by name (client-side filtering initially)
2. **Enhanced table display**: Improve visual presentation with better spacing and formatting
3. **Additional columns**: Show app count per project (requires API enhancement)
4. **Improved pagination**: Fix pagination to properly handle server-side pagination with total count
5. **Better empty states**: Improved empty state messages when no projects match filters
6. **Sorting options**: Add ability to sort by name or creation date
7. **Quick stats**: Optional summary card showing total projects count

**Future consideration** (out of scope for this change but documented):
- Server-side search API endpoint with query parameter
- Filter by date range
- Advanced filters (e.g., projects with no apps)

## Impact

- **Affected specs**: `artifact-web-ui` (Projects List View requirement)
- **Affected code**: 
  - `web-ui/src/pages/Projects.tsx` - Main component enhancement
  - `web-ui/src/api/projects.ts` - Potentially add new API methods if server-side search is needed (future)
  - `server/internal/api/project_handlers.go` - Potentially add search query parameter (future)
  - `server/internal/database/repositories.go` - Potentially add search filtering to List method (future)
- **User-facing**: Improved usability of projects page
- **Breaking changes**: None

