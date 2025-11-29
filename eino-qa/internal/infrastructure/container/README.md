# 依赖注入容器

## 概述

依赖注入容器负责初始化和管理系统的所有组件，确保组件按正确的依赖顺序初始化，并提供统一的资源管理和优雅关闭机制。

## 组件初始化顺序

容器按以下顺序初始化组件，确保依赖关系正确：

1. **日志系统** (`initLogger`)
   - 初始化 logrus logger
   - 初始化结构化日志接口
   - 配置日志级别、格式和输出

2. **指标收集器** (`initMetrics`)
   - 初始化 metrics.Collector
   - 用于收集系统运行指标

3. **Eino 客户端** (`initEinoClient`)
   - 初始化 DashScope 聊天模型
   - 初始化 DashScope 嵌入模型
   - 配置重试策略和超时

4. **Milvus 客户端** (`initMilvusClient`)
   - 连接到 Milvus 向量数据库
   - 验证连接可用性

5. **多租户管理** (`initTenantManagement`)
   - 初始化 Collection 管理器
   - 初始化 Milvus 租户管理器
   - 初始化 SQLite 数据库管理器
   - 初始化统一租户管理器

6. **仓储层** (`initRepositories`)
   - 初始化向量仓储（Milvus）
   - 初始化订单仓储（SQLite）
   - 初始化会话仓储（SQLite）

7. **AI 组件** (`initAIComponents`)
   - 初始化意图识别器
   - 初始化 RAG 检索器
   - 初始化订单查询器
   - 初始化响应生成器

8. **用例层** (`initUseCases`)
   - 初始化对话用例
   - 初始化向量管理用例

9. **中间件** (`initMiddlewares`)
   - 初始化租户识别中间件
   - 初始化安全脱敏中间件
   - 初始化日志记录中间件
   - 初始化指标收集中间件
   - 初始化错误处理中间件
   - 初始化 API Key 认证中间件

10. **处理器** (`initHandlers`)
    - 初始化对话处理器
    - 初始化向量管理处理器
    - 初始化模型管理处理器
    - 初始化健康检查处理器

11. **HTTP 服务器** (`initHTTPServer`)
    - 配置 Gin 路由
    - 创建 HTTP 服务器
    - 配置服务器参数

## 使用方法

### 创建容器

```go
import (
    "eino-qa/internal/infrastructure/config"
    "eino-qa/internal/infrastructure/container"
)

// 加载配置
cfg, err := config.Load("config/config.yaml")
if err != nil {
    log.Fatal(err)
}

// 创建容器
c, err := container.New(cfg)
if err != nil {
    log.Fatal(err)
}
defer c.Close()
```

### 启动服务器

```go
// 启动 HTTP 服务器
if err := c.Server.Run(); err != nil {
    log.Fatal(err)
}
```

### 访问组件

容器初始化后，所有组件都可以通过容器实例访问：

```go
// 访问日志
c.Logger.Info(ctx, "message", map[string]interface{}{"key": "value"})

// 访问用例
resp, err := c.ChatUseCase.Execute(ctx, req)

// 访问租户管理器
tenant, err := c.TenantManager.GetTenant(ctx, "tenant1")
```

## 优雅关闭

容器的 `Close()` 方法会按正确的顺序关闭所有组件，释放资源：

1. 关闭租户管理器（关闭所有数据库连接）
2. 关闭 Milvus 客户端
3. 关闭 Eino 客户端

```go
defer func() {
    if err := c.Close(); err != nil {
        log.Printf("Error closing container: %v", err)
    }
}()
```

## 错误处理

如果任何组件初始化失败，容器创建会立即返回错误，并提供详细的错误信息：

```go
c, err := container.New(cfg)
if err != nil {
    // 错误信息会指明哪个组件初始化失败
    log.Fatalf("Failed to initialize container: %v", err)
}
```

## 配置要求

容器需要完整的配置对象，包括：

- `Server`: HTTP 服务器配置
- `DashScope`: DashScope API 配置
- `Milvus`: Milvus 数据库配置
- `Database`: SQLite 数据库配置
- `RAG`: RAG 检索配置
- `Intent`: 意图识别配置
- `Session`: 会话管理配置
- `Security`: 安全配置
- `Logging`: 日志配置

详见 `config/config.yaml` 示例配置文件。

## 依赖关系图

```
Container
├── Logger (基础)
├── Metrics (基础)
├── EinoClient (基础)
│   ├── ChatModel
│   └── EmbedModel
├── MilvusClient (基础)
├── TenantManagement
│   ├── MilvusTenantManager
│   │   └── CollectionManager
│   ├── DBManager
│   └── TenantManager
├── Repositories
│   ├── VectorRepository (依赖: MilvusClient, EmbedModel)
│   ├── OrderRepository (依赖: DBManager)
│   └── SessionRepository (依赖: DBManager)
├── AI Components
│   ├── IntentRecognizer (依赖: ChatModel)
│   ├── RAGRetriever (依赖: ChatModel, EmbedModel, VectorRepository)
│   ├── OrderQuerier (依赖: ChatModel, OrderRepository)
│   └── ResponseGenerator (依赖: ChatModel)
├── UseCases
│   ├── ChatUseCase (依赖: AI Components, SessionRepository)
│   └── VectorUseCase (依赖: VectorRepository, EmbedModel)
├── Middlewares
├── Handlers (依赖: UseCases)
└── Server (依赖: Handlers, Middlewares)
```

## 测试

容器支持测试模式，可以注入 mock 组件：

```go
// 创建测试配置
cfg := &config.Config{
    Server: config.ServerConfig{
        Port: 8080,
        Mode: "test",
    },
    // ... 其他配置
}

// 创建容器
c, err := container.New(cfg)
if err != nil {
    t.Fatal(err)
}
defer c.Close()

// 使用容器进行测试
// ...
```

## 注意事项

1. **初始化顺序**: 不要修改组件初始化顺序，否则可能导致依赖错误
2. **资源释放**: 始终使用 `defer c.Close()` 确保资源正确释放
3. **配置验证**: 容器创建前会验证配置的有效性
4. **错误处理**: 任何组件初始化失败都会导致容器创建失败
5. **并发安全**: 容器本身不是并发安全的，应该在应用启动时创建一次
