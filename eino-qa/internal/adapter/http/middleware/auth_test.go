package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestAuthMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	apiKeys := []string{"valid-key-1", "valid-key-2"}

	tests := []struct {
		name           string
		headerName     string
		headerValue    string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "有效的 API Key (X-API-Key)",
			headerName:     "X-API-Key",
			headerValue:    "valid-key-1",
			expectedStatus: http.StatusOK,
			expectedBody:   "",
		},
		{
			name:           "有效的 API Key (Authorization Bearer)",
			headerName:     "Authorization",
			headerValue:    "Bearer valid-key-2",
			expectedStatus: http.StatusOK,
			expectedBody:   "",
		},
		{
			name:           "缺少 API Key",
			headerName:     "",
			headerValue:    "",
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   "API Key is required",
		},
		{
			name:           "无效的 API Key",
			headerName:     "X-API-Key",
			headerValue:    "invalid-key",
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   "Invalid API Key",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建测试路由
			router := gin.New()
			authMW := NewAuthMiddleware(apiKeys)
			router.Use(authMW.Handler())

			// 添加测试处理器
			router.GET("/test", func(c *gin.Context) {
				c.JSON(200, gin.H{"message": "success"})
			})

			// 创建测试请求
			req := httptest.NewRequest("GET", "/test", nil)
			if tt.headerName != "" {
				req.Header.Set(tt.headerName, tt.headerValue)
			}

			// 执行请求
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// 验证响应
			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.expectedBody != "" {
				assert.Contains(t, w.Body.String(), tt.expectedBody)
			}
		})
	}
}

func TestAuthMiddleware_NoKeysConfigured(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// 创建测试路由（没有配置 API Key）
	router := gin.New()
	authMW := NewAuthMiddleware([]string{})
	router.Use(authMW.Handler())

	router.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "success"})
	})

	// 创建测试请求（不带 API Key）
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// 验证：没有配置 API Key 时应该允许通过
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestOptionalAuthMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	apiKeys := []string{"valid-key"}

	tests := []struct {
		name           string
		headerValue    string
		expectedStatus int
	}{
		{
			name:           "有效的 API Key",
			headerValue:    "valid-key",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "没有提供 API Key（允许通过）",
			headerValue:    "",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "无效的 API Key",
			headerValue:    "invalid-key",
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建测试路由
			router := gin.New()
			optionalAuthMW := NewOptionalAuthMiddleware(apiKeys)
			router.Use(optionalAuthMW.Handler())

			router.GET("/test", func(c *gin.Context) {
				c.JSON(200, gin.H{"message": "success"})
			})

			// 创建测试请求
			req := httptest.NewRequest("GET", "/test", nil)
			if tt.headerValue != "" {
				req.Header.Set("X-API-Key", tt.headerValue)
			}

			// 执行请求
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// 验证响应
			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}
