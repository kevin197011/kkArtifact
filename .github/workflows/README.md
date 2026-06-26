# GitHub Actions 工作流

## 概览

| 工作流 | 文件 | 触发 | 作用 |
|--------|------|------|------|
| **CI** | `ci.yml` | push / PR → `main` | 单元测试、Web UI 构建、Docker 集成测试 |
| **Build and Release** | `build-and-release.yml` | 推送标签 `v*` / 手动 | 构建 GHCR 镜像、创建 GitHub Release |

## CI 任务

| Job | 内容 |
|-----|------|
| Unit Tests | `server`、`agent` 矩阵：`go build` + `go test`（Go 1.24） |
| Web UI | `npm ci` + `npm run build`（Node 20） |
| Integration | `docker compose up -d --build --wait` → 编译 agent → `ruby scripts/test_integration.rb` |

Integration 依赖 Unit Tests 与 Web UI 均通过后执行。

### 集成测试覆盖

`scripts/test_integration.rb` 自动验证（无交互）：

- 健康检查、公开清单 API
- 废弃 `/public/projects*` 端点 Deprecation 头
- 管理员登录、创建 API Token
- Agent push / pull、发布与 latest 拉取
- 审计日志 `operation` 筛选
- 权限：无 admin 不可删除；无 pull 权限不可 pull

环境变量：

| 变量 | 默认 | 说明 |
|------|------|------|
| `KK_SERVER_URL` | `http://localhost:8080` | API 地址 |
| `KK_AGENT_BIN` | `/tmp/kkartifact-agent` | Agent 二进制路径 |
| `KK_ADMIN_USER` / `KK_ADMIN_PASS` | admin / admin123 | 管理员凭据 |
| `NO_INTERACTION` | — | CI 中设为 `true` |

## 镜像（GHCR）

```
ghcr.io/<owner>/kkartifact/server:<tag>
ghcr.io/<owner>/kkartifact/web-ui:<tag>
```

Server 镜像内含多平台 Agent 二进制与安装脚本。

示例：

```bash
docker pull ghcr.io/kevin197011/kkartifact/server:v1.0.0
docker pull ghcr.io/kevin197011/kkartifact/web-ui:v1.0.0
```

## 发布流程

```bash
git tag v1.0.0
git push origin v1.0.0
# → 构建并推送 GHCR 镜像 → 创建 Release（含部署说明）
```

生产部署：

```bash
export KK_SERVER_IMAGE=ghcr.io/kevin197011/kkartifact/server:v1.0.0
export KK_WEB_UI_IMAGE=ghcr.io/kevin197011/kkartifact/web-ui:v1.0.0
docker compose -f docker-compose.prod.yml pull
docker compose -f docker-compose.prod.yml up -d
```

手动触发：Actions → **Build and Release** → Run workflow

私有 GHCR 仓库需 PAT（`read:packages`）登录：

```bash
echo $GITHUB_TOKEN | docker login ghcr.io -u USERNAME --password-stdin
```

## 本地与 CI 对应

| 本地 | CI |
|------|-----|
| `cd server && go test ./...` | Unit Tests → server |
| `cd agent && go test ./...` | Unit Tests → agent |
| `cd web-ui && npm run build` | Web UI |
| `ruby scripts/build.rb --test` | Integration（等价） |

## 脚本速查

| 脚本 | 用途 |
|------|------|
| `scripts/up.rb` | `build.rb --up`：构建镜像并启动 |
| `scripts/build.rb` | Docker 构建；`--test` 含集成测试；`--agent` 仅编译 CLI |
| `scripts/test_integration.rb` | 集成测试（服务需已运行） |
| `scripts/push-tree.rb` | G02 批量 push 包装 |
| `scripts/run-migrations.rb` | 数据库迁移 |
| `t.rb` | 转发至 `scripts/push-tree.rb` |

```bash
ruby scripts/build.rb --help   # 查看 build 选项
```
