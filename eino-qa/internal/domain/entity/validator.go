package entity

import (
	"errors"
	"regexp"
	"strings"
)

var (
	// 订单 ID 格式: #YYYYMMDDXXX
	orderIDPattern = regexp.MustCompile(`^#\d{8}\d{3}$`)

	// 会话 ID 格式: sess_YYYYMMDDHHMMSSXXXXXXXXXXXXXXXX
	sessionIDPattern = regexp.MustCompile(`^sess_\d{14}[a-z0-9]{16}$`)
)

// ValidateOrderID 验证订单 ID 格式
func ValidateOrderID(orderID string) error {
	if orderID == "" {
		return errors.New("order ID cannot be empty")
	}

	if !orderIDPattern.MatchString(orderID) {
		return errors.New("invalid order ID format, expected #YYYYMMDDXXX")
	}

	return nil
}

// ValidateSessionID 验证会话 ID 格式
func ValidateSessionID(sessionID string) error {
	if sessionID == "" {
		return errors.New("session ID cannot be empty")
	}

	if !sessionIDPattern.MatchString(sessionID) {
		return errors.New("invalid session ID format")
	}

	return nil
}

// ValidateTenantID 验证租户 ID
func ValidateTenantID(tenantID string) error {
	if tenantID == "" {
		return errors.New("tenant ID cannot be empty")
	}

	// 租户 ID 只能包含字母、数字、下划线和连字符
	if !regexp.MustCompile(`^[a-zA-Z0-9_-]+$`).MatchString(tenantID) {
		return errors.New("tenant ID can only contain letters, numbers, underscores and hyphens")
	}

	return nil
}

// ValidateConfidence 验证置信度分数
func ValidateConfidence(confidence float64) error {
	if confidence < 0 || confidence > 1 {
		return errors.New("confidence must be between 0 and 1")
	}
	return nil
}

// ExtractOrderID 从文本中提取订单 ID
func ExtractOrderID(text string) (string, error) {
	matches := orderIDPattern.FindStringSubmatch(text)
	if len(matches) == 0 {
		return "", errors.New("no order ID found in text")
	}
	return matches[0], nil
}

// SanitizeSQL 清理 SQL 查询，防止注入
func SanitizeSQL(sql string) error {
	// 转换为大写进行检查
	upperSQL := strings.ToUpper(sql)

	// 检查危险操作
	dangerousKeywords := []string{
		"DROP", "DELETE", "UPDATE", "INSERT",
		"ALTER", "CREATE", "TRUNCATE", "EXEC",
		"EXECUTE", "UNION", "--", ";",
	}

	for _, keyword := range dangerousKeywords {
		if strings.Contains(upperSQL, keyword) {
			return errors.New("SQL contains dangerous keyword: " + keyword)
		}
	}

	return nil
}

// ValidateVector 验证向量维度
func ValidateVector(vector []float32, expectedDim int) error {
	if len(vector) == 0 {
		return errors.New("vector cannot be empty")
	}

	if len(vector) != expectedDim {
		return errors.New("vector dimension mismatch")
	}

	return nil
}
