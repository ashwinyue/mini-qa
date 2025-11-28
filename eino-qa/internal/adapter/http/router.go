package http

import (
	"eino-qa/internal/adapter/http/handler"
	"eino-qa/internal/adapter/http/middleware"

	"github.com/gin-gonic/gin"
)

// RouterConfig 路由配置
type RouterConfig struct {
	// Handlers
	ChatHandler   *handler.ChatHandler
	VectorHandler *handler.VectorHandler
	HealthHandler *handler.HealthHandler
	ModelHandler  *handler.ModelHandler

	// Middlewares
	TenantMiddleware   gin.HandlerFunc
	SecurityMiddleware gin.HandlerFunc
	LoggingMiddleware  gin.HandlerFunc
	MetricsMiddleware  gin.HandlerFunc
	ErrorMiddleware    gin.HandlerFunc
	AuthMiddleware     gin.HandlerFunc

	// Config
	Mode   string // "debug", "release", "test"
	APIKey string // API Key for protected endpoints
}

// SetupRouter 配置路由
// 需求: 6.1 - HTTP API 接口
func SetupRouter(config *RouterConfig) *gin.Engine {
	// 设置 Gin 模式
	if config.Mode != "" {
		gin.SetMode(config.Mode)
	}

	// 创建 Gin 引擎
	router := gin.New()

	// 全局中间件（按顺序应用）
	// 1. 恢复中间件（处理 panic）
	router.Use(gin.Recovery())

	// 2. 租户识别中间件
	if config.TenantMiddleware != nil {
		router.Use(config.TenantMiddleware)
	}

	// 3. 日志记录中间件
	if config.LoggingMiddleware != nil {
		router.Use(config.LoggingMiddleware)
	}

	// 4. 指标收集中间件
	// 需求: 7.4 - 统计各类请求的数量和平均响应时间
	if config.MetricsMiddleware != nil {
		router.Use(config.MetricsMiddleware)
	}

	// 5. 安全脱敏中间件
	if config.SecurityMiddleware != nil {
		router.Use(config.SecurityMiddleware)
	}

	// 6. 错误处理中间件（应该在最后）
	if config.ErrorMiddleware != nil {
		router.Use(config.ErrorMiddleware)
	}

	// 健康检查路由（不需要认证）
	// 需求: 7.5 - 健康检查接口
	if config.HealthHandler != nil {
		healthGroup := router.Group("/health")
		{
			healthGroup.GET("", config.HealthHandler.HandleHealth)
			healthGroup.GET("/live", config.HealthHandler.HandleLiveness)
			healthGroup.GET("/ready", config.HealthHandler.HandleReadiness)
			healthGroup.GET("/metrics", config.HealthHandler.HandleMetrics)
		}
	}

	// 对话接口（不需要认证，但需要租户识别）
	// 需求: 6.1, 6.2, 6.3, 6.4, 6.5
	if config.ChatHandler != nil {
		router.POST("/chat", config.ChatHandler.HandleChat)
	}

	// API v1 路由组（需要 API Key 认证）
	apiV1 := router.Group("/api/v1")
	if config.AuthMiddleware != nil {
		apiV1.Use(config.AuthMiddleware)
	}
	{
		// 向量管理接口
		// 需求: 9.1, 9.2, 9.3, 9.4, 9.5
		if config.VectorHandler != nil {
			vectorGroup := apiV1.Group("/vectors")
			{
				vectorGroup.POST("/items", config.VectorHandler.HandleAddVectors)
				vectorGroup.DELETE("/items", config.VectorHandler.HandleDeleteVectors)
				vectorGroup.GET("/count", config.VectorHandler.HandleGetVectorCount)
				vectorGroup.GET("/items/:id", config.VectorHandler.HandleGetVector)
			}
		}
	}

	// 模型管理接口（需要 API Key 认证）
	if config.ModelHandler != nil {
		modelsGroup := router.Group("/models")
		if config.AuthMiddleware != nil {
			modelsGroup.Use(config.AuthMiddleware)
		}
		{
			modelsGroup.GET("", config.ModelHandler.HandleListModels)
			modelsGroup.GET("/current", config.ModelHandler.HandleGetCurrentModel)
			modelsGroup.POST("/switch", config.ModelHandler.HandleSwitchModel)
			modelsGroup.GET("/info/:type/:name", config.ModelHandler.HandleGetModelInfo)
		}
	}

	// 根路径
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"service": "eino-qa-system",
			"version": "1.0.0",
			"status":  "running",
		})
	})

	return router
}

// DefaultRouterConfig 创建默认路由配置
func DefaultRouterConfig() *RouterConfig {
	// 创建默认的安全中间件（脱敏常见敏感字段）
	securityMiddleware := middleware.NewSecurityMiddleware([]string{
		"password", "token", "api_key", "secret",
	})

	return &RouterConfig{
		Mode:               gin.ReleaseMode,
		TenantMiddleware:   middleware.TenantMiddleware(),
		SecurityMiddleware: securityMiddleware.Handler(),
		ErrorMiddleware:    middleware.ErrorHandler(),
	}
}
