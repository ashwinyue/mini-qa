package entity

import "time"

// Message 表示对话中的一条消息
type Message struct {
	ID        string
	Content   string
	Role      string // "user", "assistant", "system"
	Timestamp time.Time
	Metadata  map[string]any
}

// NewMessage 创建新的消息实例
func NewMessage(content, role string) *Message {
	return &Message{
		ID:        generateMessageID(),
		Content:   content,
		Role:      role,
		Timestamp: time.Now(),
		Metadata:  make(map[string]any),
	}
}

// Validate 验证消息的有效性
func (m *Message) Validate() error {
	if m.Content == "" {
		return ErrEmptyContent
	}

	validRoles := map[string]bool{
		"user":      true,
		"assistant": true,
		"system":    true,
	}

	if !validRoles[m.Role] {
		return ErrInvalidRole
	}

	return nil
}

// IsUser 判断是否为用户消息
func (m *Message) IsUser() bool {
	return m.Role == "user"
}

// IsAssistant 判断是否为助手消息
func (m *Message) IsAssistant() bool {
	return m.Role == "assistant"
}

// IsSystem 判断是否为系统消息
func (m *Message) IsSystem() bool {
	return m.Role == "system"
}

// generateMessageID 生成消息 ID
func generateMessageID() string {
	// 使用时间戳和随机数生成唯一 ID
	return time.Now().Format("20060102150405") + randomString(8)
}
