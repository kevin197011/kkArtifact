# GitHub Actions 工作流说明

## 概览

| 工作流 | 触发 | 作用 |
|--------|------|------|
| **CI** (`ci.yml`) | push/PR → `main` | 单元测试 + Web UI 构建 + Docker 集成测试 |
| **Build and Release** (`build-and-release.yml`) | 推送标签 `v*` / 手动 | 构建并推送 GHCR 镜像，创建 Release |

## 镜像地址（GHCR）

统一命名（小写）：

```
ghcr.io/<owner>/kkartifact/server:<tag>
ghcr.io/<owner>/kkartifact/web-ui:<tag>
```

示例：

```bash
docker pull ghcr.io/kevin197011/kkartifact/server:v1.0.0
docker pull ghcr.io/kevin197011/kkartifact/web-ui:v1.0.0
```

Server 镜像内已包含多平台 Agent 二进制与安装脚本。

## 发布流程（GitHub → 生产）

```bash
# 1. 打标签触发构建
git tag v1.0.0
git push origin v1.0.0

# 2. Actions 自动：构建镜像 → 推送 GHCR → 创建 Release

# 3. 生产服务器拉取并启动
export KK_SERVER_IMAGE=ghcr.io/kevin197011/kkartifact/server:v1.0.0
export KK_WEB_UI_IMAGE=ghcr.io/kevin197011/kkartifact/web-ui:v1.0.0
docker compose -f docker-compose.prod.yml pull
docker compose -f docker-compose.prod.yml up -d
```

手动触发：Actions → **Build and Release** → Run workflow

## 本地开发（线下，不依赖 GitHub）

```bash
# 一键构建并启动
ruby scripts/up.rb

# 构建 + 启动 + 集成测试
ruby scripts/build.rb --test

# 仅编译 agent
ruby scripts/build.rb --agent

# 强制无缓存重建 server
ruby scripts/build.rb --server --no-cache --up
```

等价 Docker 命令：

```bash
docker compose build server web-ui
docker compose up -d
ruby scripts/test_integration.rb
```

## CI 与本地测试对应关系

| 本地命令 | CI 对应 |
|----------|---------|
| `cd server && go test ./...` | CI → Unit Tests (server) |
| `cd agent && go test ./...` | CI → Unit Tests (agent) |
| `cd web-ui && npm run build` | CI → Web UI |
| `ruby scripts/build.rb --test` | CI → Integration |
