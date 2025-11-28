# 项目结构说明

## 目录结构

```
eino-qa/
├── cmd/                                    # 应用入口
│   └── server/
│       └── main.go                         # 主程序入口
│
├── internal/                               # 内部包（不对外暴露）
│   ├── domain/                             # 领域层（Domain Layer）
│   │   ├── entity/                         # 领域实体
│   │   │   ├── message.go                  # 消息实体
│   │   │   ├── intent.go                   # 意图实体
│   │   │   ├── document.go                 # 文档实体
│   │   │   └── order.go                    # 订单实体
│   │   └── repository/                     # 仓储接口定义
│   │       ├── vector.go                   # 向量仓储接口
│   │       ├── order.go                    # 订单仓储接口
│   │       └── session.go                  # 会话仓储接口
│   │
│   ├── usecase/                            # 用例层（Use Case Layer）
│   │   ├── chat/                           # 对话用例
│   │   │   └── chat_usecase.go
│   │   ├── vector/                         # 向量管理用例
│   │   │   └── vector_usecase.go
│   │   └── order/                          # 订单查询用例
│   │       └── order_usecase.go
│   │
│   ├── adapter/                            # 接口适配层（Interface Adapter Layer）
│   │   ├── http/
│   │   │   ├── handler/                    # HTTP 处理器
│   │   │   │   ├── chat_handler.go         # 对话接口处理器
│   │   │   │   ├── vector_handler.go       # 向量管理处理器
│   │   │   │   └── health_handler.go       # 健康检查处理器
│   │   │   ├── middleware/                 # HTTP 中间件
│   │   │   │   ├── tenant.go               # 租户识别中间件
│   │   │   │   ├── security.go             # 安全中间件
│   │   │   │   └── logging.go              # 日志中间件
│   │   │   └── router.go                   # 路由配置
│   │   └── presenter/                      # 数据展示层
│   │       └── chat_presenter.go
│   │
│   └── infrastructure/                     # 基础设施层（Infrastructure Layer）
│       ├── repository/                     # 仓储实现
│       │   ├── milvus/                     # Milvus 向量数据库实现
│       │   │   └── vector_repository.go
│       │   ├── sqlite/                     # SQLite 数据库实现
│       │   │   ├── order_repository.go
│       │   │   └── session_repository.go
│       │   └── memory/                     # 内存实现（用于测试）
│       │       └── session_repository.go
│       ├── ai/                             # AI 组件
│       │   └── eino/                       # Eino 框架集成
│       │       ├── intent_recognizer.go    # 意图识别器
│       │       ├── rag_retriever.go        # RAG 检索器
│       │       ├── order_querier.go        # 订单查询器
│       │       └── response_generator.go   # 响应生成器
│       ├── config/                         # 配置管理
│       │   ├── config.go                   # 配置加载
│       │   └── config_test.go              # 配置测试
│       └── logger/                         # 日志管理
│           └── logger.go
│
├── pkg/                                    # 公共包（可对外暴露）
│   ├── errors/                             # 错误处理
│   │   └── errors.go
│   └── utils/                              # 工具函数
│       └── id_generator.go                 # ID 生成器
│
├── config/                                 # 配置文件
│   └── config.yaml                         # 主配置文件
│
├── .env.example                            # 环境变量示例
├── .gitignore                              # Git 忽略文件
├── go.mod                                  # Go 模块定义
├── go.sum                                  # Go 依赖锁定
├── Makefile                                # 构建脚本
└── README.md                               # 项目说明
```

## Clean Architecture 分层说明

### 1. Domain Layer（领域层）- 最内层

**位置**: `internal/domain/`

**职责**:
- 定义核心业务实体（Entity）
- 定义仓储接口（Repository Interface）
- 定义值对象（Value Object）
- 包含业务规则和领域逻辑

**特点**:
- 不依赖任何外部框架或库
- 纯粹的业务逻辑
- 最稳定的层，变化最少

**示例**:
```go
// domain/entity/message.go
type Message struct {
    ID        string
    Content   string
    Role      string
    Timestamp time.Time
}

// domain/repository/vector.go
type VectorRepository interface {
    Search(ctx context.Context, vector []float32, topK int) ([]Document, error)
    Insert(ctx context.Context, docs []Document) error
}
```

### 2. Use Case Layer（用例层）

**位置**: `internal/usecase/`

**职责**:
- 实现应用业务逻辑
- 编排领域对象完成具体用例
- 使用 Eino ADK 进行流程编排

**特点**:
- 依赖领域层的接口
- 不依赖具体的技术实现
- 包含应用特定的业务规则

**示例**:
```go
// usecase/chat/chat_usecase.go
type ChatUseCase struct {
    intentRecognizer IntentRecognizer
    ragRetriever     RAGRetriever
    sessionRepo      domain.SessionRepository
}

func (uc *ChatUseCase) Execute(ctx context.Context, req ChatRequest) (*ChatResponse, error) {
    // 业务逻辑编排
}
```

### 3. Interface Adapter Layer（接口适配层）

**位置**: `internal/adapter/`

**职责**:
- 转换数据格式
- 适配外部接口（HTTP、gRPC 等）
- 实现中间件和过滤器

**特点**:
- 依赖用例层
- 处理数据转换和验证
- 不包含业务逻辑

**示例**:
```go
// adapter/http/handler/chat_handler.go
type ChatHandler struct {
    chatUseCase *usecase.ChatUseCase
}

func (h *ChatHandler) HandleChat(c *gin.Context) {
    // HTTP 请求处理和数据转换
}
```

### 4. Infrastructure Layer（基础设施层）- 最外层

**位置**: `internal/infrastructure/`

**职责**:
- 实现具体的技术细节
- 数据库访问实现
- 外部 API 调用
- 框架集成

**特点**:
- 依赖所有内层
- 包含具体的技术实现
- 最容易变化的层

**示例**:
```go
// infrastructure/repository/milvus/vector_repository.go
type MilvusVectorRepository struct {
    client milvus.Client
}

func (r *MilvusVectorRepository) Search(ctx context.Context, vector []float32, topK int) ([]domain.Document, error) {
    // Milvus 具体实现
}
```

## 依赖规则

```
Infrastructure Layer
        ↓ 依赖
Interface Adapter Layer
        ↓ 依赖
Use Case Layer
        ↓ 依赖
Domain Layer
```

**核心原则**:
- 依赖只能从外层指向内层
- 内层不能依赖外层
- 内层定义接口，外层实现接口（依赖倒置）

## 配置管理

### 配置文件

**config/config.yaml**: 主配置文件，支持环境变量替换

```yaml
dashscope:
  api_key: ${DASHSCOPE_API_KEY}  # 从环境变量读取
```

### 环境变量

**.env**: 本地开发环境变量（不提交到 Git）

**.env.example**: 环境变量模板（提交到 Git）

## 构建和运行

### 使用 Makefile

```bash
# 查看所有命令
make help

# 编译项目
make build

# 运行服务
make run

# 运行测试
make test

# 清理构建产物
make clean
```

### 直接使用 Go 命令

```bash
# 运行服务
go run cmd/server/main.go

# 编译
go build -o bin/server cmd/server/main.go

# 运行测试
go test ./...
```

## 测试策略

### 单元测试

- 位置：与源文件同目录，文件名以 `_test.go` 结尾
- 命名：`Test<FunctionName>`
- 运行：`go test ./...`

### 属性测试

- 使用 `gopter` 库
- 验证通用属性和不变量
- 命名：`TestProperty_<PropertyName>`

### 集成测试

- 位置：`tests/integration/`
- 测试完整的用户流程
- 需要外部依赖（Milvus、数据库）

## 下一步

项目基础设施已搭建完成，接下来的任务：

1. ✅ 项目初始化和基础设施搭建
2. ⏭️ Domain Layer 实现
3. ⏭️ Infrastructure Layer - 配置和日志
4. ⏭️ Infrastructure Layer - Milvus 集成
5. ⏭️ Infrastructure Layer - SQLite 集成
6. ⏭️ Infrastructure Layer - Eino AI 集成
7. ⏭️ Use Case Layer - Chat 用例
8. ⏭️ Use Case Layer - Vector 管理用例
9. ⏭️ Interface Adapter Layer - HTTP 中间件
10. ⏭️ Interface Adapter Layer - HTTP Handler

详细任务列表请参考 `.kiro/specs/eino-qa-system/tasks.md`。
