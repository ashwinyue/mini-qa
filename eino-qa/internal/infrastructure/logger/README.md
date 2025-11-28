# 日志系统 (Logger System)

## 概述

日志系统提供结构化日志记录功能，支持上下文感知、多级别日志和灵活的输出配置。

## 功能特性

### 1. 结构化日志
- JSON 格式输出
- 文本格式输出
- 自定义字段支持

### 2. 上下文感知
- 自动提取 trace_id
- 自动提取 tenant_id
- 自动提取 session_id
- 自动提取 request_id

### 3. 多级别日志
- Debug
- Info
- Warn
- Error

### 4. 灵活输出
- 标准输出 (stdout)
- 文件输出
- 自动创建日志目录

## 使用方法

### 创建日志实例

```go
import "eino-qa/internal/infrastructure/logger"

// JSON 格式输出到标准输出
config := logger.Config{
    Level:  "info",
    Format: "json",
    Output: "stdout",
}
log, err := logger.New(config)
if err != nil {
    panic(err)
}

// 文本格式输出到文件
config := logger.Config{
    Level:    "debug",
    Format:   "text",
    Output:   "file",
    FilePath: "./logs/app.log",
}
log, err := logger.New(config)
```

### 记录日志

```go
import "context"

ctx := context.Background()

// Info 级别
log.Info(ctx, "user logged in", map[string]interface{}{
    "user_id": "12345",
    "action":  "login",
})

// Error 级别
log.Error(ctx, "failed to connect to database", map[string]interface{}{
    "error":    err.Error(),
    "database": "milvus",
    "retry":    3,
})

// Debug 级别
log.Debug(ctx, "processing request", map[string]interface{}{
    "step":     "validation",
    "duration": "5ms",
})

// Warn 级别
log.Warn(ctx, "rate limit approaching", map[string]interface{}{
    "current": 950,
    "limit":   1000,
})
```

### 上下文字段自动提取

```go
// 创建带有上下文信息的 context
ctx := context.WithValue(context.Background(), "trace_id", "abc123")
ctx = context.WithValue(ctx, "tenant_id", "tenant1")
ctx = context.WithValue(ctx, "request_id", "req456")

// 记录日志时自动包含上下文字段
log.Info(ctx, "processing request", map[string]interface{}{
    "action": "query",
})

// 输出:
// {
//   "timestamp": "2024-11-28T10:00:00Z",
//   "level": "info",
//   "message": "processing request",
//   "trace_id": "abc123",
//   "tenant_id": "tenant1",
//   "request_id": "req456",
//   "action": "query"
// }
```

### 预设字段

```go
// 创建带有预设字段的 logger
componentLogger := log.WithFields(map[string]interface{}{
    "component": "rag_retriever",
    "version":   "1.0.0",
})

// 所有日志都会包含预设字段
componentLogger.Info(ctx, "retrieval started", map[string]interface{}{
    "query": "Python课程",
})

// 输出:
// {
//   "timestamp": "2024-11-28T10:00:00Z",
//   "level": "info",
//   "message": "retrieval started",
//   "component": "rag_retriever",
//   "version": "1.0.0",
//   "query": "Python课程"
// }
```

## 日志格式

### JSON 格式

```json
{
  "timestamp": "2024-11-28T10:00:00Z",
  "level": "info",
  "message": "request completed successfully",
  "request_id": "550e8400-e29b-41d4-a716-446655440000",
  "tenant_id": "default",
  "method": "POST",
  "path": "/chat",
  "status_code": 200,
  "duration_ms": 125
}
```

### 文本格式

```
2024-11-28T10:00:00Z INFO request completed successfully request_id=550e8400-e29b-41d4-a716-446655440000 tenant_id=default method=POST path=/chat status_code=200 duration_ms=125
```

## 与 HTTP 中间件集成

```go
import (
    "eino-qa/internal/adapter/http/middleware"
    "eino-qa/internal/infrastructure/logger"
)

// 创建日志实例
log, err := logger.New(logger.Config{
    Level:  "info",
    Format: "json",
    Output: "stdout",
})

// 创建日志中间件
loggingMiddleware := middleware.NewLoggingMiddleware(log)

// 在 Gin 中使用
router.Use(loggingMiddleware.Handler())
```

## 日志级别

### Debug
用于详细的调试信息，通常只在开发环境启用。

```go
log.Debug(ctx, "cache hit", map[string]interface{}{
    "key":   "user:12345",
    "value": "cached_data",
})
```

### Info
用于一般信息性消息，记录正常的业务流程。

```go
log.Info(ctx, "request started", map[string]interface{}{
    "method": "POST",
    "path":   "/chat",
})
```

### Warn
用于警告信息，表示潜在问题但不影响正常运行。

```go
log.Warn(ctx, "slow query detected", map[string]interface{}{
    "duration_ms": 5000,
    "query":       "SELECT * FROM orders",
})
```

### Error
用于错误信息，表示发生了错误但系统仍可继续运行。

```go
log.Error(ctx, "failed to retrieve documents", map[string]interface{}{
    "error":  err.Error(),
    "query":  "Python课程",
    "tenant": "tenant1",
})
```

## 最佳实践

### 1. 使用结构化字段

❌ 不推荐:
```go
log.Info(ctx, fmt.Sprintf("User %s logged in from %s", userID, ip), nil)
```

✅ 推荐:
```go
log.Info(ctx, "user logged in", map[string]interface{}{
    "user_id": userID,
    "ip":      ip,
})
```

### 2. 记录关键业务事件

```go
// 需求: 7.1 - 记录请求 ID、租户 ID、查询内容和处理时长
log.Info(ctx, "chat request completed", map[string]interface{}{
    "query":       "Python课程包含哪些内容？",
    "route":       "course",
    "duration_ms": 125,
})
```

### 3. 记录外部调用

```go
// 需求: 7.3 - 记录调用参数和响应时间
startTime := time.Now()
result, err := milvusClient.Search(ctx, query)
duration := time.Since(startTime)

log.Info(ctx, "milvus search completed", map[string]interface{}{
    "collection":  "knowledge_base",
    "query_size":  len(query),
    "result_size": len(result),
    "duration_ms": duration.Milliseconds(),
})
```

### 4. 记录错误详情

```go
// 需求: 7.2 - 记录错误堆栈和上下文信息
if err != nil {
    log.Error(ctx, "database query failed", map[string]interface{}{
        "error":     err.Error(),
        "query":     sqlQuery,
        "tenant_id": tenantID,
        "retry":     retryCount,
    })
}
```

### 5. 使用组件级 Logger

```go
// 为每个组件创建专用 logger
type RAGRetriever struct {
    logger logger.Logger
}

func NewRAGRetriever(baseLogger logger.Logger) *RAGRetriever {
    return &RAGRetriever{
        logger: baseLogger.WithFields(map[string]interface{}{
            "component": "rag_retriever",
        }),
    }
}

func (r *RAGRetriever) Retrieve(ctx context.Context, query string) {
    r.logger.Info(ctx, "starting retrieval", map[string]interface{}{
        "query": query,
    })
}
```

## 配置示例

### 开发环境

```go
config := logger.Config{
    Level:  "debug",
    Format: "text",
    Output: "stdout",
}
```

### 生产环境

```go
config := logger.Config{
    Level:    "info",
    Format:   "json",
    Output:   "file",
    FilePath: "/var/log/eino-qa/app.log",
}
```

### 测试环境

```go
config := logger.Config{
    Level:  "warn",
    Format: "json",
    Output: "stdout",
}
```

## 性能考虑

- JSON 格式化: ~5-10 微秒/条
- 文本格式化: ~2-5 微秒/条
- 文件写入: 异步批量写入，对性能影响最小
- 上下文提取: ~1 微秒

## 需求映射

- **需求 7.1**: 记录请求 ID、租户 ID、查询内容和处理时长
- **需求 7.2**: 记录错误堆栈和上下文信息
- **需求 7.3**: 记录调用参数和响应时间

## 相关文档

- [指标系统](../metrics/README.md)
- [HTTP 中间件](../../adapter/http/middleware/README.md)
- [安全中间件](../../adapter/http/middleware/security.go)
