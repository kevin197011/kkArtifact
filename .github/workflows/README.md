# GitHub Actions 工作流说明

## build-and-release.yml

此工作流用于构建和发布 kkArtifact 项目的所有组件。

### 功能

1. **构建并推送 Docker 镜像**
   - Server 镜像：推送到 `docker.io/<IMAGE_PREFIX>/server`
   - Web UI 镜像：推送到 `docker.io/<IMAGE_PREFIX>/web-ui`

2. **构建 Agent 二进制文件**
   - 支持多平台构建（Linux、macOS、Windows）
   - 支持多架构（amd64、arm64）

3. **发布到 GitHub Releases**
   - 自动创建 Release
   - 上传所有平台的二进制文件
   - 生成 SHA256 校验和文件

### 触发条件

- **自动触发**：推送版本标签（格式：`v*`，如 `v1.0.0`）
- **手动触发**：在 GitHub Actions 页面手动运行工作流

### 配置要求

#### Docker Hub 认证

在 GitHub 仓库的 Settings > Secrets and variables > Actions 中添加以下密钥：

- `DOCKER_USERNAME`: Docker Hub 用户名
- `DOCKER_PASSWORD`: Docker Hub 密码或访问令牌

#### 镜像仓库前缀

修改工作流文件中的 `IMAGE_PREFIX` 环境变量：

```yaml
env:
  IMAGE_PREFIX: your-dockerhub-username  # 修改为你的 Docker Hub 用户名或组织名
```

### 使用示例

#### 创建并推送版本标签

```bash
# 创建版本标签
git tag -a v1.0.0 -m "Release version 1.0.0"

# 推送标签到远程仓库
git push origin v1.0.0
```

推送标签后，GitHub Actions 将自动：
1. 构建并推送 Docker 镜像（带版本标签）
2. 构建所有平台的 Agent 二进制文件
3. 创建 GitHub Release 并上传所有文件

#### 手动触发

1. 访问 GitHub 仓库的 Actions 页面
2. 选择 "Build and Release" 工作流
3. 点击 "Run workflow" 按钮
4. 选择分支并点击 "Run workflow"

### 生成的标签

#### Docker 镜像标签

- `latest`（仅主分支）
- `v1.0.0`（版本标签）
- `1.0`（主版本.次版本）
- `1`（主版本）
- `<commit-sha>`（提交 SHA）

#### GitHub Release

- Release 名称：标签名称（如 `v1.0.0`）
- 包含文件：
  - `kkartifact-agent-linux-amd64`
  - `kkartifact-agent-linux-arm64`
  - `kkartifact-agent-darwin-amd64`
  - `kkartifact-agent-darwin-arm64`
  - `kkartifact-agent-windows-amd64.exe`
  - `checksums.txt`（SHA256 校验和）

### 注意事项

1. **首次使用**：确保已配置 Docker Hub 认证密钥
2. **镜像仓库**：修改 `IMAGE_PREFIX` 以匹配你的 Docker Hub 用户名/组织
3. **版本标签**：使用语义化版本（SemVer）格式（如 `v1.0.0`）
4. **预发布版本**：如果标签包含 `-`（如 `v1.0.0-beta.1`），Release 将被标记为预发布版本

