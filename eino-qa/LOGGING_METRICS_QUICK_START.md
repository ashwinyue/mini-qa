# 日志和指标系统快速开始

## 快速集成

### 1. 创建日志和指标实例

```go
package main

import (
    "eino-qa/internal/infrastructure/logger"
    "eino-qa/internal/infrastructure/metrics"
)

func main() {
    // 创建日志系统
    log, err := logger.New(logger.Config{
        Level:  "info",
        Format: "json",
        Output: "stdout",
    })
    if err != nil {
        panic(err)
    }

    // 创建指标系统
    metricsCollector := metrics.New(metrics.DefaultConfig())
}
```

### 2. 集成到 HTTP 服务器

```go
import (
    "eino-qa/internal/adapter/http"
    "eino-qa/internal/adapter/http/handler"
    "eino-qa/internal/adapter/http/middleware"
)

// 创建中间件
loggingMW := middleware.NewLoggingMiddleware(log)
metricsMW := middleware.NewMetricsMiddleware(metricsCollector)

// 创建健康检查处理器
healthHandler := handler.NewHealthHandler().
    WithMetricsProvider(metricsCollector)

// 配置路由
routerConfig := &http.RouterConfig{
    Mode:              gin.ReleaseMode,
    HealthHandler:     healthHandler,
    LoggingMiddleware: loggingMW.Handler(),
    MetricsMiddleware: metricsMW.Handler(),
    // ... 其他配置
}

router := http.SetupRouter(routerConfig)
router.Run(":8080")
```

### 3. 在业务代码中使用日志

```go
type MyService struct {
    logger logger.Logger
}

func (s *MyService) DoSomething(ctx context.Context) error {
    // 记录信息
    s.logger.Info(ctx, "operation started", map[string]interface{}{
        "operation": "do_something",
    })

    // 记录外部调用
    startTime := time.Now()
    result, err := s.callExternalAPI(ctx)
    duration := time.Since(startTime)

    s.logger.Info(ctx, "external api called", map[string]interface{}{
        "duration_ms": duration.Milliseconds(),
        "success":     err == nil,
    })

    // 记录错误
    if err != nil {
        s.logger.Error(ctx, "operation failed", map[string]interface{}{
            "error": err.Error(),
        })
        return err
    }

    return nil
}
```

## 可用的 API 端点

### 健康检查（包含指标）
```bash
curl http://localhost:8080/health
```

响应:
```json
{
  "status": "healthy",
  "timestamp": "2024-11-28T10:00:00Z",
  "components": {
    "milvus": {"status": "healthy"},
    "database": {"status": "healthy"}
  },
  "metrics": {
    "total_requests": 1000,
    "avg_response_time_ms": 125.5
  }
}
```

### 详细指标查询
```bash
curl http://localhost:8080/health/metrics
```

响应:
```json
{
  "total_requests": 1000,
  "success_requests": 950,
  "client_errors": 30,
  "server_errors": 20,
  "avg_response_time_ms": 125.5,
  "p95_response_time_ms": 250.0,
  "p99_response_time_ms": 450.0,
  "route_stats": {
    "/chat": {
      "count": 500,
      "avg_duration_ms": 145.8,
      "success_count": 475,
      "error_count": 25
    }
  }
}
```

### 存活检查
```bash
curl http://localhost:8080/health/live
```

### 就绪检查
```bash
curl http://localhost:8080/health/ready
```

## 日志格式示例

### 请求开始日志
```json
{
  "timestamp": "2024-11-28T10:00:00Z",
  "level": "info",
  "message": "request started",
  "request_id": "550e8400-e29b-41d4-a716-446655440000",
  "tenant_id": "default",
  "method": "POST",
  "path": "/chat",
  "client_ip": "127.0.0.1"
}
```

### 请求完成日志
```json
{
  "timestamp": "2024-11-28T10:00:01Z",
  "level": "info",
  "message": "request completed successfully",
  "request_id": "550e8400-e29b-41d4-a716-446655440000",
  "tenant_id": "default",
  "method": "POST",
  "path": "/chat",
  "status_code": 200,
  "duration_ms": 125,
  "request_size": 256,
  "response_size": 512
}
```

### 错误日志
```json
{
  "timestamp": "2024-11-28T10:00:02Z",
  "level": "error",
  "message": "database query failed",
  "request_id": "550e8400-e29b-41d4-a716-446655440000",
  "tenant_id": "default",
  "error": "connection timeout",
  "query": "SELECT * FROM orders"
}
```

## 配置选项

### 日志配置

```go
type Config struct {
    Level    string // "debug", "info", "warn", "error"
    Format   string // "json", "text"
    Output   string // "stdout", "file"
    FilePath string // 日志文件路径（当 Output="file" 时）
}
```

### 指标配置

```go
type Config struct {
    MaxResponseTimeSamples int // 保留的响应时间样本数（用于百分位数计算）
}
```

## 常见使用场景

### 场景 1: 记录业务操作
```go
logger.Info(ctx, "user action", map[string]interface{}{
    "action":  "login",
    "user_id": "12345",
})
```

### 场景 2: 记录外部服务调用
```go
startTime := time.Now()
result, err := externalService.Call(ctx, params)
duration := time.Since(startTime)

logger.Info(ctx, "external service called", map[string]interface{}{
    "service":     "milvus",
    "operation":   "search",
    "duration_ms": duration.Milliseconds(),
    "success":     err == nil,
})
```

### 场景 3: 记录错误详情
```go
if err != nil {
    logger.Error(ctx, "operation failed", map[string]interface{}{
        "error":     err.Error(),
        "operation": "rag_retrieval",
        "query":     query,
        "retry":     retryCount,
    })
}
```

### 场景 4: 创建组件级 Logger
```go
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
```

## 监控和告警建议

### 关键指标监控
- 总请求数
- 错误率 (client_errors + server_errors) / total_requests
- P95 响应时间
- 各路由的平均响应时间

### 告警阈值建议
- 错误率 > 5%
- P95 响应时间 > 500ms
- 服务端错误数 > 10/分钟

## 性能影响

- 日志记录: ~5-10 微秒/条
- 指标收集: ~1-2 微秒/请求
- 内存占用: ~100KB + (样本数 × 8 字节)

## 更多信息

- [指标系统详细文档](internal/infrastructure/metrics/README.md)
- [日志系统详细文档](internal/infrastructure/logger/README.md)
- [完整示例](examples/logging_metrics_example.go)
- [任务总结](TASK_13_SUMMARY.md)
