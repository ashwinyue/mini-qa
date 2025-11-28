package sqlite

import (
	"eino-qa/internal/domain/repository"
)

// RepositoryFactory SQLite 仓储工厂
type RepositoryFactory struct {
	dbManager *DBManager
}

// NewRepositoryFactory 创建仓储工厂
func NewRepositoryFactory(basePath string) *RepositoryFactory {
	return &RepositoryFactory{
		dbManager: NewDBManager(basePath),
	}
}

// GetOrderRepository 获取订单仓储
func (f *RepositoryFactory) GetOrderRepository(tenantID string) repository.OrderRepository {
	return NewOrderRepository(f.dbManager, tenantID)
}

// GetSessionRepository 获取会话仓储
func (f *RepositoryFactory) GetSessionRepository(tenantID string) repository.SessionRepository {
	return NewSessionRepository(f.dbManager, tenantID)
}

// GetMissedQueryRepository 获取未命中查询仓储
func (f *RepositoryFactory) GetMissedQueryRepository(tenantID string) *MissedQueryRepository {
	return NewMissedQueryRepository(f.dbManager, tenantID)
}

// GetDBManager 获取数据库管理器
func (f *RepositoryFactory) GetDBManager() *DBManager {
	return f.dbManager
}

// Close 关闭所有数据库连接
func (f *RepositoryFactory) Close() error {
	return f.dbManager.Close()
}
