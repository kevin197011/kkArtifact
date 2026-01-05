# 发布系统需求文档（kkArtifact-server & kkArtifact-agent）

> 目标：构建一套 **替代 rsync + SSH** 的现代化发布代码托管与同步系统，支持 **多项目 / 多 App / Hash 版本管理**，通过 **Token 鉴权**实现安全、可审计、可扩展的代码发布与同步能力，同时支持 **Web 管理 UI 前后端分离**、**事件系统** 和 **Webhook**。

配置文件采用：

```text
.kkartifact.yml
```

---

## 一、背景与问题

### 1.1 当前痛点

* rsync + SSH 发布方式存在以下问题：

  * 无项目级权限隔离
  * 无统一版本元数据
  * 发布审计困难
  * 不适合云原生 / Agent 化部署
  * 安全模型粗糙（基于账号而非项目）
  * 缺乏可视化管理界面
  * 历史版本管理困难，容易占用大量存储
  * 为每个 App 单独配置版本保留数量繁琐

### 1.2 建设目标

* 建立 **以“编译产物（Artifact）为核心”** 的发布体系
* 实现类似 rsync 的 **多文件增量同步能力**
* 基于 HTTP + Token 的安全访问模型
* 支持 Pull / Push 双向同步
* 支持快速发布与秒级回滚
* **支持 API 查询**：列出项目、App、版本 Hash，并按创建时间排序
* **支持 Web 管理 UI 前后端分离**，提供可视化操作与状态监控
* **事件系统与 Webhook 支持**，实现发布流程通知和自动触发
* **发布动作不使用软链接**，基于部署路径的配置文件进行增删更新
* **Push 时版本号可作为参数传递**，由 Agent 生成 manifest 并推送到服务器
* **Server 端支持全局配置历史版本保留数量**，自动清理超出数量的旧版本
* **Server 定时任务每日凌晨 3 点执行版本清理**
* **Agent 历史版本控制仅针对当前操作的 App 生效，不影响其他 App**

---

## 二、总体架构

### 2.1 架构组件

* **kkArtifact-server**：

  * 代码托管中心
  * 提供 HTTP API
  * 管理项目 / App / 版本
  * 提供事件系统和 Webhook 管理
  * Web UI 后端服务

* **kkArtifact-agent**：

  * 部署在目标机器或构建节点
  * 执行 push / pull / diff / 校验

* **Build Server（可选）**：

  * 编译代码
  * 通过 Agent 或 API 上传构建产物

* **Web 管理 UI**：

  * 前端独立部署
  * 调用 API 展示项目 / 应用 / 版本信息
  * 支持发布操作、版本回滚、审计查询、事件订阅

---

## 三、核心数据模型

### 3.1 目录结构规范（kkArtifact-server）

```text
/repos
└── {project}
    └── {app}
        └── {hash}/
            ├── bin/
            ├── config/
            └── meta.yaml
```

### 3.2 版本不可变原则

* hash 版本一旦生成：

  * 不允许覆盖
  * 不允许修改
* 新版本 = 新 hash

---

## 四、meta.yaml 规范

```yaml
project: project-a
app: app-api
version: a8f3c21d
git_commit: a8f3c21d
build_time: 2025-12-26T17:30:00
builder: build-01
files:
  - path: bin/app
    sha256: xxx
    size: 123456
```

用途：

* 审计
* 校验
* 回滚定位
* **由 Agent 自动生成，无需手动维护**

---

## 五、Token 鉴权设计

### 5.1 Token 能力模型

Token 需绑定：

* project
* app（可多个）
* permissions
* 权限范围可分为 Global / Project / App

### 5.2 权限分级

| 权限级别          | 作用范围          |
| ------------- | ------------- |
| Global Token  | 管控所有项目 / 所有应用 |
| Project Token | 管控某一个项目下所有应用  |
| App Token     | 仅管控指定应用       |

### 5.3 权限校验逻辑

1. 校验 Token 合法性与过期
2. 判断权限范围（scope）
3. 校验操作权限（push / pull / promote）
4. Repo Server 决定是否允许访问

---

## 六、kkArtifact-server 功能需求

### 6.1 API 列表

| 方法     | 路径                                      | 描述                          |
| ------ | --------------------------------------- | --------------------------- |
| GET    | /manifest/{project}/{app}/{hash}        | 获取文件清单                      |
| GET    | /file/{project}/{app}/{hash}?path=      | 下载文件                        |
| PUT    | /file/{project}/{app}/{hash}?path=      | 上传单个文件                      |
| POST   | /upload/init                            | 初始化上传                       |
| POST   | /upload/finish                          | 完成上传                        |
| POST   | /promote                                | 标记为可发布版本                    |
| GET    | /projects                               | 列出所有项目，支持按创建时间排序            |
| GET    | /projects/{project}/apps                | 列出项目下所有 App，支持按创建时间排序       |
| GET    | /projects/{project}/apps/{app}/versions | 列出 App 下所有版本 Hash，支持按创建时间排序 |
| GET    | /webhooks                               | 查询已配置 Webhook               |
| POST   | /webhooks                               | 新建 Webhook                  |
| PUT    | /webhooks/{id}                          | 更新 Webhook                  |
| DELETE | /webhooks/{id}                          | 删除 Webhook                  |

### 6.2 Repo Server 职责边界

* 不关心 Git
* 不执行编译
* 只负责：

  * 文件托管
  * 版本隔离
  * 权限校验
  * API 查询项目 / App / 版本信息
  * 事件触发与 Webhook 通知

### 6.3 历史版本管理（Server）

* **全局配置**保留最新 N 个版本（无需单独设置每个 App）
* 定时任务每日凌晨 3 点执行，清理超出版本
* 自动删除或归档最旧版本

---

## 七、kkArtifact-agent 功能需求

### 7.0 忽略文件 / 目录规则（Ignore Rules）

* 配置文件统一为 `.kkartifact.yml`
* push / pull / diff / manifest 阶段均生效
* **在 Agent 执行 push 或 pull 时必须加载 `.kkartifact.yml`**，若不存在则报错并退出，防止意外同步
* 支持 glob / 前缀 / 目录 / 文件

示例：

```yaml
ignore:
  - logs/
  - tmp/
  - '*.log'
```

### 7.1 Pull（拉取）能力

* 按 project/app/hash 拉取版本
* manifest diff 本地与远端文件
* 仅拉取缺失 / 变更文件
* 支持并发下载
* 支持 ignore 规则
* **不做软链接操作**，基于部署目录配置文件做更新
* **自动删除新版本中不存在的旧文件**（清理旧版本冗余文件）
* **部署端本地目录结构**：所有代码文件直接在 `/opt/apps/{project}/{app}/` 下

### 7.2 Push（推送）能力

* 扫描本地构建目录
* **支持传入版本号参数**，由 Agent 生成 manifest 并推送
* 多文件并发上传
* 支持幂等与断点重试
* 支持 ignore 规则
* **推送端本地目录结构**：构建产物全部在 `/local/build/path/` 下
* **不做软链接操作**，更新部署路径下的文件即可

### 7.3 历史版本控制（Agent）

* **仅针对当前操作的 App 生效，不影响其他 App**
* Push 后检查当前 App 的版本数量
* 超过全局保留数量时，删除最旧的版本
* Pull 时只同步指定版本
* 遵循 Server 全局保留策略，但作用范围仅限当前 App

### 7.4 Agent 使用示例

#### Push 使用方式

```bash
kkArtifact-agent push \
  --project myproject \
  --app myapp \
  --version 1.0.0 \
  --path /local/build/path \
  --config .kkartifact.yml
```

#### Pull 使用方式

```bash
kkArtifact-agent pull \
  --project myproject \
  --app myapp \
  --version 1.0.0 \
  --path /opt/apps/myproject/myapp \
  --config .kkartifact.yml
```

### 7.5 Agent 配置示例

```yaml
server_url: https://repo.example.com
token: xxx
concurrency: 8
chunk_size: 4MB
retain_versions: 5  # 全局保留最新 5 个版本
ignore:
  - logs/
  - tmp/
  - '*.log'
```

---

## 八、Web 管理 UI 需求

* 前后端分离：

  * **前端**：Vue / React / Svelte，可独立部署
  * **后端**：调用 kkArtifact-server API
* 功能：

  * 项目 / App / 版本浏览
  * Push / Pull 操作触发
  * 回滚操作
  * 发布状态和日志展示
  * Webhook / 事件订阅管理
  * Token 管理

---

## 九、事件系统 & Webhook

### 9.1 事件触发场景

* Push 完成
* Pull 完成
* Promote（版本可发布）
* 回滚完成
* 删除版本

### 9.2 Webhook 功能

* 支持配置不同事件的 Webhook
* 支持 HTTP POST 回调
* 支持 Slack / 钉钉 / CI/CD / 内部服务
* 用途：

  * 发布通知
  * 自动部署触发
  * 回滚触发
  * 审计与日志同步

### 9.3 Webhook 配置示例

```yaml
webhooks:
  - name: slack_notify
    event: promote
    url: https://hooks.slack.com/services/xxx/yyy/zzz
  - name: ci_trigger
    event: push
    url: https://ci.example.com/build
```

---

## 十、同步与性能要求

### 10.1 基本性能要求

* 支持大文件（>1GB）
* 支持 HTTP Range
* 并发 worker pool
* 失败可重试
* 操作幂等

### 10.2 大规模场景支持（V1 核心目标）

系统需支持以下大规模场景：

* **模块规模**：支持 2000+ 个 App 模块
* **存储容量**：支持 2TB+ 存储容量
* **并发访问**：支持 200+ 并发操作
* **查询性能**：
  * 项目列表查询：<500ms（2000+ 项目）
  * App 列表查询：<500ms（单项目 100+ App）
  * 版本列表查询：<1s（单 App 100+ 版本）
* **存储性能**：
  * 文件上传：支持 100+ MB/s 吞吐
  * 文件下载：支持 100+ MB/s 吞吐
  * 支持对象存储（S3/OSS 兼容）作为存储后端

### 10.3 大规模场景优化策略

* **存储层优化**：
  * 支持对象存储（S3/OSS）作为主要存储后端
  * 本地文件系统作为可选方案（开发/小规模）
  * 存储接口抽象，支持多种存储后端切换
* **数据库优化**：
  * 关键字段索引优化（项目、App、版本创建时间）
  * 查询分页支持（默认每页 50 条，最大 500 条）
  * 数据库连接池优化（最小 10，最大 100）
  * 慢查询监控和优化
* **缓存机制**：
  * Redis 缓存热点数据（项目列表、App 列表、最新版本）
  * 缓存失效策略（TTL + 主动失效）
  * Manifest 元数据缓存
* **API 优化**：
  * 响应压缩（Gzip）默认启用
  * 批量查询接口（支持一次查询多个 App 的版本）
  * 异步操作支持（大规模清理任务异步执行）
* **前端优化**：
  * 虚拟滚动（支持渲染大量列表项）
  * 分页加载（前端分页 + 后端分页）
  * 数据预取和缓存
  * 懒加载（按需加载详细信息）

---

## 十一、安全与审计

* 所有操作基于 Token
* Token 可吊销
* 每次 push / pull 记录：

  * 时间
  * project/app
  * hash
  * agent id
* Webhook 事件通知可记录历史触发日志

---

## 十二、开发与调测环境

### 12.1 Docker Compose 调测环境

系统应提供基于 Docker Compose 的完整开发与调测环境，支持：

* **一键启动**：通过 `docker compose up` 启动所有服务
* **服务包含**：
  * kkArtifact-server（HTTP API 服务）
  * PostgreSQL 数据库
  * Web UI 前端服务
  * （可选）Agent 测试环境
* **数据持久化**：
  * PostgreSQL 数据卷持久化
  * Artifact 存储目录挂载
  * 日志目录挂载
* **环境配置**：
  * 通过 `.env` 文件配置环境变量
  * 支持开发/测试/生产环境切换
  * 数据库连接、端口、存储路径等可配置
* **热重载**：
  * 开发模式下支持代码热重载
  * 前端支持 Vite HMR
  * 后端支持文件监控自动重启
* **日志查看**：
  * 统一日志输出到控制台
  * 支持 `docker compose logs -f` 查看实时日志
  * 支持按服务过滤日志

### 12.2 Docker Compose 配置要求

* `docker-compose.yml` 配置文件应包含：
  * Server 服务定义（端口映射、环境变量、卷挂载）
  * PostgreSQL 服务定义（数据卷、初始化脚本）
  * Web UI 服务定义（开发服务器或构建产物）
  * 网络配置（服务间通信）
* `.env.example` 提供环境变量模板
* `README.md` 包含快速启动说明

### 12.3 开发环境使用方式

```bash
# 启动所有服务
docker compose up -d

# 查看日志
docker compose logs -f

# 停止服务
docker compose down

# 重建并启动
docker compose up -d --build

# 清理所有数据（包括数据库）
docker compose down -v
```

---

## 十三、非功能性需求

### 13.1 可扩展性

* kkArtifact-server 无状态（可横向扩展）
* 支持水平扩展（多实例负载均衡）
* 存储支持对象存储（S3/OSS）以实现分布式存储
* 数据库支持读写分离（未来扩展）

### 13.2 存储后端支持

* **V1 必须支持**：
  * 对象存储（S3 兼容协议，如 AWS S3、阿里云 OSS、MinIO）
  * 本地文件系统（开发/测试环境）
* **可选支持**：
  * NFS（网络文件系统）
  * 其他对象存储后端（通过插件机制）

### 13.3 性能优化（V1 核心要求）

* **必须实现**：
  * 数据库索引优化（关键查询字段）
  * API 响应压缩（Gzip）
  * 查询分页支持（所有列表接口）
  * Redis 缓存机制（热点数据）
  * 性能监控和日志
* **建议实现**：
  * SHA256 并行计算（大文件优化）
  * 异步任务处理（大规模清理）

### 13.4 运维支持

* kkArtifact-agent 支持 systemd 运行
* API 兼容未来 Web UI / CI 平台
* **支持 Docker Compose 一键启动开发环境**
* 支持 Prometheus 指标导出
* 支持结构化日志（JSON 格式）

---

## 十四、成功标准（Definition of Done）

* 能完全替代 rsync 发布
* 支持多项目隔离
* 支持秒级回滚
* Agent 可独立运行
* 无 SSH 依赖
* 支持通过 API 查询项目 / App / 版本 Hash，并按创建时间排序
* Agent push / pull 必须存在 `.kkartifact.yml` 配置，否则报错退出
* **发布动作不使用软链接，基于部署路径配置文件进行增删更新**
* **Pull 时自动删除新版本中不存在的旧文件**
* **Push 时版本号可作为参数传递**
* **Web 管理 UI 可独立操作项目 / 版本 / 发布流程**
* **事件系统与 Webhook 可完整通知和触发外部系统**
* **Server 全局配置历史版本数量，定时任务每日凌晨 3 点执行**
* **Agent 历史版本控制仅针对当前操作的 App 生效，不影响其他 App**

---

## 十五、后续规划（非本期）

* 灰度 / 批次发布
* CI/CD 深度集成
* K8s / systemd 自动化部署
* 高级审计和报表

---

> 本文档定位：**kkArtifact V1 核心需求基线，包含 Agent、Server、Web UI 与事件系统**
