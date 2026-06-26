# kkArtifact

替代 `rsync + SSH` 的 Artifact 管理与同步系统，支持多项目 / 多应用 / Hash 版本、Token 鉴权、Web UI 与审计日志。

生产实例：<https://packages.slileisure.com/>

## 特性

- 并发 Push/Pull、断点续传、文件 Hash 校验
- 不可变版本存储，支持发布（latest）与版本保留策略
- Token 细粒度权限（push / pull / promote / admin）
- 公开版本清单（无需登录）、中文 Web 管理后台
- Agent `push-tree` 批量迁移（G02 等固定目录布局）
- Docker Compose 开发与 GHCR 镜像发布

## 快速开始

```bash
# 构建并启动（推荐）
ruby scripts/up.rb

# 或直接使用 docker compose
docker compose up -d --build
```

| 服务 | 地址 |
|------|------|
| Web UI | http://localhost:3000 |
| 公开清单 | http://localhost:3000/ |
| API | http://localhost:8080 |
| Swagger | http://localhost:3000/swagger/index.html |
| 默认账号 | `admin` / `admin123` |

```bash
docker compose logs -f    # 查看日志
docker compose down     # 停止
```

## Agent

### 安装

**安装脚本（推荐）**

```bash
# Linux / macOS
curl -fsSL http://localhost:8080/api/v1/downloads/scripts/install-agent.sh | bash

# 指定服务器地址
curl -fsSL https://packages.slileisure.com/api/v1/downloads/scripts/install-agent.sh \
  | server_url="https://packages.slileisure.com" bash
```

**其他方式**：从 [GitHub Releases](https://github.com/kevin197011/kkArtifact/releases) 下载对应平台二进制，或 `cd agent && go build -o kkartifact-agent .`

`server_url` 优先级：`server_url` 环境变量 → `SERVER_URL` → 脚本检测地址 → `http://localhost:8080`

### 配置

复制 `.kkartifact.yml.sample` 为 `.kkartifact.yml`：

```yaml
server_url: https://packages.slileisure.com
token: YOUR_TOKEN          # Web UI → Token 管理 创建
project: g02               # 默认项目，可省略 --project
concurrency: 50            # 本地 30–50，远程 50–100
ignore: []
```

| 参数 | 说明 |
|------|------|
| `server_url` | API 地址（经 Web UI 代理时用前端 URL，直连 API 用 8080） |
| `token` | API Token，需具备对应操作权限 |
| `project` | 默认项目名 |
| `concurrency` | 并发文件数，过高可能触发 TCP 缓冲区错误 |
| `ignore` | glob 忽略规则，`[]` 表示不忽略 |

### Push / Pull

```bash
# 显式指定
kkartifact-agent push --project myproject --app myapp --version v1.0.0 --path ./dist

# 从路径推断 app/version（路径末尾两段）
kkartifact-agent push --project g02 \
  --path /data/vcs/G02/tidb/G02_agent_api/<hash>

# 推送后自动发布
kkartifact-agent push --project g02 --path .../<hash> --publish

# 下载
kkartifact-agent pull --project g02 --app G02_agent_api --version <hash> --path ./deploy
```

### Push-tree（批量上传）

目录布局：`root/{app}/{version}/`，适合 G02 等批量迁移。

```bash
kkartifact-agent push-tree /data/vcs/G02/tidb \
  --project g02 \
  --skip-existing \
  --publish
```

| 参数 | 说明 |
|------|------|
| `--skip-existing` | 跳过服务端已有版本 |
| `--dry-run` | 仅打印计划 |
| `--publish` | 每个版本推送后自动发布（Token 需 `push` + `promote`） |

Ruby 包装（默认 root 为 `/data/vcs/G02/tidb`）：

```bash
ruby scripts/push-tree.rb /data/vcs/G02/tidb --project g02 --skip-existing --publish
# 或项目根目录：ruby t.rb --project g02 --skip-existing
```

## G02 批量迁移

典型本地布局：

```text
/data/vcs/G02/tidb/
  G02_agent_api/<version_hash>/
  G02_other_app/<version_hash>/
  ...
```

推荐流程：

1. Web UI 创建项目 `g02`（若不存在）
2. 创建 Token，权限勾选 **push**、**promote**（若使用 `--publish`）
3. 配置 `.kkartifact.yml`（`project: g02`）
4. 先 `--dry-run` 确认，再正式推送：

```bash
ruby scripts/push-tree.rb /data/vcs/G02/tidb --project g02 --dry-run
ruby scripts/push-tree.rb /data/vcs/G02/tidb --project g02 --skip-existing --publish
```

5. 在 https://packages.slileisure.com/ 公开清单或管理后台核对版本

## Web UI

**公开清单**（`/`）：单次请求 `GET /api/v1/public/inventory`，树形展示项目 → 应用 → 版本，支持搜索。

**管理后台**（需登录）：项目 / 应用 / 版本、Token、Webhook、配置、审计日志（支持 `operation` 筛选）、仪表盘。

## API 概要

认证：`Authorization: Bearer YOUR_TOKEN`

### Token 权限

| 权限 | 说明 |
|------|------|
| `push` | 上传（upload/init、file POST、upload/finish） |
| `pull` | 下载（manifest、file GET、latest） |
| `promote` | 发布 / 取消发布 |
| `admin` | 删除、同步存储、Token/Webhook 管理等 |

Token 可限定全局、项目或应用范围。`admin` 包含以上全部权限。

### 主要端点

**公开**

- `GET /api/v1/public/inventory` — 完整公开清单（**推荐**）
- `GET /api/v1/health` — 健康检查

已废弃（响应头 `Deprecation: true`）：`/public/projects*` 分页接口，请改用 `/public/inventory`。

**需认证**

- 项目 / 应用 / 版本：`/api/v1/projects/...`
- 上传：`POST /api/v1/upload/init`、`POST /api/v1/file/...`、`POST /api/v1/upload/finish`
- 下载：`GET /api/v1/manifest/...`、`GET /api/v1/file/...`（支持 Range）
- 发布：`POST /api/v1/publish`、`POST /api/v1/unpublish`（需 `promote`）
- 管理：`POST /api/v1/sync-storage`、Token/Webhook CRUD（需 `admin`）
- 审计：`GET /api/v1/audit-logs?operation=push&project_id=1&limit=50&offset=0`

完整文档见 Swagger UI。

## 开发

### 项目结构

```text
server/     Go API 服务
agent/      Go CLI（push / pull / push-tree）
web-ui/     React + TypeScript 管理界面
scripts/    Ruby 自动化脚本（构建、测试、批量推送）
```

### 本地命令

```bash
ruby scripts/up.rb              # 构建 + 启动
ruby scripts/build.rb --test    # 构建 + 启动 + 集成测试（与 CI 一致）
ruby scripts/build.rb --agent   # 仅编译 agent → /tmp/kkartifact-agent
ruby scripts/build.rb --server --no-cache --up

ruby scripts/test_integration.rb   # 单独跑集成测试（需服务已启动）
```

### 从源码构建

```bash
cd server && go build -o kkartifact-server .
cd agent  && go build -o kkartifact-agent .
cd web-ui && npm ci && npm run build
```

## CI/CD 与发布

```
本地                    GitHub                      生产
ruby scripts/up.rb  →   push/PR → CI            →   docker compose -f docker-compose.prod.yml up -d
ruby scripts/build.rb   tag v* → GHCR 镜像 + Release
     --test
```

推送标签触发构建：

```bash
git tag v1.0.0 && git push origin v1.0.0
```

| 镜像 | 地址 |
|------|------|
| Server（含 Agent） | `ghcr.io/kevin197011/kkartifact/server:<tag>` |
| Web UI | `ghcr.io/kevin197011/kkartifact/web-ui:<tag>` |

生产部署：

```bash
export KK_SERVER_IMAGE=ghcr.io/kevin197011/kkartifact/server:v1.0.0
export KK_WEB_UI_IMAGE=ghcr.io/kevin197011/kkartifact/web-ui:v1.0.0
docker compose -f docker-compose.prod.yml pull
docker compose -f docker-compose.prod.yml up -d
```

工作流细节见 [.github/workflows/README.md](.github/workflows/README.md)。

## 配置参考

### Server 环境变量

| 变量 | 默认 | 说明 |
|------|------|------|
| `SERVER_PORT` | 8080 | API 端口 |
| `DB_*` | 见 `.env.example` | PostgreSQL |
| `STORAGE_TYPE` | local | `local` 或 `s3` |
| `STORAGE_LOCAL_PATH` | /repos | 本地存储路径 |
| `ADMIN_USERNAME` / `ADMIN_PASSWORD` | admin / admin123 | 初始管理员 |
| `VERSION_RETENTION_LIMIT` | 5 | 全局版本保留数 |
| `JWT_SECRET` | 随机 | JWT 密钥 |
| `ENABLE_SWAGGER` | true | Swagger UI |

### Web UI

| 变量 | 默认 | 说明 |
|------|------|------|
| `WEB_UI_PORT` | 3000 | 前端端口 |
| `VITE_API_URL` | / | API 代理前缀 |

复制 `.env.example` 为 `.env` 并按环境修改。

## 常见问题

**上传报 `no buffer space available`**  
降低 `concurrency`：本地建议 30–50，远程 50–100。

**`ignore: []` 含义**  
不忽略任何文件。需要排除时添加 glob，如 `logs/`、`*.log`。

**存储与数据库不一致**  
管理后台 Projects 页点击「Sync Storage」，或 `POST /api/v1/sync-storage`（需 admin）。

**重置管理员密码**  
通过数据库更新 `users.password_hash`，或调整 `ADMIN_PASSWORD` 后重建管理员用户。

## 许可证

MIT License — Copyright (c) 2025 kk
