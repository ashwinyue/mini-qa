package tenant

import (
	"context"
	"fmt"
	"sync"

	"eino-qa/internal/domain/entity"
	"eino-qa/internal/infrastructure/repository/milvus"
	"eino-qa/internal/infrastructure/repository/sqlite"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// Manager 统一的多租户管理器
// 负责管理租户的 Milvus Collection 和 SQLite 数据库映射
type Manager struct {
	milvusTenantMgr *milvus.TenantManager
	dbManager       *sqlite.DBManager
	tenants         map[string]*entity.Tenant
	mu              sync.RWMutex
	logger          *logrus.Logger
}

// Config 租户管理器配置
type Config struct {
	MilvusTenantManager *milvus.TenantManager
	DBManager           *sqlite.DBManager
	Logger              *logrus.Logger
}

// NewManager 创建租户管理器
func NewManager(config Config) *Manager {
	if config.Logger == nil {
		config.Logger = logrus.New()
	}

	return &Manager{
		milvusTenantMgr: config.MilvusTenantManager,
		dbManager:       config.DBManager,
		tenants:         make(map[string]*entity.Tenant),
		logger:          config.Logger,
	}
}

// GetTenant 获取租户信息
// 如果租户不存在，会自动创建租户及其资源
func (m *Manager) GetTenant(ctx context.Context, tenantID string) (*entity.Tenant, error) {
	// 标准化租户 ID
	if tenantID == "" {
		tenantID = "default"
	}

	// 先尝试从缓存读取
	m.mu.RLock()
	tenant, exists := m.tenants[tenantID]
	m.mu.RUnlock()

	if exists {
		return tenant, nil
	}

	// 需要创建新租户
	m.mu.Lock()
	defer m.mu.Unlock()

	// 双重检查，防止并发创建
	if tenant, exists := m.tenants[tenantID]; exists {
		return tenant, nil
	}

	m.logger.WithField("tenant_id", tenantID).Info("initializing new tenant")

	// 创建租户实体
	tenant = entity.NewTenant(tenantID, tenantID)
	if err := tenant.Validate(); err != nil {
		return nil, fmt.Errorf("invalid tenant: %w", err)
	}

	// 确保租户的 Milvus Collection 存在
	collectionName, err := m.milvusTenantMgr.GetCollection(ctx, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get/create collection for tenant %s: %w", tenantID, err)
	}
	tenant.CollectionName = collectionName

	// 确保租户的 SQLite 数据库存在
	db, err := m.dbManager.GetDB(tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get/create database for tenant %s: %w", tenantID, err)
	}

	// 验证数据库连接
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get sql.DB for tenant %s: %w", tenantID, err)
	}
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database for tenant %s: %w", tenantID, err)
	}

	// 缓存租户信息
	m.tenants[tenantID] = tenant

	m.logger.WithFields(logrus.Fields{
		"tenant_id":  tenantID,
		"collection": collectionName,
		"database":   tenant.DatabasePath,
	}).Info("tenant initialized successfully")

	return tenant, nil
}

// GetCollection 获取租户的 Milvus Collection 名称
func (m *Manager) GetCollection(ctx context.Context, tenantID string) (string, error) {
	tenant, err := m.GetTenant(ctx, tenantID)
	if err != nil {
		return "", err
	}
	return tenant.CollectionName, nil
}

// GetDB 获取租户的数据库连接
func (m *Manager) GetDB(ctx context.Context, tenantID string) (*gorm.DB, error) {
	// 确保租户已初始化
	_, err := m.GetTenant(ctx, tenantID)
	if err != nil {
		return nil, err
	}

	// 从 DBManager 获取数据库连接
	return m.dbManager.GetDB(tenantID)
}

// TenantExists 检查租户是否存在
func (m *Manager) TenantExists(ctx context.Context, tenantID string) (bool, error) {
	// 检查缓存
	m.mu.RLock()
	_, cached := m.tenants[tenantID]
	m.mu.RUnlock()

	if cached {
		return true, nil
	}

	// 检查 Milvus Collection 是否存在
	collectionExists, err := m.milvusTenantMgr.CollectionExists(ctx, tenantID)
	if err != nil {
		return false, fmt.Errorf("failed to check collection existence: %w", err)
	}

	if !collectionExists {
		return false, nil
	}

	// 检查数据库文件是否存在（通过尝试获取连接）
	_, err = m.dbManager.GetDB(tenantID)
	if err != nil {
		return false, nil
	}

	return true, nil
}

// CreateTenant 显式创建租户及其资源
func (m *Manager) CreateTenant(ctx context.Context, tenantID, name string) (*entity.Tenant, error) {
	if tenantID == "" {
		return nil, fmt.Errorf("tenant ID cannot be empty")
	}

	// 检查租户是否已存在
	exists, err := m.TenantExists(ctx, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to check tenant existence: %w", err)
	}

	if exists {
		m.logger.WithField("tenant_id", tenantID).Info("tenant already exists")
		return m.GetTenant(ctx, tenantID)
	}

	m.logger.WithFields(logrus.Fields{
		"tenant_id": tenantID,
		"name":      name,
	}).Info("creating new tenant")

	// 创建租户实体
	tenant := entity.NewTenant(tenantID, name)
	if err := tenant.Validate(); err != nil {
		return nil, fmt.Errorf("invalid tenant: %w", err)
	}

	// 创建 Milvus Collection
	collectionName, err := m.milvusTenantMgr.GetCollection(ctx, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to create collection: %w", err)
	}
	tenant.CollectionName = collectionName

	// 创建 SQLite 数据库
	db, err := m.dbManager.GetDB(tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to create database: %w", err)
	}

	// 验证数据库连接
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get sql.DB: %w", err)
	}
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// 缓存租户信息
	m.mu.Lock()
	m.tenants[tenantID] = tenant
	m.mu.Unlock()

	m.logger.WithFields(logrus.Fields{
		"tenant_id":  tenantID,
		"collection": collectionName,
		"database":   tenant.DatabasePath,
	}).Info("tenant created successfully")

	return tenant, nil
}

// DeleteTenant 删除租户及其所有资源
func (m *Manager) DeleteTenant(ctx context.Context, tenantID string) error {
	if tenantID == "" || tenantID == "default" {
		return fmt.Errorf("cannot delete default tenant")
	}

	m.logger.WithField("tenant_id", tenantID).Info("deleting tenant")

	// 删除 Milvus Collection
	if err := m.milvusTenantMgr.DropTenantCollection(ctx, tenantID); err != nil {
		m.logger.WithError(err).Warn("failed to drop tenant collection")
		// 继续删除其他资源
	}

	// 关闭并删除数据库连接
	if err := m.dbManager.RemoveDB(tenantID); err != nil {
		m.logger.WithError(err).Warn("failed to remove tenant database")
		// 继续删除其他资源
	}

	// 从缓存中删除
	m.mu.Lock()
	delete(m.tenants, tenantID)
	m.mu.Unlock()

	m.logger.WithField("tenant_id", tenantID).Info("tenant deleted successfully")

	return nil
}

// ListTenants 列出所有已缓存的租户
func (m *Manager) ListTenants() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	tenants := make([]string, 0, len(m.tenants))
	for tenantID := range m.tenants {
		tenants = append(tenants, tenantID)
	}
	return tenants
}

// GetTenantInfo 获取租户的详细信息
func (m *Manager) GetTenantInfo(ctx context.Context, tenantID string) (*TenantInfo, error) {
	tenant, err := m.GetTenant(ctx, tenantID)
	if err != nil {
		return nil, err
	}

	info := &TenantInfo{
		TenantID:       tenant.ID,
		Name:           tenant.Name,
		CollectionName: tenant.CollectionName,
		DatabasePath:   tenant.DatabasePath,
		Metadata:       tenant.Metadata,
	}

	// 获取 Collection 统计信息
	if m.milvusTenantMgr != nil {
		exists, _ := m.milvusTenantMgr.CollectionExists(ctx, tenantID)
		info.CollectionExists = exists
	}

	// 获取数据库统计信息
	db, err := m.dbManager.GetDB(tenantID)
	if err == nil {
		sqlDB, err := db.DB()
		if err == nil {
			stats := sqlDB.Stats()
			info.DBConnections = stats.OpenConnections
		}
	}

	return info, nil
}

// ClearCache 清除租户缓存
func (m *Manager) ClearCache() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.tenants = make(map[string]*entity.Tenant)
	m.logger.Info("tenant cache cleared")
}

// Close 关闭租户管理器，释放所有资源
func (m *Manager) Close() error {
	m.logger.Info("closing tenant manager")

	// 关闭所有数据库连接
	if err := m.dbManager.Close(); err != nil {
		m.logger.WithError(err).Error("failed to close database manager")
		return err
	}

	// 清除缓存
	m.ClearCache()

	m.logger.Info("tenant manager closed successfully")
	return nil
}

// TenantInfo 租户信息
type TenantInfo struct {
	TenantID         string         `json:"tenant_id"`
	Name             string         `json:"name"`
	CollectionName   string         `json:"collection_name"`
	CollectionExists bool           `json:"collection_exists"`
	DatabasePath     string         `json:"database_path"`
	DBConnections    int            `json:"db_connections"`
	Metadata         map[string]any `json:"metadata,omitempty"`
}
