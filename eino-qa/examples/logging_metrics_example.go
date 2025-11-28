package main

import (
	"context"
	"eino-qa/internal/adapter/http"
	"eino-qa/internal/adapter/http/handler"
	"eino-qa/internal/adapter/http/middleware"
	"eino-qa/internal/infrastructure/logger"
	"eino-qa/internal/infrastructure/metrics"
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

// 本示例展示如何集成日志和指标系统
// 需求: 7.1, 7.2, 7.3, 7.4, 7.5

func main() {
	// 1. 创建日志系统
	// 需求: 7.1 - 记录请求 ID、租户 ID、查询内容和处理时长
	logConfig := logger.Config{
		Level:  "info",
		Format: "json",
		Output: "stdout",
	}
	appLogger, err := logger.New(logConfig)
	if err != nil {
		log.Fatalf("Failed to create logger: %v", err)
	}

	// 2. 创建指标收集器
	// 需求: 7.4 - 统计各类请求的数量和平均响应时间
	metricsCollector := metrics.New(metrics.DefaultConfig())

	// 3. 创建中间件
	loggingMiddleware := middleware.NewLoggingMiddleware(appLogger)
	metricsMiddleware := middleware.NewMetricsMiddleware(metricsCollector)

	// 4. 创建健康检查处理器（带指标）
	// 需求: 7.5 - 返回系统状态和关键指标快照
	healthHandler := handler.NewHealthHandler().
		WithMetricsProvider(metricsCollector).
		WithMilvusCheck(func(ctx context.Context) error {
			// 模拟 Milvus 健康检查
			return nil
		}).
		WithDBCheck(func(ctx context.Context) error {
			// 模拟数据库健康检查
			return nil
		})

	// 5. 配置路由
	routerConfig := &http.RouterConfig{
		Mode:              gin.DebugMode,
		HealthHandler:     healthHandler,
		TenantMiddleware:  middleware.TenantMiddleware(),
		LoggingMiddleware: loggingMiddleware.Handler(),
		MetricsMiddleware: metricsMiddleware.Handler(),
		SecurityMiddleware: middleware.NewSecurityMiddleware([]string{
			"password", "token", "api_key",
		}).Handler(),
		ErrorMiddleware: middleware.ErrorHandler(),
	}

	router := http.SetupRouter(routerConfig)

	// 6. 添加测试端点
	router.GET("/test", func(c *gin.Context) {
		ctx := c.Request.Context()

		// 记录业务日志
		// 需求: 7.3 - 记录调用参数和响应时间
		appLogger.Info(ctx, "processing test request", map[string]interface{}{
			"action": "test",
		})

		time.Sleep(100 * time.Millisecond) // 模拟处理

		c.JSON(200, gin.H{
			"message": "test successful",
		})
	})

	// 7. 启动服务器
	fmt.Println("Server starting on :8080")
	fmt.Println("Endpoints:")
	fmt.Println("  GET  /health         - 健康检查（包含指标）")
	fmt.Println("  GET  /health/metrics - 详细指标")
	fmt.Println("  GET  /health/live    - 存活检查")
	fmt.Println("  GET  /health/ready   - 就绪检查")
	fmt.Println("  GET  /test           - 测试端点")

	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// 示例输出：
//
// 1. 结构化日志输出（JSON 格式）:
// {
//   "timestamp": "2024-11-28T10:00:00Z",
//   "level": "info",
//   "message": "request started",
//   "request_id": "550e8400-e29b-41d4-a716-446655440000",
//   "tenant_id": "default",
//   "method": "GET",
//   "path": "/test",
//   "client_ip": "127.0.0.1"
// }
//
// 2. 请求完成日志:
// {
//   "timestamp": "2024-11-28T10:00:00Z",
//   "level": "info",
//   "message": "request completed successfully",
//   "request_id": "550e8400-e29b-41d4-a716-446655440000",
//   "tenant_id": "default",
//   "status_code": 200,
//   "duration_ms": 105
// }
//
// 3. 指标查询响应 (GET /health/metrics):
// {
//   "total_requests": 100,
//   "success_requests": 95,
//   "client_errors": 3,
//   "server_errors": 2,
//   "avg_response_time_ms": 125.5,
//   "p95_response_time_ms": 250.0,
//   "p99_response_time_ms": 450.0,
//   "route_stats": {
//     "/test": {
//       "count": 50,
//       "avg_duration_ms": 105.2,
//       "min_duration_ms": 95,
//       "max_duration_ms": 150,
//       "success_count": 50,
//       "error_count": 0
//     },
//     "/chat": {
//       "count": 50,
//       "avg_duration_ms": 145.8,
//       "min_duration_ms": 100,
//       "max_duration_ms": 500,
//       "success_count": 45,
//       "error_count": 5
//     }
//   },
//   "error_stats": {
//     "/chat:bind_error": 3,
//     "/chat:private_error": 2
//   },
//   "start_time": "2024-11-28T09:00:00Z",
//   "last_update": "2024-11-28T10:00:00Z"
// }
