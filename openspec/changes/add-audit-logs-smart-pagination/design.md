## Context

The audit logs page currently uses offset-based pagination without accurate total counts. This leads to poor UX when users need to navigate through large audit log datasets. The change needs to:
- Add total count to API responses
- Optimize count queries for performance
- Update frontend to display accurate pagination

## Goals / Non-Goals

### Goals
- Return accurate total count in audit logs API response
- Maintain query performance (<500ms for count queries on large datasets)
- Provide better UX with accurate pagination display
- Support page size selection (10, 20, 50, 100 items per page)

### Non-Goals
- Cursor-based pagination (offset-based is sufficient for current scale)
- Real-time count updates (count is calculated on each request)
- Advanced filtering beyond existing project_id/app_id filters (out of scope)

## Decisions

### Decision 1: API Response Format
**What**: Extend API response to include `total` field alongside `data` array.

**Why**: 
- Backward compatible (existing clients can ignore `total` field)
- Simple to implement and understand
- Matches common REST API pagination patterns

**Implementation**:
```json
{
  "data": [...],
  "total": 1234
}
```

### Decision 2: Count Query Strategy
**What**: Use separate COUNT query with same WHERE clause as data query.

**Why**:
- PostgreSQL COUNT queries are efficient with proper indexes
- Same filter logic ensures count matches data
- Can be optimized with indexes on `created_at`, `project_id`, `app_id`

**Alternatives considered**:
- Window functions: More complex, similar performance
- Cached counts: Adds complexity, may be stale
- Estimated counts: Inaccurate, poor UX

### Decision 3: Database Indexes
**What**: Ensure indexes exist on `audit_logs(created_at DESC)`, `audit_logs(project_id)`, `audit_logs(app_id)`.

**Why**:
- `created_at` is used for ORDER BY (most common query pattern)
- `project_id` and `app_id` are used for filtering
- Composite indexes may be added if query patterns require

**Implementation**: Verify indexes exist in migration or add if missing.

### Decision 4: Page Size Options
**What**: Support 10, 20, 50, 100 items per page (default: 50).

**Why**:
- Common page size options for data tables
- Balances performance (fewer queries) and UX (manageable page size)
- Default 50 matches current implementation

## Risks / Trade-offs

### Risk: COUNT Query Performance on Large Datasets
**Mitigation**: 
- Ensure proper database indexes
- Monitor query performance in production
- Consider approximate counts if performance becomes an issue (future enhancement)

### Trade-off: Accuracy vs Performance
- Accurate counts require COUNT query (adds ~50-200ms)
- Trade-off is acceptable for better UX
- Can optimize with materialized views or caching if needed later

## Migration Plan

1. **Backend**: Add `Count` method to repository, update handler to include total
2. **Frontend**: Update API types and component to use total count
3. **Database**: Verify/ensure indexes exist (no migration needed if already present)
4. **Testing**: Verify accuracy and performance with various dataset sizes

## Open Questions

- Should we add composite indexes for common filter combinations (project_id + created_at)?
  - **Decision**: Monitor query patterns first, add if needed based on actual usage
