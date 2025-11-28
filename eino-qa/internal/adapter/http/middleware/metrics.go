package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
)

// MetricsCollector 指标收集器接口
type MetricsCollector interface {
	RecordRequest(route string, statusCode int, duration time.Duration)
	RecordError(route string, errorType string)
}

// MetricsMiddleware 指标收集中间件
// 需求: 7.4 - 统计各类请求的数量和平均响应时间
type MetricsMiddleware struct {
	collector MetricsCollector
}

// NewMetricsMiddleware 创建指标收集中间件
func NewMetricsMiddleware(collector MetricsCollector) *MetricsMiddleware {
	return &MetricsMiddleware{
		collector: collector,
	}
}

// Handler 返回 Gin 中间件处理函数
func (mm *MetricsMiddleware) Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 记录开始时间
		startTime := time.Now()

		// 获取路由路径（用于分类统计）
		// 使用匹配的路由模式，如果没有则使用实际路径
		route := c.Request.URL.Path
		if matched := c.GetString("matched_route"); matched != "" {
			route = matched
		}

		// 处理请求
		c.Next()

		// 计算处理时长
		duration := time.Since(startTime)

		// 获取响应状态码
		statusCode := c.Writer.Status()

		// 记录请求指标
		mm.collector.RecordRequest(route, statusCode, duration)

		// 如果有错误，记录错误指标
		if len(c.Errors) > 0 {
			for _, err := range c.Errors {
				errorType := "unknown"
				if err.Type == gin.ErrorTypeBind {
					errorType = "bind_error"
				} else if err.Type == gin.ErrorTypeRender {
					errorType = "render_error"
				} else if err.Type == gin.ErrorTypePrivate {
					errorType = "private_error"
				} else if err.Type == gin.ErrorTypePublic {
					errorType = "public_error"
				} else if err.Type == gin.ErrorTypeAny {
					errorType = "any_error"
				}
				mm.collector.RecordError(route, errorType)
			}
		}
	}
}
