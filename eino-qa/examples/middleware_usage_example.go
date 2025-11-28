package main

import (
	"eino-qa/internal/adapter/http/middleware"
	"eino-qa/internal/infrastructure/config"
	"eino-qa/internal/infrastructure/logger"

	"github.com/gin-gonic/gin"
)

// 这是一个示例，展示如何使用所有中间件
func main() {
	// 加载配置
	cfg, err := config.Load("config/config.yaml")
	if err != nil {
		panic(err)
	}

	// 初始化日志
	log, err := logger.New(logger.Config{
		Level:  cfg.Logging.Level,
		Format: cfg.Logging.Format,
		Output: cfg.Logging.Output,
	})
	if err != nil {
		panic(err)
	}

	// 创建 Gin 路由器
	router := gin.New()

	// 1. 错误处理中间件（最外层）
	router.Use(middleware.ErrorHandler())

	// 2. 日志记录中间件
	// 注意：需要适配器将 logger.Logger 转换为 middleware.Logger
	loggingAdapter := &LoggerAdapter{logger: log}
	loggingMW := middleware.NewLoggingMiddleware(loggingAdapter)
	router.Use(loggingMW.Handler())

	// 3. 租户识别中间件
	router.Use(middleware.TenantMiddleware())

	// 4. 安全脱敏中间件
	securityMW := middleware.NewSecurityMiddleware(cfg.Security.SensitiveFields)
	router.Use(securityMW.Handler())

	// 公开路由（不需要 API Key）
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// 需要 API Key 的路由组
	apiGroup := router.Group("/api/v1")
	authMW := middleware.NewAuthMiddleware(cfg.Security.APIKeys)
	apiGroup.Use(authMW.Handler())
	{
		apiGroup.POST("/vectors/items", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "vector added"})
		})

		apiGroup.DELETE("/vectors/items", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "vector deleted"})
		})
	}

	// 可选 API Key 的路由组
	chatGroup := router.Group("/chat")
	optionalAuthMW := middleware.NewOptionalAuthMiddleware(cfg.Security.APIKeys)
	chatGroup.Use(optionalAuthMW.Handler())
	{
		chatGroup.POST("", func(c *gin.Context) {
			c.JSON(200, gin.H{"answer": "Hello!"})
		})
	}

	// 启动服务器
	router.Run(":8080")
}

// LoggerAdapter 适配器，将 logger.Logger 转换为 middleware.Logger
type LoggerAdapter struct {
	logger logger.Logger
}

func (la *LoggerAdapter) Info(msg string, fields map[string]interface{}) {
	// 使用空 context，因为中间件层面没有 context
	la.logger.Info(nil, msg, fields)
}

func (la *LoggerAdapter) Error(msg string, fields map[string]interface{}) {
	la.logger.Error(nil, msg, fields)
}
