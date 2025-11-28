package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
)

// SecurityMiddleware 安全脱敏中间件
// 对请求和响应中的敏感信息进行脱敏处理
type SecurityMiddleware struct {
	sensitiveFields []string
	patterns        map[string]*regexp.Regexp
}

// NewSecurityMiddleware 创建安全脱敏中间件
func NewSecurityMiddleware(sensitiveFields []string) *SecurityMiddleware {
	// 预编译常用的敏感信息正则表达式
	patterns := map[string]*regexp.Regexp{
		"password":    regexp.MustCompile(`(?i)(password|pwd|passwd)[\s]*[:=][\s]*["']?([^"'\s,}]+)["']?`),
		"id_card":     regexp.MustCompile(`\b\d{17}[\dXx]\b`),                                    // 18位身份证号
		"phone":       regexp.MustCompile(`\b1[3-9]\d{9}\b`),                                     // 11位手机号
		"email":       regexp.MustCompile(`\b[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Z|a-z]{2,}\b`), // 邮箱
		"credit_card": regexp.MustCompile(`\b\d{4}[\s-]?\d{4}[\s-]?\d{4}[\s-]?\d{4}\b`),          // 信用卡号
		"api_key":     regexp.MustCompile(`(?i)(api[_-]?key|apikey|access[_-]?token)[\s]*[:=][\s]*["']?([^"'\s,}]+)["']?`),
	}

	return &SecurityMiddleware{
		sensitiveFields: sensitiveFields,
		patterns:        patterns,
	}
}

// Handler 返回 Gin 中间件处理函数
func (sm *SecurityMiddleware) Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 脱敏请求体
		if c.Request.Body != nil && c.Request.ContentLength > 0 {
			bodyBytes, err := io.ReadAll(c.Request.Body)
			if err == nil {
				// 恢复请求体供后续使用
				c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

				// 脱敏处理（仅用于日志记录）
				sanitized := sm.sanitizeString(string(bodyBytes))
				c.Set("sanitized_request_body", sanitized)
			}
		}

		// 继续处理请求
		c.Next()
	}
}

// sanitizeString 对字符串进行脱敏处理
func (sm *SecurityMiddleware) sanitizeString(input string) string {
	result := input

	// 使用正则表达式替换敏感信息
	for _, pattern := range sm.patterns {
		result = pattern.ReplaceAllStringFunc(result, func(match string) string {
			// 保留字段名，替换值为 [REDACTED]
			if strings.Contains(match, ":") || strings.Contains(match, "=") {
				parts := regexp.MustCompile(`[:=]`).Split(match, 2)
				if len(parts) == 2 {
					return parts[0] + ": \"[REDACTED]\""
				}
			}
			return "[REDACTED]"
		})
	}

	// 脱敏 JSON 中的敏感字段
	result = sm.sanitizeJSON(result)

	return result
}

// sanitizeJSON 脱敏 JSON 中的敏感字段
func (sm *SecurityMiddleware) sanitizeJSON(input string) string {
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(input), &data); err != nil {
		// 不是有效的 JSON，返回原字符串
		return input
	}

	// 递归脱敏
	sm.sanitizeMap(data)

	// 转回 JSON
	sanitized, err := json.Marshal(data)
	if err != nil {
		return input
	}

	return string(sanitized)
}

// sanitizeMap 递归脱敏 map
func (sm *SecurityMiddleware) sanitizeMap(data map[string]interface{}) {
	for key, value := range data {
		// 检查是否是敏感字段
		if sm.isSensitiveField(key) {
			data[key] = "[REDACTED]"
			continue
		}

		// 递归处理嵌套结构
		switch v := value.(type) {
		case map[string]interface{}:
			sm.sanitizeMap(v)
		case []interface{}:
			sm.sanitizeSlice(v)
		}
	}
}

// sanitizeSlice 递归脱敏 slice
func (sm *SecurityMiddleware) sanitizeSlice(data []interface{}) {
	for i, item := range data {
		switch v := item.(type) {
		case map[string]interface{}:
			sm.sanitizeMap(v)
		case []interface{}:
			sm.sanitizeSlice(v)
		case string:
			data[i] = sm.sanitizeString(v)
		}
	}
}

// isSensitiveField 判断字段名是否敏感
func (sm *SecurityMiddleware) isSensitiveField(fieldName string) bool {
	lowerField := strings.ToLower(fieldName)

	// 检查配置的敏感字段列表
	for _, sensitive := range sm.sensitiveFields {
		if strings.Contains(lowerField, strings.ToLower(sensitive)) {
			return true
		}
	}

	// 检查常见的敏感字段
	commonSensitive := []string{
		"password", "pwd", "passwd", "secret",
		"token", "api_key", "apikey", "access_key",
		"id_card", "idcard", "identity",
		"phone", "mobile", "telephone",
		"credit_card", "card_number",
		"ssn", "social_security",
	}

	for _, sensitive := range commonSensitive {
		if strings.Contains(lowerField, sensitive) {
			return true
		}
	}

	return false
}

// SanitizeForLog 对日志内容进行脱敏（公共方法）
func (sm *SecurityMiddleware) SanitizeForLog(input string) string {
	return sm.sanitizeString(input)
}
