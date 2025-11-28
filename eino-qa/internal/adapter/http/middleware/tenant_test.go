package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestTenantMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		headerValue    string
		queryValue     string
		expectedTenant string
	}{
		{
			name:           "从请求头获取租户 ID",
			headerValue:    "tenant1",
			queryValue:     "",
			expectedTenant: "tenant1",
		},
		{
			name:           "从查询参数获取租户 ID",
			headerValue:    "",
			queryValue:     "tenant2",
			expectedTenant: "tenant2",
		},
		{
			name:           "请求头优先于查询参数",
			headerValue:    "tenant1",
			queryValue:     "tenant2",
			expectedTenant: "tenant1",
		},
		{
			name:           "使用默认租户",
			headerValue:    "",
			queryValue:     "",
			expectedTenant: "default",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建测试路由
			router := gin.New()
			router.Use(TenantMiddleware())

			// 添加测试处理器
			router.GET("/test", func(c *gin.Context) {
				tenantID, exists := c.Get("tenant_id")
				assert.True(t, exists)
				assert.Equal(t, tt.expectedTenant, tenantID)

				// 验证 context 中也设置了租户 ID
				ctxTenantID := c.Request.Context().Value("tenant_id")
				assert.Equal(t, tt.expectedTenant, ctxTenantID)

				c.JSON(200, gin.H{"tenant_id": tenantID})
			})

			// 创建测试请求
			req := httptest.NewRequest("GET", "/test", nil)
			if tt.headerValue != "" {
				req.Header.Set("X-Tenant-ID", tt.headerValue)
			}
			if tt.queryValue != "" {
				q := req.URL.Query()
				q.Add("tenant", tt.queryValue)
				req.URL.RawQuery = q.Encode()
			}

			// 执行请求
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// 验证响应
			assert.Equal(t, http.StatusOK, w.Code)
		})
	}
}
