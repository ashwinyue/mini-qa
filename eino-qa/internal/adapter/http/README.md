# HTTP Adapter Layer

HTTP 适配层实现了 Gin Web 框架的路由配置、服务器启动和优雅关闭功能。

## 目录结构

```
http/
├── router.go           # 路由配置
├── server.go           # 服务器启动和关闭
├── handler/            # HTTP 处理器
│   ├── chat_handler.go
│   ├── vector_handler.go
│   ├── health_handler.go
│   └── model_handler.go
└── middleware/         # HTTP 中间件
    ├── tenant.go
    ├── security.go
    ├── logging.go
    ├── auth.go
    └── error.go
```

## 核心组件

### 1. Router (router.go)

路由配置模块，负责设置所有 HTTP 路由和中间件。

**主要功能**:
- 配置全局中间件（租户识别、日志、安全、错误处理）
- 设置健康检查路由
- 设置对话接口路由
- 设置向量管理 API 路由
- 设置模型管理路由

**路由列表**:

| 方法 | 路径 | 描述 | 认证 |
|------|------|------|------|
| GET | / | 服务信息 | 否 |
| GET | /health | 健康检查 | 否 |
| GET | /health/live | 存活检查 | 否 |
| GET | /health/ready | 就绪检查 | 否 |
| POST | /chat | 对话接口 | 否 |
| POST | /api/v1/vectors/items | 添加向量 | 是 |
| DELETE | /api/v1/vectors/items | 删除向量 | 是 |
| GET | /api/v1/vectors/count | 获取向量数量 | 是 |
| GET | /api/v1/vectors/items/:id | 获取向量 | 是 |
| GET | /models | 列出可用模型 | 是 |
| GET | /models/current | 获取当前模型 | 是 |
| POST | /models/switch | 切换模型 | 是 |
| GET | /models/:type/:name | 获取模型信息 | 是 |

### 2. Server (server.go)

服务器启动和优雅关闭模块。

**主要功能**:
- 启动 HTTP 服务器
- 处理系统信号（SIGINT, SIGTERM）
- 优雅关闭服务器
- 配置超时和限制

**配置选项**:
- `Host`: 监听地址（默认: 0.0.0.0）
- `Port`: 监听端口（默认: 8080）
- `ReadTimeout`: 读取超时（默认: 30s）
- `WriteTimeout`: 写入超时（默认: 30s）
- `IdleTimeout`: 空闲超时（默认: 60s）
- `ShutdownTimeout`: 关闭超时（默认: 10s）
- `MaxHeaderBytes`: 最大请求头大小（默认: 1MB）

## 使用示例

### 基本使用

```go
package main

import (
    "eino-qa/internal/adapter/http"
    "eino-qa/internal/adapter/http/handler"
    "eino-qa/internal/adapter/http/middleware"
    "log"
)

func main() {
    // 创建 handlers
    chatHandler := handler.NewChatHandler(chatUseCase)
    vectorHandler := handler.NewVectorHandler(vectorUseCase)
    healthHandler := handler.NewHealthHandler()
    modelHandler := handler.NewModelHandler("qwen-turbo", "text-embedding-v2")

    // 创建 middlewares
    loggingMiddleware := middleware.NewLoggingMiddleware(logger)

    // 配置路由
    routerConfig := &http.RouterConfig{
        ChatHandler:        chatHandler,
        VectorHandler:      vectorHandler,
        HealthHandler:      healthHandler,
        ModelHandler:       modelHandler,
        TenantMiddleware:   middleware.TenantMiddleware(),
        SecurityMiddleware: middleware.SecurityMiddleware(),
        LoggingMiddleware:  loggingMiddleware.Handler(),
        ErrorMiddleware:    middleware.ErrorHandler(),
        AuthMiddleware:     middleware.AuthMiddleware("your-api-key"),
        Mode:               "release",
    }

    // 设置路由
    router := http.SetupRouter(routerConfig)

    // 创建服务器配置
    serverConfig := &http.ServerConfig{
        Host: "0.0.0.0",
        Port: 8080,
    }

    // 创建并运行服务器
    server := http.NewServer(router, serverConfig)
    if err := server.Run(); err != nil {
        log.Fatalf("Server error: %v", err)
    }
}
```

### 自定义配置

```go
// 使用默认配置
routerConfig := http.DefaultRouterConfig()

// 自定义 handlers
routerConfig.ChatHandler = handler.NewChatHandler(chatUseCase)
routerConfig.VectorHandler = handler.NewVectorHandler(vectorUseCase)

// 设置路由
router := http.SetupRouter(routerConfig)

// 使用默认服务器配置
serverConfig := http.DefaultServerConfig()
serverConfig.Port = 9090

// 创建服务器
server := http.NewServer(router, serverConfig)
```

### 手动启动和关闭

```go
// 创建服务器
server := http.NewServer(router, serverConfig)

// 在 goroutine 中启动
go func() {
    if err := server.Start(); err != nil {
        log.Printf("Server error: %v", err)
    }
}()

// 等待一段时间
time.Sleep(10 * time.Second)

// 手动关闭
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

if err := server.Shutdown(ctx); err != nil {
    log.Printf("Shutdown error: %v", err)
}
```

## 中间件顺序

中间件按以下顺序应用（从外到内）：

1. **Recovery**: Gin 内置的 panic 恢复中间件
2. **Tenant**: 租户识别中间件，提取租户 ID
3. **Logging**: 日志记录中间件，记录请求详情
4. **Security**: 安全脱敏中间件，处理敏感信息
5. **Error**: 错误处理中间件，统一错误响应格式

对于需要认证的路由，还会应用：

6. **Auth**: API Key 认证中间件

## 健康检查

系统提供三种健康检查端点：

### 1. 综合健康检查 (GET /health)

检查所有组件的健康状态，包括：
- Milvus 向量数据库
- SQLite 数据库
- DashScope API

响应示例：
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

### 2. 存活检查 (GET /health/live)

简单的存活检查，只要服务能响应就返回 200。

响应示例：
```json
{
  "status": "alive",
  "timestamp": "2024-11-28T10:00:00Z"
}
```

### 3. 就绪检查 (GET /health/ready)

检查服务是否准备好接收流量。

响应示例：
```json
{
  "status": "ready",
  "timestamp": "2024-11-28T10:00:00Z",
  "components": {
    "milvus": "ready",
    "database": "ready"
  }
}
```

## 优雅关闭

服务器支持优雅关闭，会：

1. 停止接收新请求
2. 等待现有请求完成（最多等待 ShutdownTimeout）
3. 关闭所有连接
4. 清理资源

优雅关闭由以下信号触发：
- `SIGINT` (Ctrl+C)
- `SIGTERM` (kill 命令)

## 错误处理

所有错误都通过错误处理中间件统一处理，返回格式：

```json
{
  "code": 400,
  "message": "invalid request: missing required field",
  "details": {},
  "trace_id": "550e8400-e29b-41d4-a716-446655440000"
}
```

错误类型：
- `400 Bad Request`: 请求参数错误
- `401 Unauthorized`: 认证失败
- `404 Not Found`: 资源不存在
- `500 Internal Server Error`: 服务器内部错误
- `502 Bad Gateway`: 外部服务调用失败
- `503 Service Unavailable`: 服务不可用

## 性能优化

### 超时配置

合理配置超时可以防止资源耗尽：

```go
serverConfig := &http.ServerConfig{
    ReadTimeout:  30 * time.Second,  // 读取请求超时
    WriteTimeout: 30 * time.Second,  // 写入响应超时
    IdleTimeout:  60 * time.Second,  // 空闲连接超时
}
```

### 请求大小限制

限制请求头大小可以防止恶意请求：

```go
serverConfig.MaxHeaderBytes = 1 << 20  // 1 MB
```

### 连接池

Gin 默认使用 Go 的 net/http 包，它会自动管理连接池。

## 测试

### 单元测试

```go
func TestRouter(t *testing.T) {
    // 创建测试路由
    config := http.DefaultRouterConfig()
    config.ChatHandler = handler.NewChatHandler(mockChatUseCase)
    router := http.SetupRouter(config)

    // 创建测试请求
    w := httptest.NewRecorder()
    req, _ := http.NewRequest("POST", "/chat", bytes.NewBuffer([]byte(`{"query":"test"}`)))
    req.Header.Set("Content-Type", "application/json")

    // 执行请求
    router.ServeHTTP(w, req)

    // 验证响应
    assert.Equal(t, 200, w.Code)
}
```

### 集成测试

```go
func TestServerIntegration(t *testing.T) {
    // 创建测试服务器
    router := http.SetupRouter(testConfig)
    server := http.NewServer(router, &http.ServerConfig{
        Host: "localhost",
        Port: 0, // 随机端口
    })

    // 启动服务器
    go server.Start()
    defer server.Shutdown(context.Background())

    // 发送测试请求
    resp, err := http.Post("http://localhost:8080/chat", "application/json", body)
    assert.NoError(t, err)
    assert.Equal(t, 200, resp.StatusCode)
}
```

## 需求映射

- **需求 6.1**: Gin HTTP 服务器在配置的端口上监听请求 ✓
- **需求 6.2**: 解析 JSON 请求体并提取查询参数 ✓
- **需求 6.3**: 返回 JSON 格式的响应 ✓
- **需求 6.4**: 支持流式响应（SSE） ✓
- **需求 6.5**: 返回适当的 HTTP 状态码和错误信息 ✓
- **需求 7.5**: 健康检查接口 ✓

## 注意事项

1. **生产环境**: 确保设置 `Mode` 为 `"release"`
2. **API Key**: 妥善保管 API Key，不要硬编码在代码中
3. **超时配置**: 根据实际业务需求调整超时时间
4. **日志级别**: 生产环境建议使用 INFO 级别
5. **CORS**: 如果需要跨域访问，需要添加 CORS 中间件
6. **TLS**: 生产环境建议使用 HTTPS
