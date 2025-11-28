package http

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"eino-qa/internal/adapter/http/handler"
	"eino-qa/internal/adapter/http/middleware"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func init() {
	gin.SetMode(gin.TestMode)
}

// TestSetupRouter 测试路由设置
func TestSetupRouter(t *testing.T) {
	config := DefaultRouterConfig()
	router := SetupRouter(config)

	assert.NotNil(t, router)
}

// TestRootRoute 测试根路由
func TestRootRoute(t *testing.T) {
	config := DefaultRouterConfig()
	router := SetupRouter(config)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "eino-qa-system", response["service"])
	assert.Equal(t, "running", response["status"])
}

// TestHealthRoutes 测试健康检查路由
func TestHealthRoutes(t *testing.T) {
	healthHandler := handler.NewHealthHandler()
	config := DefaultRouterConfig()
	config.HealthHandler = healthHandler
	router := SetupRouter(config)

	tests := []struct {
		name       string
		path       string
		wantStatus int
	}{
		{
			name:       "health check",
			path:       "/health",
			wantStatus: http.StatusOK,
		},
		{
			name:       "liveness check",
			path:       "/health/live",
			wantStatus: http.StatusOK,
		},
		{
			name:       "readiness check",
			path:       "/health/ready",
			wantStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", tt.path, nil)
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
}

// TestTenantMiddleware 测试租户中间件
func TestTenantMiddleware(t *testing.T) {
	config := DefaultRouterConfig()
	router := SetupRouter(config)

	// 添加测试路由
	router.GET("/test", func(c *gin.Context) {
		tenantID, exists := c.Get("tenant_id")
		assert.True(t, exists)
		c.JSON(http.StatusOK, gin.H{"tenant_id": tenantID})
	})

	tests := []struct {
		name           string
		header         string
		query          string
		expectedTenant string
	}{
		{
			name:           "tenant from header",
			header:         "tenant1",
			query:          "",
			expectedTenant: "tenant1",
		},
		{
			name:           "tenant from query",
			header:         "",
			query:          "tenant2",
			expectedTenant: "tenant2",
		},
		{
			name:           "default tenant",
			header:         "",
			query:          "",
			expectedTenant: "default",
		},
		{
			name:           "header takes precedence",
			header:         "tenant1",
			query:          "tenant2",
			expectedTenant: "tenant1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/test", nil)

			if tt.header != "" {
				req.Header.Set("X-Tenant-ID", tt.header)
			}
			if tt.query != "" {
				q := req.URL.Query()
				q.Add("tenant", tt.query)
				req.URL.RawQuery = q.Encode()
			}

			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedTenant, response["tenant_id"])
		})
	}
}

// TestAuthMiddleware 测试认证中间件
func TestAuthMiddleware(t *testing.T) {
	apiKey := "test-api-key"
	config := DefaultRouterConfig()
	authMiddleware := middleware.NewAuthMiddleware([]string{apiKey})
	config.AuthMiddleware = authMiddleware.Handler()

	// 添加一个测试 handler
	config.VectorHandler = &handler.VectorHandler{}

	router := SetupRouter(config)

	tests := []struct {
		name       string
		apiKey     string
		wantStatus int
	}{
		{
			name:       "invalid api key",
			apiKey:     "wrong-key",
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "missing api key",
			apiKey:     "",
			wantStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/api/v1/vectors/count", nil)

			if tt.apiKey != "" {
				req.Header.Set("X-API-Key", tt.apiKey)
			}

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
}

// TestErrorMiddleware 测试错误处理中间件
func TestErrorMiddleware(t *testing.T) {
	config := DefaultRouterConfig()
	router := SetupRouter(config)

	// 添加测试路由，触发错误
	router.GET("/error", func(c *gin.Context) {
		c.Error(middleware.NewBadRequestError("test error"))
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/error", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(http.StatusBadRequest), response["code"])
	assert.Contains(t, response["message"], "test error")
}

// TestChatRoute 测试对话路由（不需要认证）
func TestChatRoute(t *testing.T) {
	config := DefaultRouterConfig()
	// 添加一个测试 handler
	config.ChatHandler = &handler.ChatHandler{}
	router := SetupRouter(config)

	// 测试 POST /chat 路由存在
	w := httptest.NewRecorder()
	body := bytes.NewBufferString(`{"query":"test"}`)
	req, _ := http.NewRequest("POST", "/chat", body)
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	// 由于没有实际的 use case，会返回错误，但不应该是 404
	assert.NotEqual(t, http.StatusNotFound, w.Code)
}

// TestVectorRoutes 测试向量管理路由
func TestVectorRoutes(t *testing.T) {
	apiKey := "test-api-key"
	config := DefaultRouterConfig()
	authMiddleware := middleware.NewAuthMiddleware([]string{apiKey})
	config.AuthMiddleware = authMiddleware.Handler()

	// 添加一个测试 handler
	config.VectorHandler = &handler.VectorHandler{}

	router := SetupRouter(config)

	tests := []struct {
		name       string
		method     string
		path       string
		apiKey     string
		wantStatus int
	}{
		{
			name:       "add vectors without auth",
			method:     "POST",
			path:       "/api/v1/vectors/items",
			apiKey:     "",
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "delete vectors without auth",
			method:     "DELETE",
			path:       "/api/v1/vectors/items",
			apiKey:     "",
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "get count without auth",
			method:     "GET",
			path:       "/api/v1/vectors/count",
			apiKey:     "",
			wantStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(tt.method, tt.path, nil)

			if tt.apiKey != "" {
				req.Header.Set("X-API-Key", tt.apiKey)
			}

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
}

// TestModelRoutes 测试模型管理路由
func TestModelRoutes(t *testing.T) {
	apiKey := "test-api-key"
	config := DefaultRouterConfig()
	authMiddleware := middleware.NewAuthMiddleware([]string{apiKey})
	config.AuthMiddleware = authMiddleware.Handler()

	// 添加一个测试 handler
	config.ModelHandler = handler.NewModelHandler("qwen-turbo", "text-embedding-v2")

	router := SetupRouter(config)

	tests := []struct {
		name       string
		method     string
		path       string
		apiKey     string
		wantStatus int
	}{
		{
			name:       "list models without auth",
			method:     "GET",
			path:       "/models",
			apiKey:     "",
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "get current model without auth",
			method:     "GET",
			path:       "/models/current",
			apiKey:     "",
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "switch model without auth",
			method:     "POST",
			path:       "/models/switch",
			apiKey:     "",
			wantStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(tt.method, tt.path, nil)

			if tt.apiKey != "" {
				req.Header.Set("X-API-Key", tt.apiKey)
			}

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
}

// TestMiddlewareOrder 测试中间件顺序
func TestMiddlewareOrder(t *testing.T) {
	var executionOrder []string

	// 创建测试中间件
	testMiddleware1 := func() gin.HandlerFunc {
		return func(c *gin.Context) {
			executionOrder = append(executionOrder, "middleware1")
			c.Next()
		}
	}

	testMiddleware2 := func() gin.HandlerFunc {
		return func(c *gin.Context) {
			executionOrder = append(executionOrder, "middleware2")
			c.Next()
		}
	}

	config := &RouterConfig{
		Mode:               gin.TestMode,
		TenantMiddleware:   testMiddleware1(),
		SecurityMiddleware: testMiddleware2(),
	}

	router := SetupRouter(config)

	// 添加测试路由
	router.GET("/test", func(c *gin.Context) {
		executionOrder = append(executionOrder, "handler")
		c.JSON(http.StatusOK, gin.H{})
	})

	// 重置执行顺序
	executionOrder = []string{}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	// 验证中间件执行顺序
	assert.Equal(t, []string{"middleware1", "middleware2", "handler"}, executionOrder)
}
