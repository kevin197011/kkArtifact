# Change: G02 批量迁移工作流（公开清单、批量推送与安全加固）

## Why

G02 项目采用固定目录布局 `/data/vcs/G02/tidb/{app}/{version_hash}/`，包含 20+ 应用。现有单次 `push` 命令需手动指定 app/version，公开清单页面存在 N+1 请求，且部分管理接口缺少权限校验，不适合日常批量迁移与公开浏览场景。

## What Changes

- **公开清单 API**：`GET /api/v1/public/inventory` 单次返回完整项目→应用→版本树
- **Agent 批量推送**：新增 `push-tree` 命令，从路径推断 app/version；支持 `--skip-existing`、`--publish`
- **推送优化**：跳过未变更版本/文件；配置文件中可设置默认 `project`
- **审计日志**：API 与 Web UI 支持 `operation` 筛选；`version_delete` 显示为「删除」
- **权限加固**：push/pull/promote 细粒度校验；删除、同步存储、Token/Webhook 管理需 admin
- **CI**：PR/push 到 main 时运行 Go 测试与 Web UI 构建

## Impact

- Affected specs: artifact-api, artifact-agent, artifact-web-ui
- Affected code: server API handlers, agent CLI, web-ui InventoryPage/AuditLogs, scripts/push-tree.rb, `.github/workflows/ci.yml`
