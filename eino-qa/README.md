# Eino QA System

基于 CloudWeGo Eino 框架构建的 Go 语言智能客服系统。

## 特性

- 🚀 高性能：基于 Go 语言和 Eino 框架
- 🎯 智能路由：自动识别用户意图（课程咨询、订单查询、直接回答、人工转接）
- 📚 RAG 检索：基于 Milvus 向量数据库的知识库检索
- 🔐 多租户：支持租户级数据隔离
- 🛡️ 安全：敏感信息脱敏、API Key 验证、SQL 注入防护
- 📊 可观测：结构化日志、指标统计、健康检查

## 技术栈

- **Web 框架**: Gin
- **ORM**: GORM
- **关系数据库**: SQLite
- **向量数据库**: Milvus
- **AI 框架**: Eino ADK + Compose
- **LLM 服务**: DashScope (通义千问)

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

## 快速开始

### 前置要求

- Go 1.23+
- Milvus 2.4+
- DashScope API Key

### 安装

1. 克隆项目
```bash
git clone <repository-url>
cd eino-qa
```

2. 安装依赖
```bash
go mod download
```

3. 配置环境变量
```bash
cp .env.example .env
# 编辑 .env 文件，填入你的 API Key
```

4. 启动服务
```bash
go run cmd/server/main.go
```

服务将在 `http://localhost:8080` 启动。

## 配置

配置文件位于 `config/config.yaml`，支持通过环境变量覆盖配置项。

主要配置项：
- `server`: HTTP 服务器配置
- `dashscope`: DashScope API 配置
- `milvus`: Milvus 向量数据库配置
- `database`: SQLite 数据库配置
- `rag`: RAG 检索配置
- `security`: 安全配置

详细配置说明请参考 `config/config.yaml` 文件。

## API 文档

### 对话接口

```bash
POST /chat
Content-Type: application/json

{
  "query": "Python课程包含哪些内容？",
  "tenant_id": "default",
  "session_id": "session-123",
  "stream": false
}
```

### 向量管理

```bash
# 添加向量
POST /api/v1/vectors/items
X-API-Key: your-api-key
Content-Type: application/json

{
  "texts": ["文档内容1", "文档内容2"],
  "tenant_id": "default"
}

# 删除向量
DELETE /api/v1/vectors/items
X-API-Key: your-api-key
Content-Type: application/json

{
  "ids": ["doc-id-1", "doc-id-2"],
  "tenant_id": "default"
}
```

### 健康检查

```bash
GET /health
```

## 开发

### 运行测试

```bash
# 运行所有测试
go test ./...

# 运行单元测试
go test ./internal/...

# 运行性能测试
go test -bench=. ./...
```

### 代码规范

项目遵循 Go 标准代码规范，使用 `gofmt` 和 `golint` 进行代码格式化和检查。

## 架构

系统采用 Clean Architecture（简洁架构）模式，分为四层：

1. **Domain Layer**: 核心业务实体和业务规则
2. **Use Case Layer**: 应用业务逻辑
3. **Interface Adapter Layer**: 接口适配（HTTP Handler、中间件）
4. **Infrastructure Layer**: 外部框架和工具实现

详细架构设计请参考 `.kiro/specs/eino-qa-system/design.md`。

## 许可证

[MIT License](LICENSE)
