## 1. Backend Implementation
- [x] 1.1 Add `Count` method to `AuditRepository` that efficiently counts audit logs with optional filters (project_id, app_id)
- [x] 1.2 Modify `handleListAuditLogs` to call count query and include total in response
- [x] 1.3 Update API response structure to include `total` field alongside `data` array
- [x] 1.4 Verify database indexes exist on `audit_logs` table for efficient queries (created_at, project_id, app_id)
- [ ] 1.5 Add unit tests for count query with various filter combinations (optional, can be added later)

## 2. Frontend Implementation
- [x] 2.1 Update `AuditLog` API response type to include `total` field
- [x] 2.2 Modify `auditApi.list` to return response with `data` and `total` fields
- [x] 2.3 Update `AuditLogsPage` component to use accurate total count from API
- [x] 2.4 Add page size selector (10, 20, 50, 100 items per page) to pagination component
- [x] 2.5 Update pagination display to show accurate total (e.g., "共 1,234 条审计日志")
- [ ] 2.6 Test pagination with various page sizes and verify accurate counts (requires manual testing)

## 3. Testing and Validation
- [ ] 3.1 Test API endpoint with various filter combinations (no filter, project_id, app_id, both)
- [ ] 3.2 Verify count accuracy matches actual data in database
- [ ] 3.3 Test pagination performance with large datasets (1000+ logs)
- [ ] 3.4 Verify frontend pagination works correctly with accurate totals
- [ ] 3.5 Test page size changes and verify data reloads correctly
