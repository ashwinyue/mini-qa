package sqlite

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"

	"eino-qa/internal/domain/entity"
	"eino-qa/internal/domain/repository"
)

// OrderRepository SQLite 订单仓储实现
type OrderRepository struct {
	dbManager *DBManager
	tenantID  string
}

// NewOrderRepository 创建订单仓储
func NewOrderRepository(dbManager *DBManager, tenantID string) repository.OrderRepository {
	return &OrderRepository{
		dbManager: dbManager,
		tenantID:  tenantID,
	}
}

// getDB 获取当前租户的数据库连接
func (r *OrderRepository) getDB() (*gorm.DB, error) {
	return r.dbManager.GetDB(r.tenantID)
}

// FindByID 根据订单 ID 查询订单
func (r *OrderRepository) FindByID(ctx context.Context, orderID string) (*entity.Order, error) {
	db, err := r.getDB()
	if err != nil {
		return nil, err
	}

	var model OrderModel
	result := db.WithContext(ctx).Where("id = ? AND tenant_id = ?", orderID, r.tenantID).First(&model)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("order not found: %s", orderID)
		}
		return nil, fmt.Errorf("failed to find order: %w", result.Error)
	}

	return model.ToEntity()
}

// FindByUserID 根据用户 ID 查询订单列表
func (r *OrderRepository) FindByUserID(ctx context.Context, userID string) ([]*entity.Order, error) {
	db, err := r.getDB()
	if err != nil {
		return nil, err
	}

	var models []OrderModel
	result := db.WithContext(ctx).
		Where("user_id = ? AND tenant_id = ?", userID, r.tenantID).
		Order("created_at DESC").
		Find(&models)

	if result.Error != nil {
		return nil, fmt.Errorf("failed to find orders by user: %w", result.Error)
	}

	orders := make([]*entity.Order, 0, len(models))
	for _, model := range models {
		order, err := model.ToEntity()
		if err != nil {
			return nil, fmt.Errorf("failed to convert order model: %w", err)
		}
		orders = append(orders, order)
	}

	return orders, nil
}

// FindByStatus 根据订单状态查询订单列表
func (r *OrderRepository) FindByStatus(ctx context.Context, status entity.OrderStatus) ([]*entity.Order, error) {
	db, err := r.getDB()
	if err != nil {
		return nil, err
	}

	var models []OrderModel
	result := db.WithContext(ctx).
		Where("status = ? AND tenant_id = ?", string(status), r.tenantID).
		Order("created_at DESC").
		Find(&models)

	if result.Error != nil {
		return nil, fmt.Errorf("failed to find orders by status: %w", result.Error)
	}

	orders := make([]*entity.Order, 0, len(models))
	for _, model := range models {
		order, err := model.ToEntity()
		if err != nil {
			return nil, fmt.Errorf("failed to convert order model: %w", err)
		}
		orders = append(orders, order)
	}

	return orders, nil
}

// Create 创建新订单
func (r *OrderRepository) Create(ctx context.Context, order *entity.Order) error {
	if err := order.Validate(); err != nil {
		return fmt.Errorf("invalid order: %w", err)
	}

	// 确保租户 ID 匹配
	if order.TenantID != r.tenantID {
		return fmt.Errorf("tenant ID mismatch: expected %s, got %s", r.tenantID, order.TenantID)
	}

	db, err := r.getDB()
	if err != nil {
		return err
	}

	var model OrderModel
	if err := model.FromEntity(order); err != nil {
		return fmt.Errorf("failed to convert order entity: %w", err)
	}

	result := db.WithContext(ctx).Create(&model)
	if result.Error != nil {
		return fmt.Errorf("failed to create order: %w", result.Error)
	}

	return nil
}

// Update 更新订单
func (r *OrderRepository) Update(ctx context.Context, order *entity.Order) error {
	if err := order.Validate(); err != nil {
		return fmt.Errorf("invalid order: %w", err)
	}

	// 确保租户 ID 匹配
	if order.TenantID != r.tenantID {
		return fmt.Errorf("tenant ID mismatch: expected %s, got %s", r.tenantID, order.TenantID)
	}

	db, err := r.getDB()
	if err != nil {
		return err
	}

	var model OrderModel
	if err := model.FromEntity(order); err != nil {
		return fmt.Errorf("failed to convert order entity: %w", err)
	}

	result := db.WithContext(ctx).
		Where("id = ? AND tenant_id = ?", order.ID, r.tenantID).
		Updates(&model)

	if result.Error != nil {
		return fmt.Errorf("failed to update order: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("order not found: %s", order.ID)
	}

	return nil
}

// Delete 删除订单
func (r *OrderRepository) Delete(ctx context.Context, orderID string) error {
	db, err := r.getDB()
	if err != nil {
		return err
	}

	result := db.WithContext(ctx).
		Where("id = ? AND tenant_id = ?", orderID, r.tenantID).
		Delete(&OrderModel{})

	if result.Error != nil {
		return fmt.Errorf("failed to delete order: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("order not found: %s", orderID)
	}

	return nil
}

// List 列出所有订单（支持分页）
func (r *OrderRepository) List(ctx context.Context, offset, limit int) ([]*entity.Order, error) {
	db, err := r.getDB()
	if err != nil {
		return nil, err
	}

	var models []OrderModel
	result := db.WithContext(ctx).
		Where("tenant_id = ?", r.tenantID).
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&models)

	if result.Error != nil {
		return nil, fmt.Errorf("failed to list orders: %w", result.Error)
	}

	orders := make([]*entity.Order, 0, len(models))
	for _, model := range models {
		order, err := model.ToEntity()
		if err != nil {
			return nil, fmt.Errorf("failed to convert order model: %w", err)
		}
		orders = append(orders, order)
	}

	return orders, nil
}

// Count 获取订单总数
func (r *OrderRepository) Count(ctx context.Context) (int64, error) {
	db, err := r.getDB()
	if err != nil {
		return 0, err
	}

	var count int64
	result := db.WithContext(ctx).
		Model(&OrderModel{}).
		Where("tenant_id = ?", r.tenantID).
		Count(&count)

	if result.Error != nil {
		return 0, fmt.Errorf("failed to count orders: %w", result.Error)
	}

	return count, nil
}
