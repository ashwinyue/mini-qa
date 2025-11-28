package repository

import (
	"context"
	"eino-qa/internal/domain/entity"
	"time"
)

// SessionRepository 定义会话存储操作接口
type SessionRepository interface {
	// Save 保存会话
	// session: 会话实体
	// 返回: 错误
	Save(ctx context.Context, session *entity.Session) error

	// Load 加载会话
	// sessionID: 会话 ID
	// 返回: 会话实体和错误
	Load(ctx context.Context, sessionID string) (*entity.Session, error)

	// Delete 删除会话
	// sessionID: 会话 ID
	// 返回: 错误
	Delete(ctx context.Context, sessionID string) error

	// Exists 检查会话是否存在
	// sessionID: 会话 ID
	// 返回: 是否存在和错误
	Exists(ctx context.Context, sessionID string) (bool, error)

	// AddMessage 向会话添加消息
	// sessionID: 会话 ID
	// message: 消息实体
	// 返回: 错误
	AddMessage(ctx context.Context, sessionID string, message *entity.Message) error

	// GetMessages 获取会话的所有消息
	// sessionID: 会话 ID
	// 返回: 消息列表和错误
	GetMessages(ctx context.Context, sessionID string) ([]*entity.Message, error)

	// UpdateExpiration 更新会话过期时间
	// sessionID: 会话 ID
	// expiresAt: 新的过期时间
	// 返回: 错误
	UpdateExpiration(ctx context.Context, sessionID string, expiresAt time.Time) error

	// DeleteExpired 删除过期的会话
	// 返回: 删除的会话数量和错误
	DeleteExpired(ctx context.Context) (int, error)

	// ListByTenant 列出租户的所有会话
	// tenantID: 租户 ID
	// 返回: 会话列表和错误
	ListByTenant(ctx context.Context, tenantID string) ([]*entity.Session, error)

	// Count 获取会话总数
	// 返回: 会话数量和错误
	Count(ctx context.Context) (int64, error)
}
