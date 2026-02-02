## 1. Implementation

- [x] 1.1 在 `AuditLogs.tsx` 的「操作」列渲染逻辑中，为 `labels` 增加 `version_delete: '删除'`。
- [x] 1.2 在 `Dashboard.tsx` 的近期审计活动「操作」渲染逻辑中，为 `labels` 增加 `version_delete: '删除'`。

## 2. Validation

- [x] 2.1 手动验证：审计日志页存在 `version_delete` 记录时，该列显示为「删除」。
- [x] 2.2 手动验证：仪表盘近期活动中若出现 `version_delete`，显示为「删除」。
