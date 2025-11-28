# HTTP Handlers

本目录包含所有 HTTP 请求处理器的实现。

## 概述

Handler 层负责：
1. 解析 HTTP 请求
2. 验证请求参数
3. 调用用例层执行业务逻辑
4. 转换响应格式
5. 返回 HTTP 响应

## 处理器列表

### 1. ChatHandler (chat_handler.go)

对话处理器，处理用户对话请求。

**端点:**
- `POST /chat` - 对话接口

**功能:**
- 普通对话响应
- 流式对话响应（SSE）
- 会话管理
- 意图路由

**请求示例:**
```json
{
  "query": "Python课程包含哪些内容？",
  "tenant_id": "tenant1",
  "session_id": "session123",
  "stream": false
}
```

**响应示例:**
```json
{
  "answer": "Python课程包含...",
  "route": "course",
  "sources": [
    {
      "content": "课程内容...",
      "score": 0.95,
      "metadata": {}
    }
  ],
  "session_id": "session123",
  "metadata": {
    "intent": "course",
    "confidence": 0.98,
    "duration_ms": 234
  }
}
```

**流式响应:**
当 `stream: true` 时，使用 Server-Sent Events (SSE) 协议：

```
event: message
data: {"content": "Python"}

event: message
data: {"content": "课程"}

event: done
data: {"metadata": {...}}
```

### 2. VectorHandler (vector_handler.go)

向量管理处理器，处理知识库向量的增删查操作。

**端点:**
- `POST /api/v1/vectors/items` - 添加向量
- `DELETE /api/v1/vectors/items` - 删除向量
- `GET /api/v1/vectors/count` - 获取向量数量
- `GET /api/v1/vectors/items/:id` - 获取单个向量

**添加向量请求示例:**
```json
{
  "texts": [
    "Python是一门编程语言",
    "Go语言适合高并发场景"
  ],
  "tenant_id": "tenant1",
  "metadata": {
    "category": "programming"
  }
}
```

**添加向量响应示例:**
```json
{
  "success": true,
  "document_ids": ["doc1", "doc2"],
  "count": 2,
  "message": "successfully added 2 vectors"
}
```

**删除向量请求示例:**
```json
{
  "ids": ["doc1", "doc2"],
  "tenant_id": "tenant1"
}
```

### 3. HealthHandler (health_handler.go)

健康检查处理器，用于监控系统状态。

**端点:**
- `GET /health` - 完整健康检查
- `GET /health/live` - 存活检查
- `GET /health/ready` - 就绪检查

**健康检查响应示例:**
```json
{
  "status": "healthy",
  "timestamp": "2024-11-28T10:00:00Z",
  "components": {
    "milvus": {
      "status": "healthy"
    },
    "database": {
      "status": "healthy"
    },
    "dashscope": {
      "status": "healthy"
    }
  }
}
```

**状态说明:**
- `healthy` - 所有组件正常
- `degraded` - 部分组件异常但服务可用
- `unhealthy` - 服务不可用

### 4. ModelHandler (model_handler.go)

模型管理处理器，用于查看和切换 AI 模型。

**端点:**
- `GET /models` - 列出所有可用模型
- `GET /models/current` - 获取当前使用的模型
- `POST /models/switch` - 切换模型
- `GET /models/:type/:name` - 获取模型信息

**列出模型响应示例:**
```json
{
  "success": true,
  "models": [
    {
      "type": "chat",
      "name": "qwen-turbo",
      "current": true
    },
    {
      "type": "chat",
      "name": "qwen-plus",
      "current": false
    }
  ],
  "current": {
    "chat": "qwen-turbo",
    "embedding": "text-embedding-v2"
  }
}
```

**切换模型请求示例:**
```json
{
  "type": "chat",
  "model": "qwen-plus"
}
```

## 使用方式

### 初始化处理器

```go
// 创建处理器
chatHandler := handler.NewChatHandler(chatUseCase)
vectorHandler := handler.NewVectorHandler(vectorUseCase)
healthHandler := handler.NewHealthHandler().
    WithMilvusCheck(milvusHealthCheck).
    WithDBCheck(dbHealthCheck)
modelHandler := handler.NewModelHandler("qwen-turbo", "text-embedding-v2")

// 注册路由
router := gin.Default()

// 对话路由
router.POST("/chat", chatHandler.HandleChat)

// 向量管理路由
vectorGroup := router.Group("/api/v1/vectors")
{
    vectorGroup.POST("/items", vectorHandler.HandleAddVectors)
    vectorGroup.DELETE("/items", vectorHandler.HandleDeleteVectors)
    vectorGroup.GET("/count", vectorHandler.HandleGetVectorCount)
    vectorGroup.GET("/items/:id", vectorHandler.HandleGetVector)
}

// 健康检查路由
healthGroup := router.Group("/health")
{
    healthGroup.GET("", healthHandler.HandleHealth)
    healthGroup.GET("/live", healthHandler.HandleLiveness)
    healthGroup.GET("/ready", healthHandler.HandleReadiness)
}

// 模型管理路由
modelGroup := router.Group("/models")
{
    modelGroup.GET("", modelHandler.HandleListModels)
    modelGroup.GET("/current", modelHandler.HandleGetCurrentModel)
    modelGroup.POST("/switch", modelHandler.HandleSwitchModel)
    modelGroup.GET("/:type/:name", modelHandler.HandleGetModelInfo)
}
```

## 错误处理

所有处理器使用统一的错误处理机制：

1. 使用 `c.Error()` 记录错误
2. 错误由 `ErrorHandler` 中间件统一处理
3. 返回标准化的错误响应

**错误响应格式:**
```json
{
  "code": 400,
  "message": "invalid request: query is required",
  "details": null,
  "trace_id": "req-123456"
}
```

## 中间件集成

处理器依赖以下中间件：

1. **TenantMiddleware** - 提取租户 ID
2. **AuthMiddleware** - API Key 验证（向量管理接口）
3. **SecurityMiddleware** - 敏感信息脱敏
4. **LoggingMiddleware** - 请求日志记录
5. **ErrorHandler** - 统一错误处理

## 测试

每个处理器都应该有对应的测试文件：

- `chat_handler_test.go`
- `vector_handler_test.go`
- `health_handler_test.go`
- `model_handler_test.go`

测试应该覆盖：
- 正常请求处理
- 参数验证
- 错误处理
- 边界情况

## 性能考虑

1. **流式响应**: 对于长时间运行的对话，使用流式响应减少首字节时间
2. **超时控制**: 使用 context 超时避免请求挂起
3. **并发处理**: Gin 框架自动处理并发请求
4. **资源清理**: 使用 defer 确保资源正确释放

## 安全考虑

1. **输入验证**: 所有输入都经过验证
2. **租户隔离**: 通过中间件确保租户数据隔离
3. **API Key**: 敏感接口需要 API Key 验证
4. **错误信息**: 不暴露内部实现细节

## 需求映射

- **需求 6.1**: Gin HTTP 服务器监听请求
- **需求 6.2**: 解析 JSON 请求体
- **需求 6.3**: 返回 JSON 响应
- **需求 6.4**: 流式响应支持（SSE）
- **需求 6.5**: 错误处理和状态码
- **需求 7.5**: 健康检查接口
- **需求 9.1**: 向量管理 API Key 验证
