package sqlite

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"
)

// MissedQuery 未命中查询实体
type MissedQuery struct {
	ID        uint
	TenantID  string
	Query     string
	Intent    string
	CreatedAt time.Time
}

// MissedQueryRepository 未命中查询仓储
type MissedQueryRepository struct {
	dbManager *DBManager
	tenantID  string
}

// NewMissedQueryRepository 创建未命中查询仓储
func NewMissedQueryRepository(dbManager *DBManager, tenantID string) *MissedQueryRepository {
	return &MissedQueryRepository{
		dbManager: dbManager,
		tenantID:  tenantID,
	}
}

// getDB 获取当前租户的数据库连接
func (r *MissedQueryRepository) getDB() (*gorm.DB, error) {
	return r.dbManager.GetDB(r.tenantID)
}

// Create 创建未命中查询记录
func (r *MissedQueryRepository) Create(ctx context.Context, query, intent string) error {
	db, err := r.getDB()
	if err != nil {
		return err
	}

	model := MissedQueryModel{
		TenantID: r.tenantID,
		Query:    query,
		Intent:   intent,
	}

	result := db.WithContext(ctx).Create(&model)
	if result.Error != nil {
		return fmt.Errorf("failed to create missed query: %w", result.Error)
	}

	return nil
}

// List 列出未命中查询（支持分页）
func (r *MissedQueryRepository) List(ctx context.Context, offset, limit int) ([]*MissedQuery, error) {
	db, err := r.getDB()
	if err != nil {
		return nil, err
	}

	var models []MissedQueryModel
	result := db.WithContext(ctx).
		Where("tenant_id = ?", r.tenantID).
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&models)

	if result.Error != nil {
		return nil, fmt.Errorf("failed to list missed queries: %w", result.Error)
	}

	queries := make([]*MissedQuery, 0, len(models))
	for _, model := range models {
		queries = append(queries, &MissedQuery{
			ID:        model.ID,
			TenantID:  model.TenantID,
			Query:     model.Query,
			Intent:    model.Intent,
			CreatedAt: model.CreatedAt,
		})
	}

	return queries, nil
}

// Count 获取未命中查询总数
func (r *MissedQueryRepository) Count(ctx context.Context) (int64, error) {
	db, err := r.getDB()
	if err != nil {
		return 0, err
	}

	var count int64
	result := db.WithContext(ctx).
		Model(&MissedQueryModel{}).
		Where("tenant_id = ?", r.tenantID).
		Count(&count)

	if result.Error != nil {
		return 0, fmt.Errorf("failed to count missed queries: %w", result.Error)
	}

	return count, nil
}

// DeleteOlderThan 删除指定时间之前的记录
func (r *MissedQueryRepository) DeleteOlderThan(ctx context.Context, before time.Time) (int, error) {
	db, err := r.getDB()
	if err != nil {
		return 0, err
	}

	result := db.WithContext(ctx).
		Where("created_at < ? AND tenant_id = ?", before, r.tenantID).
		Delete(&MissedQueryModel{})

	if result.Error != nil {
		return 0, fmt.Errorf("failed to delete old missed queries: %w", result.Error)
	}

	return int(result.RowsAffected), nil
}
