package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware API Key 验证中间件
// 验证请求头中的 API Key 是否有效
type AuthMiddleware struct {
	apiKeys map[string]bool
}

// NewAuthMiddleware 创建 API Key 验证中间件
func NewAuthMiddleware(apiKeys []string) *AuthMiddleware {
	// 将 API Key 列表转换为 map 以提高查找效率
	keyMap := make(map[string]bool)
	for _, key := range apiKeys {
		if key != "" {
			keyMap[key] = true
		}
	}

	return &AuthMiddleware{
		apiKeys: keyMap,
	}
}

// Handler 返回 Gin 中间件处理函数
func (am *AuthMiddleware) Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 如果没有配置 API Key，跳过验证
		if len(am.apiKeys) == 0 {
			c.Next()
			return
		}

		// 从请求头获取 API Key
		apiKey := c.GetHeader("X-API-Key")

		// 如果请求头中没有，尝试从 Authorization 头获取
		if apiKey == "" {
			auth := c.GetHeader("Authorization")
			if strings.HasPrefix(auth, "Bearer ") {
				apiKey = strings.TrimPrefix(auth, "Bearer ")
			}
		}

		// 验证 API Key
		if apiKey == "" {
			c.JSON(http.StatusUnauthorized, ErrorResponse{
				Code:    http.StatusUnauthorized,
				Message: "API Key is required",
				TraceID: getTraceID(c),
			})
			c.Abort()
			return
		}

		if !am.apiKeys[apiKey] {
			c.JSON(http.StatusUnauthorized, ErrorResponse{
				Code:    http.StatusUnauthorized,
				Message: "Invalid API Key",
				TraceID: getTraceID(c),
			})
			c.Abort()
			return
		}

		// 验证通过，继续处理
		c.Next()
	}
}

// OptionalAuthMiddleware 可选的 API Key 验证中间件
// 如果提供了 API Key 则验证，否则允许通过
type OptionalAuthMiddleware struct {
	apiKeys map[string]bool
}

// NewOptionalAuthMiddleware 创建可选的 API Key 验证中间件
func NewOptionalAuthMiddleware(apiKeys []string) *OptionalAuthMiddleware {
	keyMap := make(map[string]bool)
	for _, key := range apiKeys {
		if key != "" {
			keyMap[key] = true
		}
	}

	return &OptionalAuthMiddleware{
		apiKeys: keyMap,
	}
}

// Handler 返回 Gin 中间件处理函数
func (oam *OptionalAuthMiddleware) Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 如果没有配置 API Key，跳过验证
		if len(oam.apiKeys) == 0 {
			c.Next()
			return
		}

		// 从请求头获取 API Key
		apiKey := c.GetHeader("X-API-Key")

		// 如果请求头中没有，尝试从 Authorization 头获取
		if apiKey == "" {
			auth := c.GetHeader("Authorization")
			if strings.HasPrefix(auth, "Bearer ") {
				apiKey = strings.TrimPrefix(auth, "Bearer ")
			}
		}

		// 如果没有提供 API Key，允许通过
		if apiKey == "" {
			c.Next()
			return
		}

		// 如果提供了 API Key，验证其有效性
		if !oam.apiKeys[apiKey] {
			c.JSON(http.StatusUnauthorized, ErrorResponse{
				Code:    http.StatusUnauthorized,
				Message: "Invalid API Key",
				TraceID: getTraceID(c),
			})
			c.Abort()
			return
		}

		// 验证通过，继续处理
		c.Next()
	}
}
