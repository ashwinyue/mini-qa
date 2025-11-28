# 指标系统 (Metrics System)

## 概述

指标系统提供了请求统计、性能监控和系统健康状态的实时数据收集功能。

## 功能特性

### 1. 请求统计
- 总请求数
- 成功请求数 (2xx)
- 客户端错误数 (4xx)
- 服务端错误数 (5xx)

### 2. 性能指标
- 平均响应时间
- P95 响应时间
- P99 响应时间
- 按路由分类的详细统计

### 3. 错误追踪
- 按路由和错误类型分类的错误统计
- 错误趋势分析

## 使用方法

### 创建指标收集器

```go
import "eino-qa/internal/infrastructure/metrics"

// 使用默认配置
collector := metrics.New(metrics.DefaultConfig())

// 自定义配置
config := metrics.Config{
    MaxResponseTimeSamples: 5000, // 保留 5000 个响应时间样本
}
collector := metrics.New(config)
```

### 记录请求

```go
// 记录请求指标
route := "/api/v1/chat"
statusCode := 200
duration := 150 * time.Millisecond

collector.RecordRequest(route, statusCode, duration)
```

### 记录错误

```go
// 记录错误
route := "/api/v1/chat"
errorType := "bind_error"

collector.RecordError(route, errorType)
```

### 获取统计信息

```go
// 获取完整统计信息
stats := collector.GetStats()

fmt.Printf("总请求数: %d\n", stats.TotalRequests)
fmt.Printf("平均响应时间: %.2f ms\n", stats.AvgResponseTime)
fmt.Printf("P95 响应时间: %.2f ms\n", stats.P95ResponseTime)

// 查看路由统计
for route, routeStat := range stats.RouteStats {
    fmt.Printf("路由 %s: %d 请求, 平均 %.2f ms\n", 
        route, routeStat.Count, routeStat.AvgDuration)
}
```

### 重置统计

```go
// 重置所有统计信息
collector.Reset()
```

## 与 HTTP 中间件集成

```go
import (
    "eino-qa/internal/adapter/http/middleware"
    "eino-qa/internal/infrastructure/metrics"
)

// 创建指标收集器
collector := metrics.New(metrics.DefaultConfig())

// 创建指标中间件
metricsMiddleware := middleware.NewMetricsMiddleware(collector)

// 在 Gin 中使用
router.Use(metricsMiddleware.Handler())
```

## 与健康检查集成

```go
import (
    "eino-qa/internal/adapter/http/handler"
    "eino-qa/internal/infrastructure/metrics"
)

// 创建指标收集器
collector := metrics.New(metrics.DefaultConfig())

// 创建健康检查处理器并注入指标提供者
healthHandler := handler.NewHealthHandler().
    WithMetricsProvider(collector)

// 健康检查响应将包含指标信息
// GET /health 返回:
// {
//   "status": "healthy",
//   "timestamp": "2024-11-28T10:00:00Z",
//   "components": {...},
//   "metrics": {
//     "total_requests": 1000,
//     "avg_response_time_ms": 125.5,
//     ...
//   }
// }
```

## API 端点

### GET /health/metrics

返回详细的系统指标信息。

**响应示例:**

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
      "min_duration_ms": 100,
      "max_duration_ms": 500,
      "success_count": 475,
      "error_count": 25
    },
    "/api/v1/vectors/items": {
      "count": 300,
      "avg_duration_ms": 95.2,
      "min_duration_ms": 50,
      "max_duration_ms": 200,
      "success_count": 295,
      "error_count": 5
    }
  },
  "error_stats": {
    "/chat:bind_error": 15,
    "/chat:private_error": 10,
    "/api/v1/vectors/items:auth_error": 5
  },
  "start_time": "2024-11-28T09:00:00Z",
  "last_update": "2024-11-28T10:00:00Z"
}
```

## 性能考虑

### 内存使用

指标系统在内存中保留响应时间样本用于计算百分位数。默认保留 10,000 个样本，可以通过配置调整：

```go
config := metrics.Config{
    MaxResponseTimeSamples: 5000, // 减少内存使用
}
```

### 并发安全

所有指标操作都是线程安全的，使用读写锁保护共享数据。

### 性能影响

- 记录请求: ~1-2 微秒
- 获取统计: ~10-50 微秒（取决于路由数量）
- 内存占用: ~100KB + (样本数 × 8 字节)

## 最佳实践

1. **合理设置样本数**: 根据流量大小调整 `MaxResponseTimeSamples`
2. **定期导出指标**: 对于长期运行的服务，定期导出并重置指标
3. **监控关键路由**: 重点关注高流量路由的性能指标
4. **设置告警**: 基于 P95/P99 响应时间设置性能告警

## 需求映射

- **需求 7.4**: 统计各类请求的数量和平均响应时间
- **需求 7.5**: 返回系统状态和关键指标快照

## 相关文档

- [日志系统](../logger/README.md)
- [HTTP 中间件](../../adapter/http/middleware/README.md)
- [健康检查](../../adapter/http/handler/health_handler.go)
