package middleware

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestErrorHandler_CustomErrors(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		err            error
		expectedStatus int
		expectedMsg    string
	}{
		{
			name:           "ValidationError",
			err:            NewValidationError("Validation failed", map[string]string{"field": "error"}),
			expectedStatus: http.StatusBadRequest,
			expectedMsg:    "Validation failed",
		},
		{
			name:           "NotFoundError",
			err:            NewNotFoundError("Resource not found"),
			expectedStatus: http.StatusNotFound,
			expectedMsg:    "Resource not found",
		},
		{
			name:           "UnauthorizedError",
			err:            NewUnauthorizedError("Unauthorized access"),
			expectedStatus: http.StatusUnauthorized,
			expectedMsg:    "Unauthorized access",
		},
		{
			name:           "ForbiddenError",
			err:            NewForbiddenError("Access forbidden"),
			expectedStatus: http.StatusForbidden,
			expectedMsg:    "Access forbidden",
		},
		{
			name:           "BadRequestError",
			err:            NewBadRequestError("Bad request"),
			expectedStatus: http.StatusBadRequest,
			expectedMsg:    "Bad request",
		},
		{
			name:           "ServiceError",
			err:            NewServiceError("Service unavailable", "milvus"),
			expectedStatus: http.StatusBadGateway,
			expectedMsg:    "Service unavailable",
		},
		{
			name:           "Generic Error",
			err:            errors.New("generic error"),
			expectedStatus: http.StatusInternalServerError,
			expectedMsg:    "generic error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建测试路由
			router := gin.New()
			router.Use(ErrorHandler())

			router.GET("/test", func(c *gin.Context) {
				c.Error(tt.err)
			})

			// 创建测试请求
			req := httptest.NewRequest("GET", "/test", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// 验证响应
			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.Contains(t, w.Body.String(), tt.expectedMsg)
		})
	}
}

func TestErrorHandler_Panic(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// 创建测试路由
	router := gin.New()
	router.Use(ErrorHandler())

	router.GET("/test", func(c *gin.Context) {
		panic("something went wrong")
	})

	// 创建测试请求
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// 验证响应：应该捕获 panic 并返回 500
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "Internal server error")
}

func TestErrorHandler_NoError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// 创建测试路由
	router := gin.New()
	router.Use(ErrorHandler())

	router.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "success"})
	})

	// 创建测试请求
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// 验证响应：没有错误时应该正常返回
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "success")
}

func TestErrorResponse_Structure(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// 创建测试路由
	router := gin.New()
	router.Use(ErrorHandler())

	router.GET("/test", func(c *gin.Context) {
		c.Set("request_id", "test-request-id")
		c.Error(NewNotFoundError("Resource not found"))
	})

	// 创建测试请求
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// 验证响应结构
	assert.Equal(t, http.StatusNotFound, w.Code)
	body := w.Body.String()
	assert.Contains(t, body, "code")
	assert.Contains(t, body, "message")
	assert.Contains(t, body, "trace_id")
	assert.Contains(t, body, "test-request-id")
}
