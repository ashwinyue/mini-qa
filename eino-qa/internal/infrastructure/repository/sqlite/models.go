package sqlite

import (
	"encoding/json"
	"time"

	"eino-qa/internal/domain/entity"
)

// OrderModel GORM 订单模型
type OrderModel struct {
	ID         string    `gorm:"primaryKey;type:varchar(50)"`
	UserID     string    `gorm:"type:varchar(100);index;not null"`
	CourseName string    `gorm:"type:varchar(200);not null"`
	Amount     float64   `gorm:"type:decimal(10,2);not null"`
	Status     string    `gorm:"type:varchar(20);index;not null"`
	TenantID   string    `gorm:"type:varchar(100);index;not null"`
	Metadata   string    `gorm:"type:text"`
	CreatedAt  time.Time `gorm:"autoCreateTime"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime"`
}

// TableName 指定表名
func (OrderModel) TableName() string {
	return "orders"
}

// ToEntity 转换为领域实体
func (m *OrderModel) ToEntity() (*entity.Order, error) {
	order := &entity.Order{
		ID:         m.ID,
		UserID:     m.UserID,
		CourseName: m.CourseName,
		Amount:     m.Amount,
		Status:     entity.OrderStatus(m.Status),
		TenantID:   m.TenantID,
		CreatedAt:  m.CreatedAt,
		UpdatedAt:  m.UpdatedAt,
		Metadata:   make(map[string]any),
	}

	// 解析 Metadata JSON
	if m.Metadata != "" {
		if err := json.Unmarshal([]byte(m.Metadata), &order.Metadata); err != nil {
			return nil, err
		}
	}

	return order, nil
}

// FromEntity 从领域实体创建
func (m *OrderModel) FromEntity(order *entity.Order) error {
	m.ID = order.ID
	m.UserID = order.UserID
	m.CourseName = order.CourseName
	m.Amount = order.Amount
	m.Status = string(order.Status)
	m.TenantID = order.TenantID
	m.CreatedAt = order.CreatedAt
	m.UpdatedAt = order.UpdatedAt

	// 序列化 Metadata
	if order.Metadata != nil && len(order.Metadata) > 0 {
		metadataBytes, err := json.Marshal(order.Metadata)
		if err != nil {
			return err
		}
		m.Metadata = string(metadataBytes)
	}

	return nil
}

// SessionModel GORM 会话模型
type SessionModel struct {
	ID        string    `gorm:"primaryKey;type:varchar(100)"`
	TenantID  string    `gorm:"type:varchar(100);index;not null"`
	Messages  string    `gorm:"type:text"`
	Metadata  string    `gorm:"type:text"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
	ExpiresAt time.Time `gorm:"index;not null"`
}

// TableName 指定表名
func (SessionModel) TableName() string {
	return "sessions"
}

// ToEntity 转换为领域实体
func (m *SessionModel) ToEntity() (*entity.Session, error) {
	session := &entity.Session{
		ID:        m.ID,
		TenantID:  m.TenantID,
		Messages:  make([]*entity.Message, 0),
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
		ExpiresAt: m.ExpiresAt,
		Metadata:  make(map[string]any),
	}

	// 解析 Messages JSON
	if m.Messages != "" {
		if err := json.Unmarshal([]byte(m.Messages), &session.Messages); err != nil {
			return nil, err
		}
	}

	// 解析 Metadata JSON
	if m.Metadata != "" {
		if err := json.Unmarshal([]byte(m.Metadata), &session.Metadata); err != nil {
			return nil, err
		}
	}

	return session, nil
}

// FromEntity 从领域实体创建
func (m *SessionModel) FromEntity(session *entity.Session) error {
	m.ID = session.ID
	m.TenantID = session.TenantID
	m.CreatedAt = session.CreatedAt
	m.UpdatedAt = session.UpdatedAt
	m.ExpiresAt = session.ExpiresAt

	// 序列化 Messages
	if session.Messages != nil && len(session.Messages) > 0 {
		messagesBytes, err := json.Marshal(session.Messages)
		if err != nil {
			return err
		}
		m.Messages = string(messagesBytes)
	}

	// 序列化 Metadata
	if session.Metadata != nil && len(session.Metadata) > 0 {
		metadataBytes, err := json.Marshal(session.Metadata)
		if err != nil {
			return err
		}
		m.Metadata = string(metadataBytes)
	}

	return nil
}

// MissedQueryModel GORM 未命中查询模型
type MissedQueryModel struct {
	ID        uint      `gorm:"primaryKey;autoIncrement"`
	TenantID  string    `gorm:"type:varchar(100);index;not null"`
	Query     string    `gorm:"type:text;not null"`
	Intent    string    `gorm:"type:varchar(50)"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

// TableName 指定表名
func (MissedQueryModel) TableName() string {
	return "missed_queries"
}
