## 1. Public Inventory API

- [x] 1.1 新增 `GET /api/v1/public/inventory` 返回完整清单
- [x] 1.2 Web UI 清单页改为单次请求加载
- [x] 1.3 Dashboard 使用 summary API 避免 fan-out

## 2. Agent Batch Push

- [x] 2.1 `push` 从 `--path` 推断 app/version
- [x] 2.2 新增 `push-tree` 命令
- [x] 2.3 新增 `--publish`、`--skip-existing`、`--dry-run`
- [x] 2.4 跳过未变更版本与未变更文件
- [x] 2.5 `scripts/push-tree.rb` 包装脚本

## 3. Security & Permissions

- [x] 3.1 push/pull/promote 权限校验与 Token 作用域
- [x] 3.2 admin 限制：删除、sync-storage、Token CRUD、Webhook CRUD
- [x] 3.3 Token 创建移至 protected 路由

## 4. Audit & UI

- [x] 4.1 审计日志 `operation` 筛选（API + Web UI）
- [x] 4.2 `version_delete` 显示为「删除」

## 5. CI & Docs

- [x] 5.1 新增 `.github/workflows/ci.yml`
- [x] 5.2 README 更新 push-tree、权限、审计 API 文档

## 6. Validation

- [x] 6.1 `go test ./...` 通过（server/agent）
- [x] 6.2 `npm run build` 通过（web-ui）
