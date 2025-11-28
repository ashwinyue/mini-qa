# HTTP Adapter 快速开始

## 最简单的使用方式

```go
package main

import (
    "eino-qa/internal/adapter/http"
    "log"
)

func main() {
    // 1. 使用默认配置创建路由
    routerConfig := http.DefaultRouterConfig()
    
    // 2. 添加你的 handlers（可选）
    // routerConfig.ChatHandler = handler.NewChatHandler(chatUseCase)
    // routerConfig.VectorHandler = handler.NewVectorHandler(vectorUseCase)
    
    // 3. 设置路由
    router := http.SetupRouter(routerConfig)
    
    // 4. 使用默认服务器配置
    serverConfig := http.DefaultServerConfig()
    
    // 5. 创建并运行服务器
    server := http.NewServer(router, serverConfig)
    if err := server.Run(); err != nil {
        log.Fatalf("Server error: %v", err)
    }
}
```

## 自定义端口

```go
serverConfig := http.DefaultServerConfig()
serverConfig.Port = 9090  // 修改端口

server := http.NewServer(router, serverConfig)
server.Run()
```

## 添加 API Key 认证

```go
import "eino-qa/internal/adapter/http/middleware"

routerConfig := http.DefaultRouterConfig()

// 创建认证中间件
authMiddleware := middleware.NewAuthMiddleware([]string{
    "your-api-key-1",
    "your-api-key-2",
})
routerConfig.AuthMiddleware = authMiddleware.Handler()

router := http.SetupRouter(routerConfig)
```

## 添加日志记录

```go
import "eino-qa/internal/adapter/http/middleware"

// 创建日志中间件
loggingMiddleware := middleware.NewLoggingMiddleware(logger)
routerConfig.LoggingMiddleware = loggingMiddleware.Handler()

router := http.SetupRouter(routerConfig)
```

## 测试健康检查

启动服务器后，访问：

```bash
# 综合健康检查
curl http://localhost:8080/health

# 存活检查
curl http://localhost:8080/health/live

# 就绪检查
curl http://localhost:8080/health/ready
```

## 测试对话接口

```bash
curl -X POST http://localhost:8080/chat \
  -H "Content-Type: application/json" \
  -H "X-Tenant-ID: tenant1" \
  -d '{"query": "你好"}'
```

## 测试向量管理接口（需要 API Key）

```bash
# 添加向量
curl -X POST http://localhost:8080/api/v1/vectors/items \
  -H "Content-Type: application/json" \
  -H "X-API-Key: your-api-key" \
  -H "X-Tenant-ID: tenant1" \
  -d '{"texts": ["文本1", "文本2"]}'

# 获取向量数量
curl http://localhost:8080/api/v1/vectors/count \
  -H "X-API-Key: your-api-key" \
  -H "X-Tenant-ID: tenant1"
```

## 优雅关闭

服务器会自动处理 Ctrl+C (SIGINT) 和 kill (SIGTERM) 信号：

```bash
# 启动服务器
./server

# 按 Ctrl+C 优雅关闭
# 或者
kill -TERM <pid>
```

## 常见问题

### Q: 如何修改超时时间？

```go
serverConfig := &http.ServerConfig{
    Host:         "0.0.0.0",
    Port:         8080,
    ReadTimeout:  60 * time.Second,  // 修改读取超时
    WriteTimeout: 60 * time.Second,  // 修改写入超时
}
```

### Q: 如何禁用某个中间件？

```go
routerConfig := http.DefaultRouterConfig()
routerConfig.SecurityMiddleware = nil  // 禁用安全中间件
```

### Q: 如何添加自定义中间件？

```go
customMiddleware := func(c *gin.Context) {
    // 你的逻辑
    c.Next()
}

router := http.SetupRouter(routerConfig)
router.Use(customMiddleware)  // 添加到所有路由
```

### Q: 如何在生产环境运行？

```go
routerConfig := http.DefaultRouterConfig()
routerConfig.Mode = gin.ReleaseMode  // 设置为 release 模式

router := http.SetupRouter(routerConfig)
```

## 下一步

查看完整文档：[README.md](./README.md)
