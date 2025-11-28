package middleware

import (
	"bytes"
	"io"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// LoggingMiddleware 日志记录中间件
// 记录请求的详细信息，包括请求 ID、租户 ID、查询内容和处理时长
type LoggingMiddleware struct {
	logger Logger
}

// Logger 日志接口（简化版，适配实际的 logger）
type Logger interface {
	Info(msg string, fields map[string]interface{})
	Error(msg string, fields map[string]interface{})
}

// NewLoggingMiddleware 创建日志记录中间件
func NewLoggingMiddleware(logger Logger) *LoggingMiddleware {
	return &LoggingMiddleware{
		logger: logger,
	}
}

// Handler 返回 Gin 中间件处理函数
func (lm *LoggingMiddleware) Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 生成请求 ID
		requestID := uuid.New().String()
		c.Set("request_id", requestID)

		// 记录开始时间
		startTime := time.Now()

		// 获取租户 ID（由 TenantMiddleware 设置）
		tenantID, _ := c.Get("tenant_id")

		// 读取请求体（用于日志记录）
		if c.Request.Body != nil {
			bodyBytes, err := io.ReadAll(c.Request.Body)
			if err == nil {
				// 恢复请求体供后续使用
				c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
			}
		}

		// 使用 ResponseWriter 包装器捕获响应
		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw

		// 记录请求开始
		lm.logger.Info("request started", map[string]interface{}{
			"request_id": requestID,
			"tenant_id":  tenantID,
			"method":     c.Request.Method,
			"path":       c.Request.URL.Path,
			"query":      c.Request.URL.RawQuery,
			"client_ip":  c.ClientIP(),
			"user_agent": c.Request.UserAgent(),
		})

		// 处理请求
		c.Next()

		// 计算处理时长
		duration := time.Since(startTime)

		// 获取响应状态码
		statusCode := c.Writer.Status()

		// 构建日志字段
		logFields := map[string]interface{}{
			"request_id":    requestID,
			"tenant_id":     tenantID,
			"method":        c.Request.Method,
			"path":          c.Request.URL.Path,
			"status_code":   statusCode,
			"duration_ms":   duration.Milliseconds(),
			"client_ip":     c.ClientIP(),
			"request_size":  c.Request.ContentLength,
			"response_size": blw.body.Len(),
		}

		// 如果有错误，记录错误信息
		if len(c.Errors) > 0 {
			logFields["errors"] = c.Errors.String()
			lm.logger.Error("request completed with errors", logFields)
		} else {
			// 根据状态码决定日志级别
			if statusCode >= 500 {
				lm.logger.Error("request completed with server error", logFields)
			} else if statusCode >= 400 {
				lm.logger.Info("request completed with client error", logFields)
			} else {
				lm.logger.Info("request completed successfully", logFields)
			}
		}

		// 如果请求体包含查询内容，记录查询详情（从脱敏后的内容获取）
		if sanitizedBody, exists := c.Get("sanitized_request_body"); exists {
			lm.logger.Info("request details", map[string]interface{}{
				"request_id":     requestID,
				"tenant_id":      tenantID,
				"sanitized_body": sanitizedBody,
			})
		}
	}
}

// bodyLogWriter 用于捕获响应体的 ResponseWriter 包装器
type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

// Write 重写 Write 方法以捕获响应体
func (w *bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// WriteString 重写 WriteString 方法
func (w *bodyLogWriter) WriteString(s string) (int, error) {
	w.body.WriteString(s)
	return w.ResponseWriter.WriteString(s)
}
