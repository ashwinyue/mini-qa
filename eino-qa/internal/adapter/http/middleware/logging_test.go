package middleware

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// MockLogger 用于测试的模拟 Logger
type MockLogger struct {
	InfoCalls  []LogCall
	ErrorCalls []LogCall
}

type LogCall struct {
	Message string
	Fields  map[string]interface{}
}

func (m *MockLogger) Info(msg string, fields map[string]interface{}) {
	m.InfoCalls = append(m.InfoCalls, LogCall{
		Message: msg,
		Fields:  fields,
	})
}

func (m *MockLogger) Error(msg string, fields map[string]interface{}) {
	m.ErrorCalls = append(m.ErrorCalls, LogCall{
		Message: msg,
		Fields:  fields,
	})
}

func TestLoggingMiddleware_BasicRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockLogger := &MockLogger{}
	loggingMW := NewLoggingMiddleware(mockLogger)

	// 创建测试路由
	router := gin.New()
	router.Use(TenantMiddleware()) // 先设置租户
	router.Use(loggingMW.Handler())

	router.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "success"})
	})

	// 创建测试请求
	req := httptest.NewRequest("GET", "/test?param=value", nil)
	req.Header.Set("X-Tenant-ID", "tenant1")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// 验证响应
	assert.Equal(t, http.StatusOK, w.Code)

	// 验证日志记录
	assert.GreaterOrEqual(t, len(mockLogger.InfoCalls), 2) // 至少有 started 和 completed

	// 验证请求开始日志
	startLog := mockLogger.InfoCalls[0]
	assert.Equal(t, "request started", startLog.Message)
	assert.Equal(t, "tenant1", startLog.Fields["tenant_id"])
	assert.Equal(t, "GET", startLog.Fields["method"])
	assert.Equal(t, "/test", startLog.Fields["path"])
	assert.NotNil(t, startLog.Fields["request_id"])

	// 验证请求完成日志
	completedLog := mockLogger.InfoCalls[1]
	assert.Equal(t, "request completed successfully", completedLog.Message)
	assert.Equal(t, 200, completedLog.Fields["status_code"])
	assert.NotNil(t, completedLog.Fields["duration_ms"])
}

func TestLoggingMiddleware_RequestID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockLogger := &MockLogger{}
	loggingMW := NewLoggingMiddleware(mockLogger)

	// 创建测试路由
	router := gin.New()
	router.Use(loggingMW.Handler())

	var capturedRequestID string
	router.GET("/test", func(c *gin.Context) {
		requestID, exists := c.Get("request_id")
		assert.True(t, exists)
		capturedRequestID = requestID.(string)
		c.JSON(200, gin.H{"request_id": requestID})
	})

	// 创建测试请求
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// 验证 request_id 被设置
	assert.NotEmpty(t, capturedRequestID)

	// 验证日志中包含相同的 request_id
	assert.GreaterOrEqual(t, len(mockLogger.InfoCalls), 1)
	assert.Equal(t, capturedRequestID, mockLogger.InfoCalls[0].Fields["request_id"])
}

func TestLoggingMiddleware_WithRequestBody(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockLogger := &MockLogger{}
	loggingMW := NewLoggingMiddleware(mockLogger)

	// 创建测试路由
	router := gin.New()
	router.Use(loggingMW.Handler())

	router.POST("/test", func(c *gin.Context) {
		// 验证请求体仍然可以被读取
		var body map[string]interface{}
		err := c.ShouldBindJSON(&body)
		assert.NoError(t, err)
		assert.Equal(t, "test", body["key"])

		c.JSON(200, gin.H{"message": "success"})
	})

	// 创建包含请求体的请求
	bodyBytes := []byte(`{"key": "test"}`)
	req := httptest.NewRequest("POST", "/test", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// 验证响应
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestLoggingMiddleware_WithErrors(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockLogger := &MockLogger{}
	loggingMW := NewLoggingMiddleware(mockLogger)

	// 创建测试路由
	router := gin.New()
	router.Use(loggingMW.Handler())

	router.GET("/test", func(c *gin.Context) {
		c.Error(NewNotFoundError("Resource not found"))
		c.JSON(404, gin.H{"error": "not found"})
	})

	// 创建测试请求
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// 验证错误日志被记录
	assert.GreaterOrEqual(t, len(mockLogger.ErrorCalls), 1)
	errorLog := mockLogger.ErrorCalls[0]
	assert.Equal(t, "request completed with errors", errorLog.Message)
	assert.NotNil(t, errorLog.Fields["errors"])
}

func TestLoggingMiddleware_ServerError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockLogger := &MockLogger{}
	loggingMW := NewLoggingMiddleware(mockLogger)

	// 创建测试路由
	router := gin.New()
	router.Use(loggingMW.Handler())

	router.GET("/test", func(c *gin.Context) {
		c.JSON(500, gin.H{"error": "internal server error"})
	})

	// 创建测试请求
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// 验证服务器错误被记录为 Error 级别
	assert.GreaterOrEqual(t, len(mockLogger.ErrorCalls), 1)
	errorLog := mockLogger.ErrorCalls[0]
	assert.Equal(t, "request completed with server error", errorLog.Message)
	assert.Equal(t, 500, errorLog.Fields["status_code"])
}

func TestLoggingMiddleware_ClientError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockLogger := &MockLogger{}
	loggingMW := NewLoggingMiddleware(mockLogger)

	// 创建测试路由
	router := gin.New()
	router.Use(loggingMW.Handler())

	router.GET("/test", func(c *gin.Context) {
		c.JSON(400, gin.H{"error": "bad request"})
	})

	// 创建测试请求
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// 验证客户端错误被记录为 Info 级别
	assert.GreaterOrEqual(t, len(mockLogger.InfoCalls), 2)
	completedLog := mockLogger.InfoCalls[1]
	assert.Equal(t, "request completed with client error", completedLog.Message)
	assert.Equal(t, 400, completedLog.Fields["status_code"])
}

func TestLoggingMiddleware_WithSanitizedBody(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockLogger := &MockLogger{}
	loggingMW := NewLoggingMiddleware(mockLogger)

	// 创建测试路由
	router := gin.New()
	router.Use(loggingMW.Handler())

	router.POST("/test", func(c *gin.Context) {
		// 模拟 SecurityMiddleware 设置的脱敏内容
		c.Set("sanitized_request_body", `{"password": "[REDACTED]"}`)
		c.JSON(200, gin.H{"message": "success"})
	})

	// 创建测试请求
	bodyBytes := []byte(`{"password": "secret123"}`)
	req := httptest.NewRequest("POST", "/test", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// 验证脱敏后的内容被记录
	found := false
	for _, call := range mockLogger.InfoCalls {
		if call.Message == "request details" {
			found = true
			assert.Equal(t, `{"password": "[REDACTED]"}`, call.Fields["sanitized_body"])
		}
	}
	assert.True(t, found, "应该记录脱敏后的请求详情")
}

func TestLoggingMiddleware_DurationTracking(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockLogger := &MockLogger{}
	loggingMW := NewLoggingMiddleware(mockLogger)

	// 创建测试路由
	router := gin.New()
	router.Use(loggingMW.Handler())

	router.GET("/test", func(c *gin.Context) {
		// 模拟一些处理时间
		c.JSON(200, gin.H{"message": "success"})
	})

	// 创建测试请求
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// 验证处理时长被记录
	assert.GreaterOrEqual(t, len(mockLogger.InfoCalls), 2)
	completedLog := mockLogger.InfoCalls[1]
	durationMs, exists := completedLog.Fields["duration_ms"]
	assert.True(t, exists)
	assert.GreaterOrEqual(t, durationMs.(int64), int64(0))
}

func TestBodyLogWriter_Write(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// 创建一个简单的测试来验证 bodyLogWriter 的功能
	router := gin.New()
	mockLogger := &MockLogger{}
	loggingMW := NewLoggingMiddleware(mockLogger)
	router.Use(loggingMW.Handler())

	router.GET("/test", func(c *gin.Context) {
		c.String(200, "test response")
	})

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// 验证响应被正确写入
	assert.Equal(t, "test response", w.Body.String())

	// 验证日志中记录了响应大小
	assert.GreaterOrEqual(t, len(mockLogger.InfoCalls), 2)
	completedLog := mockLogger.InfoCalls[1]
	responseSize, exists := completedLog.Fields["response_size"]
	assert.True(t, exists)
	assert.Greater(t, responseSize.(int), 0)
}
