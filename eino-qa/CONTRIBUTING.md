# 贡献指南

感谢您对 Eino QA System 的关注！我们欢迎各种形式的贡献。

## 如何贡献

### 报告问题

如果您发现了 bug 或有功能建议：

1. 在 [Issues](https://github.com/your-repo/issues) 中搜索，确保问题未被报告
2. 创建新的 Issue，提供详细信息：
   - Bug 报告：复现步骤、预期行为、实际行为、环境信息
   - 功能建议：使用场景、预期效果、可能的实现方案

### 提交代码

1. **Fork 项目**
   ```bash
   # 在 GitHub 上 Fork 项目
   git clone https://github.com/your-username/eino-qa.git
   cd eino-qa
   ```

2. **创建分支**
   ```bash
   git checkout -b feature/your-feature-name
   # 或
   git checkout -b fix/your-bug-fix
   ```

3. **开发和测试**
   ```bash
   # 安装依赖
   make deps
   
   # 运行测试
   make test
   
   # 代码格式化
   make fmt
   
   # 代码检查
   make lint
   ```

4. **提交更改**
   ```bash
   git add .
   git commit -m "feat: add new feature"
   # 或
   git commit -m "fix: fix bug description"
   ```

5. **推送到 GitHub**
   ```bash
   git push origin feature/your-feature-name
   ```

6. **创建 Pull Request**
   - 在 GitHub 上创建 Pull Request
   - 填写 PR 模板，描述您的更改
   - 等待代码审查

## 开发规范

### 代码风格

- 遵循 [Uber Go Style Guide](https://github.com/uber-go/guide/blob/master/style.md)
- 使用 `gofmt` 格式化代码
- 使用 `golangci-lint` 进行代码检查
- 保持代码简洁、可读

### 提交信息规范

使用 [Conventional Commits](https://www.conventionalcommits.org/) 规范：

```
<type>(<scope>): <subject>

<body>

<footer>
```

**类型 (type)**:
- `feat`: 新功能
- `fix`: Bug 修复
- `docs`: 文档更新
- `style`: 代码格式（不影响代码运行）
- `refactor`: 重构（既不是新功能也不是 Bug 修复）
- `perf`: 性能优化
- `test`: 测试相关
- `chore`: 构建过程或辅助工具的变动

**示例**:
```
feat(chat): add streaming response support

Implement Server-Sent Events (SSE) for streaming chat responses.
This improves user experience by showing incremental results.

Closes #123
```

### 测试要求

- 新功能必须包含单元测试
- Bug 修复必须包含回归测试
- 测试覆盖率不低于 80%
- 所有测试必须通过

```bash
# 运行测试
make test

# 查看覆盖率
make test-coverage
```

### 文档要求

- 新功能需要更新相关文档
- API 变更需要更新 API 文档
- 重要变更需要更新 CHANGELOG.md

## 项目结构

```
eino-qa/
├── cmd/                    # 应用入口
├── internal/               # 内部包
│   ├── domain/            # 领域层
│   ├── usecase/           # 用例层
│   ├── adapter/           # 适配层
│   └── infrastructure/    # 基础设施层
├── pkg/                   # 公共包
├── config/                # 配置文件
├── docs/                  # 文档
└── examples/              # 示例代码
```

详细说明请参考 [PROJECT_STRUCTURE.md](PROJECT_STRUCTURE.md)。

## 开发环境设置

### 前置要求

- Go 1.23+
- Docker 和 Docker Compose
- Make
- Git

### 设置步骤

```bash
# 1. 克隆项目
git clone https://github.com/your-username/eino-qa.git
cd eino-qa

# 2. 安装依赖
make deps

# 3. 配置环境变量
cp .env.example .env
# 编辑 .env 文件

# 4. 启动 Milvus
make milvus-up

# 5. 运行测试
make test

# 6. 运行服务
make run
```

## 代码审查流程

1. **自动检查**
   - 代码格式检查
   - 单元测试
   - 代码覆盖率
   - 静态分析

2. **人工审查**
   - 代码质量
   - 设计合理性
   - 测试完整性
   - 文档完整性

3. **合并要求**
   - 所有自动检查通过
   - 至少一位维护者批准
   - 无未解决的评论
   - 与主分支无冲突

## 发布流程

1. 更新版本号（遵循语义化版本）
2. 更新 CHANGELOG.md
3. 创建 Git tag
4. 构建和测试
5. 发布到 GitHub Releases
6. 更新文档

## 社区准则

### 行为准则

- 尊重他人，保持友善
- 欢迎新手，耐心解答
- 建设性的批评和建议
- 专注于技术讨论

### 沟通渠道

- GitHub Issues - Bug 报告和功能建议
- GitHub Discussions - 一般讨论和问答
- Pull Requests - 代码审查和讨论

## 常见问题

### Q: 如何选择合适的 Issue？

A: 查找标签为 `good first issue` 或 `help wanted` 的 Issue。

### Q: 我的 PR 多久会被审查？

A: 通常在 1-3 个工作日内。如果超过一周未响应，可以在 PR 中评论提醒。

### Q: 如何运行特定的测试？

A: 使用 `go test -run TestName ./path/to/package`

### Q: 代码风格检查失败怎么办？

A: 运行 `make fmt` 自动格式化代码，然后运行 `make lint` 检查。

## 获取帮助

如有任何问题，请：
1. 查看 [文档](docs/)
2. 搜索 [Issues](https://github.com/your-repo/issues)
3. 创建新的 Issue
4. 发送邮件至：dev@example.com

## 致谢

感谢所有贡献者的付出！

---

**维护团队**: Eino QA Team  
**最后更新**: 2024-11-29
