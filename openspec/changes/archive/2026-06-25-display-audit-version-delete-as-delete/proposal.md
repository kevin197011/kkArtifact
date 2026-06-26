# Change: 审计日志页将 version_delete 显示为「删除」

## Why

审计日志接口返回的操作类型包含 `version_delete`（版本删除），当前前端未为其配置中文标签，列表中会直接显示英文 `version_delete`。用户希望在审计日志（及展示审计的仪表盘）中看到统一的中文「删除」，与其他操作（推送、拉取、发布等）一致。

## What Changes

- 在审计日志页（`/audit-logs`）的「操作」列中，将 `version_delete` 显示为「删除」。
- 在仪表盘（Dashboard）中展示近期审计记录时，对 `version_delete` 同样显示为「删除」，保持与审计日志页一致。

## Impact

- Affected specs: artifact-web-ui
- Affected code: `web-ui/src/pages/AuditLogs.tsx`（操作列 labels 映射）、`web-ui/src/pages/Dashboard.tsx`（近期活动操作 labels 映射）
