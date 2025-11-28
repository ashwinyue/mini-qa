package sqlite

import (
	"context"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"

	"eino-qa/internal/domain/entity"
	"eino-qa/internal/domain/repository"
)

// SessionRepository SQLite 会话仓储实现
type SessionRepository struct {
	dbManager *DBManager
	tenantID  string
}

// NewSessionRepository 创建会话仓储
func NewSessionRepository(dbManager *DBManager, tenantID string) repository.SessionRepository {
	return &SessionRepository{
		dbManager: dbManager,
		tenantID:  tenantID,
	}
}

// getDB 获取当前租户的数据库连接
func (r *SessionRepository) getDB() (*gorm.DB, error) {
	return r.dbManager.GetDB(r.tenantID)
}

// Save 保存会话
func (r *SessionRepository) Save(ctx context.Context, session *entity.Session) error {
	if err := session.Validate(); err != nil {
		return fmt.Errorf("invalid session: %w", err)
	}

	// 确保租户 ID 匹配
	if session.TenantID != r.tenantID {
		return fmt.Errorf("tenant ID mismatch: expected %s, got %s", r.tenantID, session.TenantID)
	}

	db, err := r.getDB()
	if err != nil {
		return err
	}

	var model SessionModel
	if err := model.FromEntity(session); err != nil {
		return fmt.Errorf("failed to convert session entity: %w", err)
	}

	// 使用 Save 方法，如果存在则更新，不存在则创建
	result := db.WithContext(ctx).Save(&model)
	if result.Error != nil {
		return fmt.Errorf("failed to save session: %w", result.Error)
	}

	return nil
}

// Load 加载会话
func (r *SessionRepository) Load(ctx context.Context, sessionID string) (*entity.Session, error) {
	db, err := r.getDB()
	if err != nil {
		return nil, err
	}

	var model SessionModel
	result := db.WithContext(ctx).
		Where("id = ? AND tenant_id = ?", sessionID, r.tenantID).
		First(&model)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("session not found: %s", sessionID)
		}
		return nil, fmt.Errorf("failed to load session: %w", result.Error)
	}

	return model.ToEntity()
}

// Delete 删除会话
func (r *SessionRepository) Delete(ctx context.Context, sessionID string) error {
	db, err := r.getDB()
	if err != nil {
		return err
	}

	result := db.WithContext(ctx).
		Where("id = ? AND tenant_id = ?", sessionID, r.tenantID).
		Delete(&SessionModel{})

	if result.Error != nil {
		return fmt.Errorf("failed to delete session: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("session not found: %s", sessionID)
	}

	return nil
}

// Exists 检查会话是否存在
func (r *SessionRepository) Exists(ctx context.Context, sessionID string) (bool, error) {
	db, err := r.getDB()
	if err != nil {
		return false, err
	}

	var count int64
	result := db.WithContext(ctx).
		Model(&SessionModel{}).
		Where("id = ? AND tenant_id = ?", sessionID, r.tenantID).
		Count(&count)

	if result.Error != nil {
		return false, fmt.Errorf("failed to check session existence: %w", result.Error)
	}

	return count > 0, nil
}

// AddMessage 向会话添加消息
func (r *SessionRepository) AddMessage(ctx context.Context, sessionID string, message *entity.Message) error {
	if err := message.Validate(); err != nil {
		return fmt.Errorf("invalid message: %w", err)
	}

	// 加载现有会话
	session, err := r.Load(ctx, sessionID)
	if err != nil {
		return err
	}

	// 添加消息
	if err := session.AddMessage(message); err != nil {
		return fmt.Errorf("failed to add message to session: %w", err)
	}

	// 保存更新后的会话
	return r.Save(ctx, session)
}

// GetMessages 获取会话的所有消息
func (r *SessionRepository) GetMessages(ctx context.Context, sessionID string) ([]*entity.Message, error) {
	session, err := r.Load(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	return session.GetMessages(), nil
}

// UpdateExpiration 更新会话过期时间
func (r *SessionRepository) UpdateExpiration(ctx context.Context, sessionID string, expiresAt time.Time) error {
	db, err := r.getDB()
	if err != nil {
		return err
	}

	result := db.WithContext(ctx).
		Model(&SessionModel{}).
		Where("id = ? AND tenant_id = ?", sessionID, r.tenantID).
		Updates(map[string]interface{}{
			"expires_at": expiresAt,
			"updated_at": time.Now(),
		})

	if result.Error != nil {
		return fmt.Errorf("failed to update session expiration: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("session not found: %s", sessionID)
	}

	return nil
}

// DeleteExpired 删除过期的会话
func (r *SessionRepository) DeleteExpired(ctx context.Context) (int, error) {
	db, err := r.getDB()
	if err != nil {
		return 0, err
	}

	result := db.WithContext(ctx).
		Where("expires_at < ? AND tenant_id = ?", time.Now(), r.tenantID).
		Delete(&SessionModel{})

	if result.Error != nil {
		return 0, fmt.Errorf("failed to delete expired sessions: %w", result.Error)
	}

	return int(result.RowsAffected), nil
}

// ListByTenant 列出租户的所有会话
func (r *SessionRepository) ListByTenant(ctx context.Context, tenantID string) ([]*entity.Session, error) {
	// 确保租户 ID 匹配
	if tenantID != r.tenantID {
		return nil, fmt.Errorf("tenant ID mismatch: expected %s, got %s", r.tenantID, tenantID)
	}

	db, err := r.getDB()
	if err != nil {
		return nil, err
	}

	var models []SessionModel
	result := db.WithContext(ctx).
		Where("tenant_id = ?", tenantID).
		Order("updated_at DESC").
		Find(&models)

	if result.Error != nil {
		return nil, fmt.Errorf("failed to list sessions: %w", result.Error)
	}

	sessions := make([]*entity.Session, 0, len(models))
	for _, model := range models {
		session, err := model.ToEntity()
		if err != nil {
			return nil, fmt.Errorf("failed to convert session model: %w", err)
		}
		sessions = append(sessions, session)
	}

	return sessions, nil
}

// Count 获取会话总数
func (r *SessionRepository) Count(ctx context.Context) (int64, error) {
	db, err := r.getDB()
	if err != nil {
		return 0, err
	}

	var count int64
	result := db.WithContext(ctx).
		Model(&SessionModel{}).
		Where("tenant_id = ?", r.tenantID).
		Count(&count)

	if result.Error != nil {
		return 0, fmt.Errorf("failed to count sessions: %w", result.Error)
	}

	return count, nil
}
