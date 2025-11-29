# 设计文档

## 概述

本设计文档描述了基于 Eino 框架的 Go 语言智能客服系统的技术架构和实现方案。系统采用 **ADK 为主、Compose 为辅**的混合架构模式，结合 Gin、GORM、SQLite 和 Milvus 等技术栈，实现高性能、可扩展的智能对话服务。

### 核心目标

1. **高性能**: 利用 Go 语言的并发特性和 Eino 框架的优化，实现低延迟、高吞吐的对话处理
2. **可维护性**: 通过 ADK 的高层抽象简化复杂编排逻辑，提升代码可读性和可维护性
3. **可扩展性**: 模块化设计支持灵活添加新的意图类型、工具和数据源
4. **多租户隔离**: 基于 Milvus Collection 和独立 SQLite 数据库实现租户级资源隔离

### 技术栈

- **Web 框架**: Gin - 高性能 HTTP 路由和中间件支持
- **ORM**: GORM - 类型安全的数据库操作
- **关系数据库**: SQLite - 轻量级嵌入式数据库，每租户独立文件
- **向量数据库**: Milvus - 高性能向量检索，支持租户级 Collection 隔离
- **AI 框架**: Eino ADK + Compose - 智能体编排和图处理
- **LLM 服务**: DashScope (通义千问) - 聊天模型和嵌入模型

## 架构

### 整体架构

系统采用 **Clean Architecture（简洁架构）** 模式，遵循依赖倒置原则，从内到外分为：

1. **Domain Layer（领域层）**: 核心业务实体和业务规则，不依赖任何外部框架
2. **Use Case Layer（用例层）**: 应用业务逻辑，编排领域对象完成具体用例
3. **Interface Adapter Layer（接口适配层）**: 转换数据格式，适配外部接口
4. **Infrastructure Layer（基础设施层）**: 外部框架和工具的具体实现

### Clean Architecture 分层详解

#### 1. Domain Layer（领域层 - 最内层）

**职责**: 定义核心业务实体和业务规则，完全独立于框架和外部依赖

```go
// domain/entity/
type Message struct {
    ID        string
    Content   string
    Role      string
    Timestamp time.Time
}

type Intent struct {
    Type       string  // "course", "order", "direct", "handoff"
    Confidence float64
}

type Document struct {
    ID       string
    Content  string
    Metadata map[string]any
    Score    float64
}

type Order struct {
    ID         string
    UserID     string
    CourseName string
    Amount     float64
    Status     string
    CreatedAt  time.Time
}

// domain/repository/ (接口定义，实现在 infrastructure 层)
type VectorRepository interface {
    Search(ctx context.Context, vector []float32, topK int) ([]Document, error)
    Insert(ctx context.Context, docs []Document) error
    Delete(ctx context.Context, ids []string) error
}

type OrderRepository interface {
    FindByID(ctx context.Context, orderID string) (*Order, error)
    FindByUserID(ctx context.Context, userID string) ([]Order, error)
}

type SessionRepository interface {
    Save(ctx context.Context, sessionID string, messages []Message) error
    Load(ctx context.Context, sessionID string) ([]Message, error)
    Delete(ctx context.Context, sessionID string) error
}
```

#### 2. Use Case Layer（用例层）

**职责**: 实现应用业务逻辑，使用 Eino ADK 编排对话流程

```go
// usecase/chat/
type ChatUseCase struct {
    intentRecognizer IntentRecognizer
    ragRetriever     RAGRetriever
    orderQuerier     OrderQuerier
    responseGen      ResponseGenerator
    sessionRepo      domain.SessionRepository
}

func (uc *ChatUseCase) Execute(ctx context.Context, req ChatRequest) (*ChatResponse, error) {
    // 1. 加载会话历史
    history, _ := uc.sessionRepo.Load(ctx, req.SessionID)
    
    // 2. 识别意图
    intent, err := uc.intentRecognizer.Recognize(ctx, req.Query, history)
    if err != nil {
        return nil, err
    }
    
    // 3. 根据意图路由到不同处理流程
    var answer string
    var sources []domain.Document
    
    switch intent.Type {
    case "course":
        answer, sources, err = uc.ragRetriever.Retrieve(ctx, req.Query)
    case "order":
        answer, err = uc.orderQuerier.Query(ctx, req.Query)
    case "direct":
        answer, err = uc.responseGen.Generate(ctx, req.Query, history)
    case "handoff":
        answer = "正在为您转接人工客服..."
    }
    
    // 4. 保存会话
    uc.sessionRepo.Save(ctx, req.SessionID, append(history, 
        domain.Message{Content: req.Query, Role: "user"},
        domain.Message{Content: answer, Role: "assistant"},
    ))
    
    return &ChatResponse{
        Answer:  answer,
        Route:   intent.Type,
        Sources: sources,
    }, nil
}

// usecase/vector/
type VectorManagementUseCase struct {
    vectorRepo domain.VectorRepository
    embedder   Embedder
}

func (uc *VectorManagementUseCase) AddVectors(ctx context.Context, texts []string) error {
    // 生成向量
    vectors, err := uc.embedder.Embed(ctx, texts)
    if err != nil {
        return err
    }
    
    // 构建文档
    docs := make([]domain.Document, len(texts))
    for i, text := range texts {
        docs[i] = domain.Document{
            ID:      generateID(),
            Content: text,
            Vector:  vectors[i],
        }
    }
    
    // 插入向量库
    return uc.vectorRepo.Insert(ctx, docs)
}
```

#### 3. Interface Adapter Layer（接口适配层）

**职责**: 转换数据格式，适配 HTTP、gRPC 等外部接口

```go
// adapter/http/handler/
type ChatHandler struct {
    chatUseCase *usecase.ChatUseCase
}

func (h *ChatHandler) HandleChat(c *gin.Context) {
    var req ChatRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, ErrorResponse{Message: "invalid request"})
        return
    }
    
    // 提取租户信息
    tenantID := c.GetHeader("X-Tenant-ID")
    if tenantID == "" {
        tenantID = "default"
    }
    
    // 调用用例
    ctx := context.WithValue(c.Request.Context(), "tenant_id", tenantID)
    resp, err := h.chatUseCase.Execute(ctx, req)
    if err != nil {
        c.JSON(500, ErrorResponse{Message: err.Error()})
        return
    }
    
    c.JSON(200, resp)
}

// adapter/http/middleware/
func TenantMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        tenantID := c.GetHeader("X-Tenant-ID")
        if tenantID == "" {
            tenantID = c.Query("tenant")
        }
        if tenantID == "" {
            tenantID = "default"
        }
        c.Set("tenant_id", tenantID)
        c.Next()
    }
}

func SecurityMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // 脱敏处理
        // API Key 验证
        c.Next()
    }
}
```

#### 4. Infrastructure Layer（基础设施层 - 最外层）

**职责**: 实现具体的技术细节，如数据库访问、外部 API 调用

```go
// infrastructure/repository/milvus/
type MilvusVectorRepository struct {
    client     milvus.Client
    collection string
}

func (r *MilvusVectorRepository) Search(ctx context.Context, vector []float32, topK int) ([]domain.Document, error) {
    // 调用 Milvus SDK
    results, err := r.client.Search(ctx, r.collection, vector, topK)
    if err != nil {
        return nil, err
    }
    
    // 转换为领域对象
    docs := make([]domain.Document, len(results))
    for i, result := range results {
        docs[i] = domain.Document{
            ID:      result.ID,
            Content: result.Fields["content"].(string),
            Score:   result.Score,
        }
    }
    return docs, nil
}

// infrastructure/repository/sqlite/
type SQLiteOrderRepository struct {
    db *gorm.DB
}

func (r *SQLiteOrderRepository) FindByID(ctx context.Context, orderID string) (*domain.Order, error) {
    var order OrderModel
    err := r.db.WithContext(ctx).Where("id = ?", orderID).First(&order).Error
    if err != nil {
        return nil, err
    }
    
    // 转换为领域对象
    return &domain.Order{
        ID:         order.ID,
        UserID:     order.UserID,
        CourseName: order.CourseName,
        Amount:     order.Amount,
        Status:     order.Status,
        CreatedAt:  order.CreatedAt,
    }, nil
}

// infrastructure/ai/eino/
type EinoIntentRecognizer struct {
    chatModel model.ChatModel
    prompt    *prompt.ChatTemplate
}

func (r *EinoIntentRecognizer) Recognize(ctx context.Context, query string, history []domain.Message) (*domain.Intent, error) {
    // 使用 Eino ChatModel
    messages := buildMessages(query, history)
    resp, err := r.chatModel.Generate(ctx, messages)
    if err != nil {
        return nil, err
    }
    
    // 解析意图
    intent := parseIntent(resp.Content)
    return &domain.Intent{
        Type:       intent.Type,
        Confidence: intent.Confidence,
    }, nil
}
```

### 依赖关系图

```
┌─────────────────────────────────────────────────────────────┐
│                Infrastructure Layer                          │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐      │
│  │ Gin Server   │  │ Milvus Repo  │  │ SQLite Repo  │      │
│  │ Eino AI      │  │ DashScope    │  │ GORM         │      │
│  └──────────────┘  └──────────────┘  └──────────────┘      │
└────────────────────┬────────────────────────────────────────┘
                     │ 依赖
┌────────────────────▼────────────────────────────────────────┐
│              Interface Adapter Layer                         │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐      │
│  │ HTTP Handler │  │ Middleware   │  │ Presenter    │      │
│  └──────────────┘  └──────────────┘  └──────────────┘      │
└────────────────────┬────────────────────────────────────────┘
                     │ 依赖
┌────────────────────▼────────────────────────────────────────┐
│                  Use Case Layer                              │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐      │
│  │ Chat UseCase │  │ Vector Mgmt  │  │ Order Query  │      │
│  └──────────────┘  └──────────────┘  └──────────────┘      │
└────────────────────┬────────────────────────────────────────┘
                     │ 依赖
┌────────────────────▼────────────────────────────────────────┐
│                    Domain Layer                              │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐      │
│  │ Entity       │  │ Repository   │  │ Value Object │      │
│  │              │  │ Interface    │  │              │      │
│  └──────────────┘  └──────────────┘  └──────────────┘      │
└─────────────────────────────────────────────────────────────┘
```

### 目录结构

```
eino-qa/
├── cmd/
│   └── server/
│       └── main.go                 # 应用入口
├── internal/
│   ├── domain/                     # 领域层
│   │   ├── entity/
│   │   │   ├── message.go
│   │   │   ├── intent.go
│   │   │   ├── document.go
│   │   │   └── order.go
│   │   └── repository/             # 仓储接口
│   │       ├── vector.go
│   │       ├── order.go
│   │       └── session.go
│   ├── usecase/                    # 用例层
│   │   ├── chat/
│   │   │   └── chat_usecase.go
│   │   ├── vector/
│   │   │   └── vector_usecase.go
│   │   └── order/
│   │       └── order_usecase.go
│   ├── adapter/                    # 接口适配层
│   │   ├── http/
│   │   │   ├── handler/
│   │   │   │   ├── chat_handler.go
│   │   │   │   ├── vector_handler.go
│   │   │   │   └── health_handler.go
│   │   │   ├── middleware/
│   │   │   │   ├── tenant.go
│   │   │   │   ├── security.go
│   │   │   │   └── logging.go
│   │   │   └── router.go
│   │   └── presenter/
│   │       └── chat_presenter.go
│   └── infrastructure/             # 基础设施层
│       ├── repository/
│       │   ├── milvus/
│       │   │   └── vector_repository.go
│       │   ├── sqlite/
│       │   │   ├── order_repository.go
│       │   │   └── session_repository.go
│       │   └── memory/
│       │       └── session_repository.go
│       ├── ai/
│       │   └── eino/
│       │       ├── intent_recognizer.go
│       │       ├── rag_retriever.go
│       │       ├── order_querier.go
│       │       └── response_generator.go
│       ├── config/
│       │   └── config.go
│       └── logger/
│           └── logger.go
├── pkg/                            # 公共工具包
│   ├── errors/
│   │   └── errors.go
│   └── utils/
│       └── id_generator.go
├── config/
│   └── config.yaml
├── go.mod
└── go.sum
```

### 架构图

```
┌─────────────────────────────────────────────────────────────┐
│                        API 层 (Gin)                          │
│  /chat  /health  /vectors/add  /vectors/delete  /models     │
└────────────────────┬────────────────────────────────────────┘
                     │
┌────────────────────▼────────────────────────────────────────┐
│                   编排层 (Eino ADK)                          │
│  ┌──────────────────────────────────────────────────────┐   │
│  │  Sequential Agent (主对话流程)                        │   │
│  │    ├─ Intent Recognition Agent                       │   │
│  │    ├─ Parallel Agent (信息收集)                      │   │
│  │    │    ├─ RAG Retrieval Agent                       │   │
│  │    │    └─ Order Query Agent                         │   │
│  │    └─ Response Generation Agent                      │   │
│  └──────────────────────────────────────────────────────┘   │
└────────────────────┬────────────────────────────────────────┘
                     │
┌────────────────────▼────────────────────────────────────────┐
│                 执行层 (Eino Compose)                        │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐         │
│  │ Intent Node │  │  RAG Node   │  │ Order Node  │         │
│  └─────────────┘  └─────────────┘  └─────────────┘         │
└────────────────────┬────────────────────────────────────────┘
                     │
┌────────────────────▼────────────────────────────────────────┐
│                      工具层                                  │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐      │
│  │ Embedding    │  │ Vector Search│  │ SQL Generator│      │
│  │ Generator    │  │              │  │              │      │
│  └──────────────┘  └──────────────┘  └──────────────┘      │
└────────────────────┬────────────────────────────────────────┘
                     │
┌────────────────────▼────────────────────────────────────────┐
│                     存储层                                   │
│  ┌──────────────────────┐  ┌──────────────────────┐         │
│  │  Milvus              │  │  SQLite              │         │
│  │  (向量数据库)         │  │  (关系数据库)         │         │
│  │  - Collection per    │  │  - DB file per       │         │
│  │    tenant            │  │    tenant            │         │
│  └──────────────────────┘  └──────────────────────┘         │
└─────────────────────────────────────────────────────────────┘
```

### 混合架构模式说明

**ADK 使用场景**:
- 主对话流程编排（Sequential Agent）
- 多源信息并行收集（Parallel Agent）
- 意图识别和路由决策

**Compose 使用场景**:
- 单个智能体内部的细粒度控制
- 性能关键路径（如 RAG 检索流程）
- 需要精确控制数据流转的场景

## 组件和接口

### 核心组件

#### 1. HTTP Server (Gin)

```go
type Server struct {
    router  *gin.Engine
    agent   adk.Agent
    config  *Config
}

// 主要路由
// POST /chat - 对话接口
// GET /health - 健康检查
// POST /api/v1/vectors/items - 添加向量
// DELETE /api/v1/vectors/items - 删除向量
// POST /models/switch - 切换模型
```

#### 2. Main Agent (ADK Sequential)

```go
type MainAgent struct {
    intentAgent    adk.Agent  // 意图识别智能体
    infoAgent      adk.Agent  // 信息收集智能体（并行）
    responseAgent  adk.Agent  // 响应生成智能体
}
```

#### 3. Intent Recognition Agent

```go
type IntentAgent struct {
    chatModel  model.ChatModel
    prompt     *prompt.ChatTemplate
}

// 输出结构
type IntentResult struct {
    Intent     string  // "course", "order", "direct", "handoff"
    Confidence float64
}
```

#### 4. RAG Agent (Compose)

```go
type RAGAgent struct {
    embedder     embedding.Embedder
    vectorStore  indexer.VectorIndexer  // Milvus
    chatModel    model.ChatModel
}

// 检索流程
// Query -> Embed -> Search -> Rerank -> Generate
```

#### 5. Order Query Agent (Compose)

```go
type OrderAgent struct {
    db         *gorm.DB
    chatModel  model.ChatModel
    sqlGen     *SQLGenerator
}

// 查询流程
// Query -> Extract OrderID -> Generate SQL -> Validate -> Execute -> Format
```

#### 6. Multi-tenant Manager

```go
type TenantManager struct {
    collections map[string]string  // tenantID -> collectionName
    databases   map[string]*gorm.DB
}

func (tm *TenantManager) GetCollection(tenantID string) (string, error)
func (tm *TenantManager) GetDB(tenantID string) (*gorm.DB, error)
```

### 接口定义

#### Chat Request/Response

```go
type ChatRequest struct {
    Query     string   `json:"query" binding:"required"`
    TenantID  string   `json:"tenant_id"`
    SessionID string   `json:"session_id"`
    Stream    bool     `json:"stream"`
}

type ChatResponse struct {
    Answer    string            `json:"answer"`
    Route     string            `json:"route"`
    Sources   []Source          `json:"sources,omitempty"`
    Metadata  map[string]any    `json:"metadata"`
}

type Source struct {
    Content  string  `json:"content"`
    Score    float64 `json:"score"`
    Metadata map[string]any `json:"metadata"`
}
```

#### Vector Management

```go
type VectorAddRequest struct {
    Texts     []string `json:"texts" binding:"required"`
    TenantID  string   `json:"tenant_id"`
}

type VectorDeleteRequest struct {
    IDs       []string `json:"ids" binding:"required"`
    TenantID  string   `json:"tenant_id"`
}
```

## 数据模型

### Milvus Schema

```go
// Collection Schema (per tenant)
type KnowledgeDocument struct {
    ID        string    `milvus:"primary_key"`
    Vector    []float32 `milvus:"dim:1536"`
    Content   string
    Metadata  string    // JSON encoded
    TenantID  string
    CreatedAt int64
}
```

### SQLite Models (GORM)

```go
// 订单表
type Order struct {
    ID          string    `gorm:"primaryKey"`
    UserID      string    `gorm:"index"`
    CourseName  string
    Amount      float64
    Status      string    // "pending", "paid", "refunded", "cancelled"
    CreatedAt   time.Time
    UpdatedAt   time.Time
}

// 会话表
type Session struct {
    ID        string    `gorm:"primaryKey"`
    TenantID  string    `gorm:"index"`
    Messages  string    // JSON encoded message history
    CreatedAt time.Time
    UpdatedAt time.Time
}

// 未命中记录表
type MissedQuery struct {
    ID        uint      `gorm:"primaryKey"`
    TenantID  string    `gorm:"index"`
    Query     string
    Intent    string
    CreatedAt time.Time
}
```

### 配置模型

```go
type Config struct {
    Server   ServerConfig
    DashScope DashScopeConfig
    Milvus   MilvusConfig
    Database DatabaseConfig
}

type ServerConfig struct {
    Port     int
    Mode     string  // "debug", "release"
}

type DashScopeConfig struct {
    APIKey      string
    ChatModel   string  // "qwen-turbo", "qwen-plus"
    EmbedModel  string  // "text-embedding-v2"
}

type MilvusConfig struct {
    Host       string
    Port       int
    Username   string
    Password   string
}

type DatabaseConfig struct {
    BasePath   string  // SQLite 数据库文件基础路径
}
```


## 正确性属性

*属性是一个特征或行为，应该在系统的所有有效执行中保持为真——本质上是关于系统应该做什么的形式化陈述。属性作为人类可读规范和机器可验证正确性保证之间的桥梁。*

### Property 1: 意图分类完整性
*对于任意*用户查询，意图识别系统返回的意图类型必须是 {"course", "order", "direct", "handoff"} 集合中的一个
**验证需求: 2.1**

### Property 2: 意图识别结果结构完整性
*对于任意*意图识别操作，返回结果必须同时包含意图类型（Intent）和置信度分数（Confidence）字段
**验证需求: 2.2**

### Property 3: 对话流程使用正确的编排引擎
*对于任意*用户对话请求，系统必须使用 Eino Compose 或 ADK 进行流程编排处理
**验证需求: 1.2**

### Property 4: 流式响应支持
*对于任意*对话流程执行，当客户端请求流式响应时，系统必须以分块方式返回内容而非一次性返回完整响应
**验证需求: 1.3**

### Property 5: RAG 向量生成
*对于任意*课程咨询查询，RAG 系统必须调用嵌入模型生成查询向量
**验证需求: 3.1**

### Property 6: RAG 搜索执行
*对于任意*生成的查询向量，RAG 系统必须在 Milvus 向量数据库中执行相似度搜索
**验证需求: 3.2**

### Property 7: RAG 结果数量限制
*对于任意*相似度搜索操作，返回的文档片段数量必须不超过配置的 K 值
**验证需求: 3.3**

### Property 8: RAG 文档传递给 LLM
*对于任意*检索到的相关文档，RAG 系统必须将文档内容和原始查询一起传递给 LLM 进行答案生成
**验证需求: 3.4**

### Property 9: 订单 ID 提取
*对于任意*包含订单 ID 格式（如 #20251114001）的用户输入，订单查询系统必须能够正确提取订单 ID
**验证需求: 4.1**

### Property 10: SQL 安全性验证
*对于任意*生成的 SQL 查询语句，系统必须验证其不包含危险操作（DROP、DELETE、UPDATE 等）
**验证需求: 4.3, 8.5**

### Property 11: SQL 注入防护
*对于任意*订单 ID 输入（包括恶意构造的输入），生成的 SQL 查询必须使用参数化查询或安全转义，不存在 SQL 注入漏洞
**验证需求: 4.2, 4.3**

### Property 12: 订单查询结果格式化
*对于任意*数据库查询返回的结构化数据，订单查询系统必须将其格式化为自然语言响应而非返回原始 JSON 或表格数据
**验证需求: 4.5**

### Property 13: 租户 ID 提取
*对于任意*包含租户标识的 HTTP 请求（请求头 X-Tenant-ID 或查询参数 tenant），多租户系统必须正确提取租户 ID
**验证需求: 5.1**

### Property 14: 租户 Collection 映射
*对于任意*确定的租户 ID，多租户系统必须使用该租户对应的 Milvus Collection 进行向量操作
**验证需求: 5.3**

### Property 15: 租户数据库隔离
*对于任意*租户的数据库操作，多租户系统必须使用该租户独立的 SQLite 数据库文件，不同租户的数据不能交叉访问
**验证需求: 5.5**

### Property 16: HTTP 请求解析
*对于任意*符合 JSON 格式的 POST /chat 请求体，Gin HTTP 服务器必须成功解析并提取查询参数
**验证需求: 6.2**

### Property 17: HTTP 响应格式
*对于任意*成功处理的对话请求，HTTP 服务器返回的 JSON 响应必须包含 answer 和 metadata 字段
**验证需求: 6.3**

### Property 18: HTTP 错误响应
*对于任意*处理失败的请求，HTTP 服务器必须返回适当的 HTTP 状态码（4xx 或 5xx）和包含错误信息的响应体
**验证需求: 6.5**

### Property 19: 日志记录完整性
*对于任意*系统处理的请求，日志系统必须记录请求 ID、租户 ID、查询内容和处理时长
**验证需求: 7.1**

### Property 20: 外部调用日志
*对于任意*外部服务调用（LLM、Milvus、数据库），日志系统必须记录调用参数和响应时间
**验证需求: 7.3**

### Property 21: 指标统计
*对于任意*类型的请求（course、order、direct、handoff），指标系统必须统计该类型的请求数量和平均响应时间
**验证需求: 7.4**

### Property 22: 敏感信息脱敏
*对于任意*包含敏感字段（密码、身份证号、手机号）的日志内容，安全系统必须将敏感信息替换为占位符（如 [REDACTED]）
**验证需求: 8.1**

### Property 23: 输入脱敏
*对于任意*包含敏感信息的用户输入，安全系统必须在处理前将敏感信息替换为占位符
**验证需求: 8.2**

### Property 24: API Key 验证
*对于任意*访问管理接口（/api/v1/vectors/*）的请求，安全系统必须验证请求头中的 API Key 是否有效
**验证需求: 8.3, 9.1**

### Property 25: 向量生成
*对于任意*文本输入，向量管理系统必须使用 Eino 嵌入模型生成对应的向量表示
**验证需求: 9.2**

### Property 26: 向量插入
*对于任意*生成的向量，向量管理系统必须成功将其插入到对应租户的 Milvus Collection 中
**验证需求: 9.3**

### Property 27: 向量删除
*对于任意*有效的文档 ID，向量管理系统必须能够从 Milvus Collection 中删除对应的向量
**验证需求: 9.4**

### Property 28: 向量操作结果
*对于任意*向量操作（添加或删除），系统必须返回包含操作结果状态和受影响记录数的响应
**验证需求: 9.5**

### Property 29: 会话 ID 唯一性
*对于任意*两次独立的新对话请求，会话管理系统生成的会话 ID 必须不同
**验证需求: 10.1**

### Property 30: 消息历史追加
*对于任意*用户发送的消息或系统生成的响应，会话管理系统必须将其追加到对应会话的历史记录中
**验证需求: 10.2, 10.3**

### Property 31: 上下文加载
*对于任意*后续消息处理，会话管理系统必须加载该会话的完整历史记录作为上下文
**验证需求: 10.4**

## 错误处理

### 错误分类

1. **客户端错误 (4xx)**
   - 400 Bad Request: 请求参数格式错误、缺少必需字段
   - 401 Unauthorized: API Key 验证失败
   - 404 Not Found: 资源不存在（如订单 ID 不存在）
   - 429 Too Many Requests: 请求频率超限

2. **服务端错误 (5xx)**
   - 500 Internal Server Error: 未预期的系统错误
   - 502 Bad Gateway: 外部服务（DashScope、Milvus）调用失败
   - 503 Service Unavailable: 系统过载或维护中

### 错误处理策略

```go
type ErrorResponse struct {
    Code    int    `json:"code"`
    Message string `json:"message"`
    Details any    `json:"details,omitempty"`
    TraceID string `json:"trace_id"`
}

// 错误处理中间件
func ErrorHandler() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Next()
        
        if len(c.Errors) > 0 {
            err := c.Errors.Last()
            
            // 根据错误类型返回适当的状态码
            statusCode := determineStatusCode(err.Err)
            
            c.JSON(statusCode, ErrorResponse{
                Code:    statusCode,
                Message: err.Error(),
                TraceID: c.GetString("trace_id"),
            })
        }
    }
}
```

### 重试策略

```go
type RetryConfig struct {
    MaxAttempts int
    InitialDelay time.Duration
    MaxDelay     time.Duration
    Multiplier   float64
}

// 指数退避重试
func RetryWithBackoff(ctx context.Context, config RetryConfig, fn func() error) error {
    delay := config.InitialDelay
    
    for attempt := 0; attempt < config.MaxAttempts; attempt++ {
        err := fn()
        if err == nil {
            return nil
        }
        
        // 判断是否可重试
        if !isRetryable(err) {
            return err
        }
        
        // 等待后重试
        select {
        case <-time.After(delay):
            delay = time.Duration(float64(delay) * config.Multiplier)
            if delay > config.MaxDelay {
                delay = config.MaxDelay
            }
        case <-ctx.Done():
            return ctx.Err()
        }
    }
    
    return fmt.Errorf("max retry attempts exceeded")
}
```

### 降级策略

1. **LLM 调用失败**: 返回预设的兜底回复，引导用户联系人工客服
2. **Milvus 不可用**: 降级到基于关键词的简单匹配
3. **数据库查询失败**: 返回错误提示，记录到未命中表
4. **嵌入模型失败**: 使用缓存的向量或跳过向量检索

## 测试策略

### 单元测试

使用 Go 标准库 `testing` 和 `testify` 进行单元测试：

```go
// 测试意图识别
func TestIntentRecognition(t *testing.T) {
    agent := NewIntentAgent(mockChatModel)
    
    tests := []struct {
        name     string
        query    string
        expected string
    }{
        {"课程咨询", "Python课程包含哪些内容？", "course"},
        {"订单查询", "查询订单#20251114001", "order"},
        {"直接回答", "你好", "direct"},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result, err := agent.Recognize(context.Background(), tt.query)
            assert.NoError(t, err)
            assert.Equal(t, tt.expected, result.Intent)
        })
    }
}
```

### 属性测试

使用 `gopter` 库进行基于属性的测试：

```go
import "github.com/leanovate/gopter"
import "github.com/leanovate/gopter/gen"
import "github.com/leanovate/gopter/prop"

// Property 1: 意图分类完整性
func TestProperty_IntentClassificationCompleteness(t *testing.T) {
    properties := gopter.NewProperties(nil)
    
    properties.Property("意图必须是有效值之一", prop.ForAll(
        func(query string) bool {
            agent := NewIntentAgent(mockChatModel)
            result, err := agent.Recognize(context.Background(), query)
            if err != nil {
                return false
            }
            
            validIntents := []string{"course", "order", "direct", "handoff"}
            return contains(validIntents, result.Intent)
        },
        gen.AnyString(),
    ))
    
    properties.TestingRun(t, gopter.ConsoleReporter(false))
}

// Property 11: SQL 注入防护
func TestProperty_SQLInjectionProtection(t *testing.T) {
    properties := gopter.NewProperties(nil)
    
    properties.Property("SQL 查询必须安全", prop.ForAll(
        func(orderID string) bool {
            agent := NewOrderAgent(mockDB, mockChatModel)
            sql, err := agent.GenerateSQL(orderID)
            if err != nil {
                return true // 拒绝无效输入也是正确的
            }
            
            // 验证 SQL 不包含危险模式
            dangerous := []string{"DROP", "DELETE", "UPDATE", "--", ";"}
            for _, pattern := range dangerous {
                if strings.Contains(strings.ToUpper(sql), pattern) {
                    return false
                }
            }
            
            // 验证使用了参数化查询
            return strings.Contains(sql, "?") || strings.Contains(sql, "$")
        },
        gen.AnyString(),
    ))
    
    properties.TestingRun(t, gopter.ConsoleReporter(false))
}

// Property 29: 会话 ID 唯一性
func TestProperty_SessionIDUniqueness(t *testing.T) {
    properties := gopter.NewProperties(nil)
    
    properties.Property("会话 ID 必须唯一", prop.ForAll(
        func(n int) bool {
            manager := NewSessionManager()
            ids := make(map[string]bool)
            
            // 生成 n 个会话 ID
            for i := 0; i < n; i++ {
                id := manager.CreateSession()
                if ids[id] {
                    return false // 发现重复
                }
                ids[id] = true
            }
            
            return true
        },
        gen.IntRange(10, 1000),
    ))
    
    properties.TestingRun(t, gopter.ConsoleReporter(false))
}
```

### 集成测试

测试完整的对话流程：

```go
func TestIntegration_ChatFlow(t *testing.T) {
    // 启动测试服务器
    server := setupTestServer(t)
    defer server.Close()
    
    // 测试课程咨询流程
    t.Run("课程咨询", func(t *testing.T) {
        resp := sendChatRequest(t, server.URL, ChatRequest{
            Query:    "Python课程包含哪些内容？",
            TenantID: "test",
        })
        
        assert.Equal(t, "course", resp.Route)
        assert.NotEmpty(t, resp.Answer)
        assert.NotEmpty(t, resp.Sources)
    })
    
    // 测试订单查询流程
    t.Run("订单查询", func(t *testing.T) {
        resp := sendChatRequest(t, server.URL, ChatRequest{
            Query:    "查询订单#20251114001",
            TenantID: "test",
        })
        
        assert.Equal(t, "order", resp.Route)
        assert.Contains(t, resp.Answer, "订单")
    })
}
```

### 性能测试

使用 `testing` 包的 Benchmark 功能：

```go
func BenchmarkChatEndpoint(b *testing.B) {
    server := setupTestServer(b)
    defer server.Close()
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        sendChatRequest(b, server.URL, ChatRequest{
            Query:    "Python课程包含哪些内容？",
            TenantID: "test",
        })
    }
}

func BenchmarkRAGRetrieval(b *testing.B) {
    agent := setupRAGAgent(b)
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        agent.Retrieve(context.Background(), "Python课程")
    }
}
```

### 测试覆盖率目标

- 单元测试覆盖率: ≥ 80%
- 属性测试: 覆盖所有关键正确性属性
- 集成测试: 覆盖所有主要用户流程
- 性能测试: P95 响应时间 < 500ms

