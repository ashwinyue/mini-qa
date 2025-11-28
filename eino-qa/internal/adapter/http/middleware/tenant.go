package middleware

import (
	"context"

	"github.com/gin-gonic/gin"
)

// TenantMiddleware 租户识别中间件
// 从请求头或查询参数中提取租户 ID，并设置到 context 中
func TenantMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 优先从请求头获取租户 ID
		tenantID := c.GetHeader("X-Tenant-ID")

		// 如果请求头中没有，尝试从查询参数获取
		if tenantID == "" {
			tenantID = c.Query("tenant")
		}

		// 如果都没有，使用默认租户
		if tenantID == "" {
			tenantID = "default"
		}

		// 将租户 ID 设置到 Gin context 中
		c.Set("tenant_id", tenantID)

		// 同时设置到 request context 中，供后续使用
		ctx := context.WithValue(c.Request.Context(), "tenant_id", tenantID)
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}
