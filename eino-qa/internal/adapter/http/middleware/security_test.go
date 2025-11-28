package middleware

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestSecurityMiddleware_SanitizeString(t *testing.T) {
	securityMW := NewSecurityMiddleware([]string{"password", "secret"})

	tests := []struct {
		name        string
		input       string
		contains    []string
		notContains []string
	}{
		{
			name:        "脱敏密码字段",
			input:       `{"password": "secret123"}`,
			contains:    []string{"[REDACTED]"},
			notContains: []string{"secret123"},
		},
		{
			name:        "脱敏身份证号",
			input:       "我的身份证号是 110101199001011234",
			contains:    []string{"[REDACTED]"},
			notContains: []string{"110101199001011234"},
		},
		{
			name:        "脱敏手机号",
			input:       "联系电话：13812345678",
			contains:    []string{"[REDACTED]"},
			notContains: []string{"13812345678"},
		},
		{
			name:        "脱敏邮箱",
			input:       "邮箱：user@example.com",
			contains:    []string{"[REDACTED]"},
			notContains: []string{"user@example.com"},
		},
		{
			name:     "不脱敏普通文本",
			input:    "这是一段普通的文本",
			contains: []string{"这是一段普通的文本"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := securityMW.sanitizeString(tt.input)

			for _, s := range tt.contains {
				assert.Contains(t, result, s)
			}

			for _, s := range tt.notContains {
				assert.NotContains(t, result, s)
			}
		})
	}
}

func TestSecurityMiddleware_Handler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	securityMW := NewSecurityMiddleware([]string{"password"})

	// 创建测试路由
	router := gin.New()
	router.Use(securityMW.Handler())

	router.POST("/test", func(c *gin.Context) {
		// 检查是否设置了脱敏后的请求体
		sanitized, exists := c.Get("sanitized_request_body")
		assert.True(t, exists)

		// 验证脱敏后的内容不包含敏感信息
		sanitizedStr := sanitized.(string)
		assert.Contains(t, sanitizedStr, "[REDACTED]")
		assert.NotContains(t, sanitizedStr, "secret123")

		c.JSON(200, gin.H{"message": "success"})
	})

	// 创建包含敏感信息的请求
	body := []byte(`{"password": "secret123", "username": "user1"}`)
	req := httptest.NewRequest("POST", "/test", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	// 执行请求
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// 验证响应
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestSecurityMiddleware_IsSensitiveField(t *testing.T) {
	securityMW := NewSecurityMiddleware([]string{"custom_secret"})

	tests := []struct {
		name      string
		fieldName string
		expected  bool
	}{
		{"密码字段", "password", true},
		{"密码字段（大写）", "Password", true},
		{"密码字段（下划线）", "user_password", true},
		{"API Key", "api_key", true},
		{"Token", "access_token", true},
		{"自定义敏感字段", "custom_secret", true},
		{"普通字段", "username", false},
		{"普通字段", "email_address", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := securityMW.isSensitiveField(tt.fieldName)
			assert.Equal(t, tt.expected, result)
		})
	}
}
