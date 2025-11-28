package entity

import "time"

// OrderStatus 定义订单状态
type OrderStatus string

const (
	// OrderStatusPending 待支付
	OrderStatusPending OrderStatus = "pending"
	// OrderStatusPaid 已支付
	OrderStatusPaid OrderStatus = "paid"
	// OrderStatusRefunded 已退款
	OrderStatusRefunded OrderStatus = "refunded"
	// OrderStatusCancelled 已取消
	OrderStatusCancelled OrderStatus = "cancelled"
)

// Order 表示订单实体
type Order struct {
	ID         string
	UserID     string
	CourseName string
	Amount     float64
	Status     OrderStatus
	TenantID   string
	CreatedAt  time.Time
	UpdatedAt  time.Time
	Metadata   map[string]any
}

// NewOrder 创建新的订单实例
func NewOrder(userID, courseName string, amount float64, tenantID string) *Order {
	now := time.Now()
	return &Order{
		ID:         generateOrderID(),
		UserID:     userID,
		CourseName: courseName,
		Amount:     amount,
		Status:     OrderStatusPending,
		TenantID:   tenantID,
		CreatedAt:  now,
		UpdatedAt:  now,
		Metadata:   make(map[string]any),
	}
}

// Validate 验证订单的有效性
func (o *Order) Validate() error {
	if o.UserID == "" {
		return ErrEmptyUserID
	}

	if o.CourseName == "" {
		return ErrEmptyCourseName
	}

	if o.Amount < 0 {
		return ErrInvalidAmount
	}

	if o.TenantID == "" {
		return ErrEmptyTenantID
	}

	validStatuses := map[OrderStatus]bool{
		OrderStatusPending:   true,
		OrderStatusPaid:      true,
		OrderStatusRefunded:  true,
		OrderStatusCancelled: true,
	}

	if !validStatuses[o.Status] {
		return ErrInvalidOrderStatus
	}

	return nil
}

// IsPending 判断订单是否待支付
func (o *Order) IsPending() bool {
	return o.Status == OrderStatusPending
}

// IsPaid 判断订单是否已支付
func (o *Order) IsPaid() bool {
	return o.Status == OrderStatusPaid
}

// IsRefunded 判断订单是否已退款
func (o *Order) IsRefunded() bool {
	return o.Status == OrderStatusRefunded
}

// IsCancelled 判断订单是否已取消
func (o *Order) IsCancelled() bool {
	return o.Status == OrderStatusCancelled
}

// UpdateStatus 更新订单状态
func (o *Order) UpdateStatus(status OrderStatus) error {
	validStatuses := map[OrderStatus]bool{
		OrderStatusPending:   true,
		OrderStatusPaid:      true,
		OrderStatusRefunded:  true,
		OrderStatusCancelled: true,
	}

	if !validStatuses[status] {
		return ErrInvalidOrderStatus
	}

	o.Status = status
	o.UpdatedAt = time.Now()
	return nil
}

// AddMetadata 添加元数据
func (o *Order) AddMetadata(key string, value any) {
	if o.Metadata == nil {
		o.Metadata = make(map[string]any)
	}
	o.Metadata[key] = value
}

// generateOrderID 生成订单 ID
func generateOrderID() string {
	// 格式: #YYYYMMDDXXX
	return generateUniqueID("#"+time.Now().Format("20060102"), 6)
}
