package entity

import "time"

// Session 表示对话会话
type Session struct {
	ID        string
	TenantID  string
	Messages  []*Message
	CreatedAt time.Time
	UpdatedAt time.Time
	ExpiresAt time.Time
	Metadata  map[string]any
}

// NewSession 创建新的会话实例
func NewSession(tenantID string, ttl time.Duration) *Session {
	now := time.Now()
	return &Session{
		ID:        generateSessionID(),
		TenantID:  tenantID,
		Messages:  make([]*Message, 0),
		CreatedAt: now,
		UpdatedAt: now,
		ExpiresAt: now.Add(ttl),
		Metadata:  make(map[string]any),
	}
}

// Validate 验证会话的有效性
func (s *Session) Validate() error {
	if s.ID == "" {
		return ErrEmptySessionID
	}

	if s.TenantID == "" {
		return ErrEmptyTenantID
	}

	return nil
}

// IsExpired 判断会话是否过期
func (s *Session) IsExpired() bool {
	return time.Now().After(s.ExpiresAt)
}

// AddMessage 添加消息到会话
func (s *Session) AddMessage(message *Message) error {
	if s.IsExpired() {
		return ErrSessionExpired
	}

	if err := message.Validate(); err != nil {
		return err
	}

	s.Messages = append(s.Messages, message)
	s.UpdatedAt = time.Now()
	return nil
}

// GetMessages 获取所有消息
func (s *Session) GetMessages() []*Message {
	return s.Messages
}

// GetMessageCount 获取消息数量
func (s *Session) GetMessageCount() int {
	return len(s.Messages)
}

// GetLastMessage 获取最后一条消息
func (s *Session) GetLastMessage() *Message {
	if len(s.Messages) == 0 {
		return nil
	}
	return s.Messages[len(s.Messages)-1]
}

// GetUserMessages 获取所有用户消息
func (s *Session) GetUserMessages() []*Message {
	userMessages := make([]*Message, 0)
	for _, msg := range s.Messages {
		if msg.IsUser() {
			userMessages = append(userMessages, msg)
		}
	}
	return userMessages
}

// GetAssistantMessages 获取所有助手消息
func (s *Session) GetAssistantMessages() []*Message {
	assistantMessages := make([]*Message, 0)
	for _, msg := range s.Messages {
		if msg.IsAssistant() {
			assistantMessages = append(assistantMessages, msg)
		}
	}
	return assistantMessages
}

// ExtendExpiration 延长会话过期时间
func (s *Session) ExtendExpiration(duration time.Duration) {
	s.ExpiresAt = time.Now().Add(duration)
	s.UpdatedAt = time.Now()
}

// AddMetadata 添加元数据
func (s *Session) AddMetadata(key string, value any) {
	if s.Metadata == nil {
		s.Metadata = make(map[string]any)
	}
	s.Metadata[key] = value
}

// generateSessionID 生成会话 ID
func generateSessionID() string {
	return generateUniqueID("sess_", 32)
}
