# GitHub Actions CI/CD 工作流说明

本项目使用 GitFlow 模型，并配置了完整的 GitHub Actions CI/CD 流水线，支持自动构建和推送 Docker 镜像到 GitHub Container Registry (ghcr.io)。

## 工作流概览

### 1. 主分支 CI/CD (`ci-cd.yml`)

**触发条件:**
- 推送到 `main` 分支
- 向 `main` 或 `develop` 分支的 Pull Request

**功能:**
- ✅ 运行单元测试和代码检查
- 🐳 **仅 main 分支**: 构建并推送 Docker 镜像 (AMD64 + ARM64)
- 🔒 **仅 main 分支**: 安全扫描 (Trivy)
- 📦 **仅 main 分支**: 推送到 GitHub Container Registry

**镜像标签 (仅 main 分支):**
- `ghcr.io/用户名/higress-graphql-federation:main`
- `ghcr.io/用户名/higress-graphql-federation:latest`

### 1.1. Develop 分支 CI (`develop-ci.yml`)

**触发条件:**
- 推送到 `develop` 分支
- 向 `develop` 分支的 Pull Request

**功能:**
- ✅ 运行单元测试和代码检查
- 📊 代码覆盖率检查
- 📋 上传覆盖率报告

### 2. Feature 分支 CI (`feature-ci.yml`)

**触发条件:**
- 推送到 `feature/**` 或 `feat/**` 分支
- 向 `develop` 或 `main` 分支的 Pull Request

**功能:**
- ✅ 运行单元测试和代码检查
- 📊 代码覆盖率检查
- 🐳 构建测试镜像
- 🧪 集成测试

**镜像标签:**
- `ghcr.io/用户名/higress-graphql-federation:feature-分支名`
- `ghcr.io/用户名/higress-graphql-federation:分支名-commit-hash`

### 3. Release 分支发布 (`release.yml`)

**触发条件:**
- 推送到 `release/**` 分支
- 手动触发 (workflow_dispatch)

**功能:**
- ✅ 运行完整的测试套件
- 🔒 安全扫描
- 🐳 构建并推送发布镜像
- 📝 自动生成 GitHub Release
- 📋 生成变更日志

**镜像标签:**
- `ghcr.io/用户名/higress-graphql-federation:v1.0.0`
- `ghcr.io/用户名/higress-graphql-federation:v1.0.0-amd64`
- `ghcr.io/用户名/higress-graphql-federation:v1.0.0-arm64`

### 4. Hotfix 分支紧急修复 (`hotfix.yml`)

**触发条件:**
- 推送到 `hotfix/**` 分支
- 手动触发 (workflow_dispatch)

**功能:**
- ✅ 运行完整的测试套件
- 🔒 安全扫描
- 🐳 构建并推送热修复镜像
- 📝 创建紧急发布
- 🔔 团队通知

**镜像标签:**
- `ghcr.io/用户名/higress-graphql-federation:v1.0.1`
- `ghcr.io/用户名/higress-graphql-federation:v1.0.1-amd64`
- `ghcr.io/用户名/higress-graphql-federation:v1.0.1-arm64`

### 5. 依赖更新 (`dependencies.yml`)

**触发条件:**
- 手动触发 (workflow_dispatch)

**功能:**
- 🔄 更新 Go 模块依赖 (需要人工确认)
- 📝 可选择是否自动创建 Pull Request
- ⚠️ **需要人工审核**: 所有依赖更新都需要人工确认

## GitFlow 工作流

### 开发流程

1. **Feature 开发**
   ```bash
   git checkout -b feature/new-feature develop
   # 开发功能...
   git push origin feature/new-feature
   # 创建 Pull Request 到 develop
   ```

2. **Release 准备**
   ```bash
   git checkout -b release/v1.0.0 develop
   # 准备发布...
   git push origin release/v1.0.0
   # 自动触发发布流程
   ```

3. **Hotfix 紧急修复**
   ```bash
   git checkout -b hotfix/v1.0.1 main
   # 修复问题...
   git push origin hotfix/v1.0.1
   # 自动触发热修复流程
   ```

### 分支策略

- **`main`**: 生产环境代码，稳定版本
- **`develop`**: 开发主分支，集成最新功能
- **`feature/*`**: 功能开发分支
- **`release/*`**: 发布准备分支
- **`hotfix/*`**: 紧急修复分支

## Docker 镜像使用

### 拉取镜像

```bash
# 拉取最新版本
docker pull ghcr.io/用户名/higress-graphql-federation:latest

# 拉取特定版本
docker pull ghcr.io/用户名/higress-graphql-federation:v1.0.0

# 拉取特定架构
docker pull ghcr.io/用户名/higress-graphql-federation:v1.0.0-amd64
docker pull ghcr.io/用户名/higress-graphql-federation:v1.0.0-arm64
```

### 使用镜像

```bash
# 运行容器
docker run -d \
  --name higress-graphql-federation \
  -p 8080:8080 \
  ghcr.io/用户名/higress-graphql-federation:latest
```

## 权限配置

### GitHub 仓库权限

确保以下权限已正确配置：

1. **Actions**: 启用 GitHub Actions
2. **Packages**: 启用 GitHub Packages (Container Registry)
3. **Secrets**: 配置必要的 secrets

### 必要的 Secrets

- `GITHUB_TOKEN`: 自动提供，用于推送镜像和创建 Release

### 权限说明

工作流需要以下权限：
- `contents: read/write`: 读取代码，创建 Release
- `packages: write`: 推送 Docker 镜像 (仅 main 分支)
- `security-events: write`: 上传安全扫描结果 (仅 main 分支)

## 重要变更说明

### 🔄 镜像构建策略调整

- **只有合并到 `main` 分支才会构建和推送 Docker 镜像**
- `develop` 分支只运行测试，不构建镜像
- 这样可以减少不必要的镜像构建，节省资源

### ⚠️ 依赖更新策略调整

- **所有依赖更新都需要人工确认**
- Dependabot 不会自动创建 Pull Request
- 需要手动触发依赖更新工作流
- 可以选择是否自动创建 PR 或仅显示变更摘要

## 监控和通知

### 工作流状态

- ✅ 绿色: 所有检查通过
- ❌ 红色: 测试失败或构建错误
- 🟡 黄色: 工作流进行中

### 通知方式

- GitHub 通知
- Pull Request 状态检查
- Release 创建通知

## 故障排除

### 常见问题

1. **构建失败**
   - 检查 Go 版本兼容性
   - 验证 Dockerfile 语法
   - 检查依赖项

2. **推送失败**
   - 验证 GitHub Token 权限
   - 检查 Container Registry 访问权限

3. **测试失败**
   - 检查测试代码
   - 验证环境配置

### 调试步骤

1. 查看工作流日志
2. 检查分支权限
3. 验证 Secrets 配置
4. 测试本地构建

## 最佳实践

1. **分支命名**: 遵循 GitFlow 约定
2. **提交信息**: 使用语义化提交信息
3. **测试覆盖**: 确保充分的测试覆盖
4. **安全扫描**: 定期检查安全漏洞 (仅 main 分支)
5. **依赖更新**: 定期手动检查和更新依赖项
6. **镜像管理**: 只在 main 分支构建镜像，避免资源浪费
7. **人工审核**: 所有依赖更新都需要人工确认

## 相关链接

- [GitHub Actions 文档](https://docs.github.com/en/actions)
- [GitHub Container Registry](https://docs.github.com/en/packages/working-with-a-github-packages-registry/working-with-the-container-registry)
- [GitFlow 模型](https://nvie.com/posts/a-successful-git-branching-model/)
- [Docker 多架构构建](https://docs.docker.com/buildx/working-with-buildx/)
