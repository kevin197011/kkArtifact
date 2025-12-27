# kkArtifact

> 现代化代码发布与同步系统 - 替代 rsync + SSH 的解决方案

## 项目简介

kkArtifact 是一套现代化的代码托管与同步系统，支持多项目/多 App/Hash 版本管理，通过 Token 鉴权实现安全、可审计、可扩展的代码发布与同步能力。

### 核心功能

- **kkArtifact-server**: HTTP API 服务，提供文件托管、版本管理、Web UI 后端
- **kkArtifact-agent**: CLI 工具，支持 push/pull 操作和 manifest 生成
- **Web UI**: 基于 React + Ant Design 的管理界面（规划中）
- **事件系统**: 支持 Webhook 的事件驱动架构
- **大规模支持**: 支持 2000+ App，2TB+ 存储容量（S3/OSS 对象存储）
- **版本管理**: 不可变版本、Hash 版本控制、自动版本清理

## 快速开始

### 前置要求

- Docker 20.10+
- Docker Compose 2.0+
- Go 1.21+ (本地开发)
- Node.js 18+ (Web UI 开发，可选)

### 使用 Docker Compose

1. 复制环境变量配置：
```bash
cp .env.example .env
```

2. 启动所有服务：
```bash
docker compose up -d
```

3. 运行数据库迁移（首次启动）：
```bash
docker compose exec server ./migrate -direction=up
```

4. 查看日志：
```bash
docker compose logs -f
```

5. 停止服务：
```bash
docker compose down
```

### 本地开发

#### Server 开发

```bash
cd server
go mod download
go run main.go
```

#### Agent 开发

```bash
cd agent
go mod download
go run main.go --help
```

#### 数据库迁移

```bash
cd server
export MIGRATIONS_PATH=./migrations/migrations
go run cmd/migrate/main.go -direction=up
```

## 项目结构

```
.
├── server/          # kkArtifact-server (Go)
│   ├── internal/    # 内部包
│   │   ├── api/     # HTTP API handlers
│   │   ├── auth/    # 认证授权
│   │   ├── cache/   # Redis 缓存
│   │   ├── config/  # 配置管理
│   │   ├── database/# 数据库层
│   │   ├── events/  # 事件系统
│   │   ├── storage/ # 存储层
│   │   └── ...
│   └── migrations/  # 数据库迁移
├── agent/           # kkArtifact-agent (Go)
│   └── internal/
│       ├── cli/     # CLI 命令
│       ├── client/  # API 客户端
│       ├── config/  # 配置管理
│       └── manifest/# Manifest 生成
├── web-ui/          # Web Management UI (React + TypeScript，规划中)
├── openspec/        # OpenSpec 规范和提案
├── docker-compose.yml
├── Makefile
└── README.md
```

## 配置说明

主要配置通过环境变量管理，详见 `.env.example`：

- **Server**: 端口、主机配置
- **Database**: PostgreSQL 连接配置
- **Redis**: 缓存配置（可选）
- **Storage**: 存储后端配置（S3/OSS 或本地文件系统）

### Agent 配置

创建 `.kkartifact.yml` 配置文件：

```yaml
server_url: http://localhost:8080
token: your-token-here
project: myproject
app: myapp
ignore:
  - "*.log"
  - ".git/**"
  - "node_modules/**"
retain_versions: 5  # 可选，参考 server 全局配置
```

## 使用示例

### Push 操作

```bash
# 使用命令行参数
kkartifact-agent push \
  --project myproject \
  --app myapp \
  --version a8f3c21d \
  --path ./dist

# 或使用配置文件
kkartifact-agent push --version a8f3c21d --path ./dist
```

### Pull 操作

```bash
kkartifact-agent pull \
  --project myproject \
  --app myapp \
  --version a8f3c21d \
  --deploy-path /opt/myapp
```

## API 端点

### 认证
所有 API 端点需要 Bearer Token 认证：
```
Authorization: Bearer <token>
```

### 主要端点

- `GET /api/v1/projects` - 列出所有项目
- `GET /api/v1/projects/:project/apps` - 列出项目的所有 App
- `GET /api/v1/projects/:project/apps/:app/versions` - 列出 App 的所有版本
- `GET /api/v1/manifest/:project/:app/:hash` - 获取版本 manifest
- `POST /api/v1/upload/init` - 初始化上传会话
- `POST /api/v1/file/:project/:app/:hash` - 上传文件
- `POST /api/v1/upload/finish` - 完成上传
- `POST /api/v1/promote` - 标记版本为已发布
- `GET /api/v1/config` - 获取全局配置
- `PUT /api/v1/config` - 更新全局配置
- `GET /api/v1/webhooks` - 列出所有 webhook
- `POST /api/v1/webhooks` - 创建 webhook
- `GET /api/v1/audit-logs` - 列出审计日志

## 构建

使用 Makefile：

```bash
make build-all      # 构建所有组件
make build-server   # 仅构建 server
make build-agent    # 仅构建 agent
make test           # 运行测试
```

## 开发计划

本项目采用 OpenSpec 规范进行开发，详见 `openspec/changes/implement-kkartifact-system/`。

### 当前状态

✅ **已完成**：
- 核心存储层（S3/本地文件系统）
- 数据库层和元数据管理
- HTTP API 层（主要端点）
- 认证和授权框架
- Webhook 管理
- Agent CLI（push/pull）
- 配置管理
- 版本清理和定时任务
- Gzip 压缩和 CORS 支持

🚧 **进行中**：
- Web UI 前端实现
- 更完善的错误处理和日志
- 性能优化和缓存集成

## License

MIT License - Copyright (c) 2025 kk
