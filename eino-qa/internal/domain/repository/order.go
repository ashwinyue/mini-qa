package repository

import (
	"context"
	"eino-qa/internal/domain/entity"
)

// OrderRepository 定义订单数据库操作接口
type OrderRepository interface {
	// FindByID 根据订单 ID 查询订单
	// orderID: 订单 ID
	// 返回: 订单实体和错误
	FindByID(ctx context.Context, orderID string) (*entity.Order, error)

	// FindByUserID 根据用户 ID 查询订单列表
	// userID: 用户 ID
	// 返回: 订单列表和错误
	FindByUserID(ctx context.Context, userID string) ([]*entity.Order, error)

	// FindByStatus 根据订单状态查询订单列表
	// status: 订单状态
	// 返回: 订单列表和错误
	FindByStatus(ctx context.Context, status entity.OrderStatus) ([]*entity.Order, error)

	// Create 创建新订单
	// order: 订单实体
	// 返回: 错误
	Create(ctx context.Context, order *entity.Order) error

	// Update 更新订单
	// order: 订单实体
	// 返回: 错误
	Update(ctx context.Context, order *entity.Order) error

	// Delete 删除订单
	// orderID: 订单 ID
	// 返回: 错误
	Delete(ctx context.Context, orderID string) error

	// List 列出所有订单（支持分页）
	// offset: 偏移量
	// limit: 限制数量
	// 返回: 订单列表和错误
	List(ctx context.Context, offset, limit int) ([]*entity.Order, error)

	// Count 获取订单总数
	// 返回: 订单数量和错误
	Count(ctx context.Context) (int64, error)
}
