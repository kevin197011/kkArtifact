# Change: Add Smart Pagination for Audit Logs

## Why

The current audit logs page uses simple offset-based pagination without accurate total count, leading to:
- Inaccurate pagination display (showing "at least X logs" instead of exact count)
- Poor user experience when navigating through large audit log datasets
- No efficient way to know total number of logs or jump to specific pages
- Performance concerns with large datasets (offset pagination becomes slow for deep pages)

This change implements smart pagination with accurate total counts, optimized queries, and better UX for browsing audit logs.

## What Changes

- **MODIFIED**: Audit logs API endpoint to return total count along with paginated results
- **MODIFIED**: Audit logs repository to support efficient count queries with filters
- **MODIFIED**: Frontend audit logs page to display accurate pagination with total count
- **ADDED**: Database index optimization for audit logs queries (if not already present)
- **ADDED**: Support for page size selection in frontend (10, 20, 50, 100 items per page)

**BREAKING**: None - API response format is extended (backward compatible), frontend behavior is improved.

## Impact

- **Affected specs**: `artifact-api` (audit logs API), `artifact-web-ui` (audit logs page)
- **Affected code**: 
  - `server/internal/api/audit_handlers.go` - Add total count to response
  - `server/internal/database/audit_repository.go` - Add count query method
  - `web-ui/src/pages/AuditLogs.tsx` - Update pagination with accurate totals
  - `web-ui/src/api/audit.ts` - Update API response type
- **User impact**: Users can now see accurate total counts and navigate audit logs more efficiently
