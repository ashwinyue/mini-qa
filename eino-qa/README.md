# Eino QA System

基于 CloudWeGo Eino 框架构建的 Go 语言智能客服系统，采用 Clean Architecture 设计模式，提供高性能、可扩展的智能对话服务。

## ✨ 特性

- 🚀 **高性能**：基于 Go 语言和 Eino 框架，支持高并发处理
- 🎯 **智能路由**：自动识别用户意图（课程咨询、订单查询、直接回答、人工转接）
- 📚 **RAG 检索**：基于 Milvus 向量数据库的知识库检索增强生成
- 🔐 **多租户**：支持租户级数据隔离，独立的向量 Collection 和数据库
- 🛡️ **安全**：敏感信息脱敏、API Key 验证、SQL 注入防护
- 📊 **可观测**：结构化日志、指标统计、健康检查
- 🔄 **流式响应**：支持 SSE 流式输出，提升用户体验
- 🏗️ **Clean Architecture**：清晰的分层架构，易于维护和扩展

## 🛠️ 技术栈

- **Web 框架**: Gin - 高性能 HTTP 路由
- **ORM**: GORM - 类型安全的数据库操作
- **关系数据库**: SQLite - 轻量级嵌入式数据库
- **向量数据库**: Milvus - 高性能向量检索
- **AI 框架**: Eino ADK + Compose - 智能体编排
- **LLM 服务**: DashScope (通义千问) - 聊天和嵌入模型

## 项目结构

```
eino-qa/
├── cmd/
│   └── server/          # 应用入口
├── internal/
│   ├── domain/          # 领域层（实体、仓储接口）
│   ├── usecase/         # 用例层（业务逻辑）
│   ├── adapter/         # 接口适配层（HTTP Handler、中间件）
│   └── infrastructure/  # 基础设施层（数据库、AI 组件）
│       ├── config/      # 配置管理
│       ├── logger/      # 日志系统
│       ├── ai/          # AI 组件（Eino）
│       └── repository/  # 数据仓储实现
│           └── milvus/  # ✅ Milvus 向量数据库集成
├── pkg/                 # 公共工具包
├── config/              # 配置文件
├── examples/            # 使用示例
└── go.mod
```

## 🚀 快速开始

### 前置要求

- Go 1.23+
- Docker 和 Docker Compose
- DashScope API Key（[获取地址](https://dashscope.console.aliyun.com/)）

### 5 分钟快速启动

```bash
# 1. 克隆项目
git clone <repository-url>
cd eino-qa

# 2. 配置环境变量
cp .env.example .env
# 编辑 .env 文件，填入你的 DashScope API Key

# 3. 启动 Milvus
docker-compose -f docker-compose.milvus.yml up -d

# 4. 等待 Milvus 启动（约 30 秒）
sleep 30

# 5. 运行服务
make run
```

服务将在 `http://localhost:8080` 启动。

### 验证安装

```bash
# 健康检查
curl http://localhost:8080/health

# 测试对话
curl -X POST http://localhost:8080/chat \
  -H "Content-Type: application/json" \
  -d '{"query": "你好"}'
```

详细的启动指南请参考 [QUICKSTART.md](QUICKSTART.md) 或 [STARTUP_GUIDE.md](STARTUP_GUIDE.md)。

## ⚙️ 配置

配置文件位于 `config/config.yaml`，支持通过环境变量覆盖配置项。

### 主要配置项

```yaml
server:
  port: 8080              # HTTP 服务端口
  mode: debug             # 运行模式：debug, release

dashscope:
  api_key: ${DASHSCOPE_API_KEY}  # API Key（从环境变量读取）
  chat_model: qwen-turbo          # 聊天模型
  embed_model: text-embedding-v2  # 嵌入模型

milvus:
  host: localhost         # Milvus 主机地址
  port: 19530            # Milvus 端口

database:
  base_path: ./data/db   # SQLite 数据库文件路径

rag:
  top_k: 5               # 检索返回的文档数量
  score_threshold: 0.7   # 相似度阈值

security:
  api_keys:              # 向量管理 API Key
    - ${API_KEY_1}
    - ${API_KEY_2}
```

详细配置说明请参考 `config/config.yaml` 文件。

## 📖 API 文档

### 对话接口

发送用户查询，获取智能回复：

```bash
curl -X POST http://localhost:8080/chat \
  -H "Content-Type: application/json" \
  -d '{
    "query": "Python课程包含哪些内容？",
    "tenant_id": "default",
    "session_id": "session-123"
  }'
```

**响应示例**：

```json
{
  "answer": "Python 课程包含以下内容：\n1. 基础语法\n2. 数据结构\n3. 面向对象编程",
  "route": "course",
  "session_id": "session-123",
  "sources": [
    {
      "content": "Python 课程包含基础语法、数据结构...",
      "score": 0.95
    }
  ],
  "metadata": {
    "intent": "course",
    "confidence": 0.92,
    "duration_ms": 234
  }
}
```

### 向量管理

添加文档到知识库：

```bash
curl -X POST http://localhost:8080/api/v1/vectors/items \
  -H "Content-Type: application/json" \
  -H "X-API-Key: your_api_key" \
  -d '{
    "texts": ["Python 课程包含基础语法、数据结构等内容"],
    "tenant_id": "default"
  }'
```

删除文档：

```bash
curl -X DELETE http://localhost:8080/api/v1/vectors/items \
  -H "Content-Type: application/json" \
  -H "X-API-Key: your_api_key" \
  -d '{
    "ids": ["doc-uuid-001"],
    "tenant_id": "default"
  }'
```

### 健康检查

```bash
curl http://localhost:8080/health
```

完整的 API 文档请参考 [docs/API_DOCUMENTATION.md](docs/API_DOCUMENTATION.md)。

## 🔧 开发

### 常用命令

```bash
# 编译项目
make build

# 运行服务
make run

# 运行测试
make test

# 代码格式化
make fmt

# 代码检查
make lint

# 查看测试覆盖率
make test-coverage

# 清理构建产物
make clean
```

### 运行测试

```bash
# 运行所有测试
go test ./...

# 运行单元测试
go test ./internal/...

# 运行集成测试（需要 Milvus）
make test-integration

# 运行性能测试
go test -bench=. ./...

# 生成测试覆盖率报告
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### 代码规范

项目遵循 Go 标准代码规范：
- 使用 `gofmt` 进行代码格式化
- 使用 `golangci-lint` 进行代码检查
- 遵循 [Uber Go Style Guide](https://github.com/uber-go/guide/blob/master/style.md)

## 🏗️ 架构

系统采用 **Clean Architecture（简洁架构）** 模式，遵循依赖倒置原则：

```
┌─────────────────────────────────────────┐
│     Infrastructure Layer                │
│  (Gin, Milvus, SQLite, Eino)           │
└────────────────┬────────────────────────┘
                 │ 依赖
┌────────────────▼────────────────────────┐
│     Interface Adapter Layer             │
│  (HTTP Handler, Middleware)             │
└────────────────┬────────────────────────┘
                 │ 依赖
┌────────────────▼────────────────────────┐
│     Use Case Layer                      │
│  (Chat, Vector Management)              │
└────────────────┬────────────────────────┘
                 │ 依赖
┌────────────────▼────────────────────────┐
│     Domain Layer                        │
│  (Entity, Repository Interface)         │
└─────────────────────────────────────────┘
```

### 核心特点

- **依赖倒置**：内层定义接口，外层实现接口
- **关注点分离**：每层职责清晰，易于测试和维护
- **可扩展性**：易于添加新功能和替换技术实现
- **可测试性**：各层独立测试，支持 Mock

详细架构设计请参考：
- [设计文档](.kiro/specs/eino-qa-system/design.md)
- [项目结构说明](PROJECT_STRUCTURE.md)

## 📚 文档

- [快速开始](QUICKSTART.md) - 5 分钟快速上手
- [启动指南](STARTUP_GUIDE.md) - 详细的启动和配置说明
- [API 文档](docs/API_DOCUMENTATION.md) - 完整的 API 接口文档
- [部署指南](docs/DEPLOYMENT_GUIDE.md) - 生产环境部署指南
- [项目结构](PROJECT_STRUCTURE.md) - 目录结构和分层说明
- [设计文档](.kiro/specs/eino-qa-system/design.md) - 架构设计和技术选型
- [需求文档](.kiro/specs/eino-qa-system/requirements.md) - 功能需求和验收标准

## 🤝 贡献

欢迎贡献代码、报告问题或提出建议！

1. Fork 项目
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 开启 Pull Request

## 📄 许可证

本项目采用 MIT 许可证 - 详见 [LICENSE](LICENSE) 文件

## 🙏 致谢

- [CloudWeGo Eino](https://github.com/cloudwego/eino) - AI 应用开发框架
- [Milvus](https://milvus.io/) - 向量数据库
- [Gin](https://gin-gonic.com/) - Web 框架
- [GORM](https://gorm.io/) - ORM 库

## 📞 联系方式

如有问题或建议，请：
- 提交 [Issue](https://github.com/your-repo/issues)
- 查看 [文档](docs/)
- 发送邮件至：support@example.com

---

**版本**: v1.0.0  
**最后更新**: 2024-11-29
