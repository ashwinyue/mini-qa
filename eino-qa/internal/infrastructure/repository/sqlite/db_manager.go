package sqlite

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DBManager 管理多租户数据库连接
type DBManager struct {
	basePath string
	dbs      map[string]*gorm.DB
	mu       sync.RWMutex
	config   *gorm.Config
}

// NewDBManager 创建数据库管理器
func NewDBManager(basePath string) *DBManager {
	return &DBManager{
		basePath: basePath,
		dbs:      make(map[string]*gorm.DB),
		config: &gorm.Config{
			Logger: logger.Default.LogMode(logger.Silent),
		},
	}
}

// GetDB 获取租户的数据库连接
func (m *DBManager) GetDB(tenantID string) (*gorm.DB, error) {
	// 先尝试读锁获取已存在的连接
	m.mu.RLock()
	db, exists := m.dbs[tenantID]
	m.mu.RUnlock()

	if exists {
		return db, nil
	}

	// 使用写锁创建新连接
	m.mu.Lock()
	defer m.mu.Unlock()

	// 双重检查，防止并发创建
	if db, exists := m.dbs[tenantID]; exists {
		return db, nil
	}

	// 创建租户数据库文件路径
	dbPath := m.getDBPath(tenantID)

	// 确保目录存在
	if err := os.MkdirAll(filepath.Dir(dbPath), 0755); err != nil {
		return nil, fmt.Errorf("failed to create database directory: %w", err)
	}

	// 打开数据库连接
	db, err := gorm.Open(sqlite.Open(dbPath), m.config)
	if err != nil {
		return nil, fmt.Errorf("failed to open database for tenant %s: %w", tenantID, err)
	}

	// 自动迁移表结构
	if err := m.autoMigrate(db); err != nil {
		return nil, fmt.Errorf("failed to migrate database for tenant %s: %w", tenantID, err)
	}

	// 缓存连接
	m.dbs[tenantID] = db

	return db, nil
}

// getDBPath 获取租户数据库文件路径
func (m *DBManager) getDBPath(tenantID string) string {
	// 使用租户 ID 作为数据库文件名
	// 格式: {basePath}/{tenantID}.db
	return filepath.Join(m.basePath, fmt.Sprintf("%s.db", tenantID))
}

// autoMigrate 自动迁移表结构
func (m *DBManager) autoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&OrderModel{},
		&SessionModel{},
		&MissedQueryModel{},
	)
}

// Close 关闭所有数据库连接
func (m *DBManager) Close() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	var errs []error
	for tenantID, db := range m.dbs {
		sqlDB, err := db.DB()
		if err != nil {
			errs = append(errs, fmt.Errorf("failed to get sql.DB for tenant %s: %w", tenantID, err))
			continue
		}

		if err := sqlDB.Close(); err != nil {
			errs = append(errs, fmt.Errorf("failed to close database for tenant %s: %w", tenantID, err))
		}
	}

	// 清空连接缓存
	m.dbs = make(map[string]*gorm.DB)

	if len(errs) > 0 {
		return fmt.Errorf("errors closing databases: %v", errs)
	}

	return nil
}

// RemoveDB 移除租户的数据库连接（用于测试或清理）
func (m *DBManager) RemoveDB(tenantID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	db, exists := m.dbs[tenantID]
	if !exists {
		return nil
	}

	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get sql.DB for tenant %s: %w", tenantID, err)
	}

	if err := sqlDB.Close(); err != nil {
		return fmt.Errorf("failed to close database for tenant %s: %w", tenantID, err)
	}

	delete(m.dbs, tenantID)

	return nil
}

// ListTenants 列出所有已连接的租户
func (m *DBManager) ListTenants() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	tenants := make([]string, 0, len(m.dbs))
	for tenantID := range m.dbs {
		tenants = append(tenants, tenantID)
	}

	return tenants
}

// SetLogLevel 设置日志级别
func (m *DBManager) SetLogLevel(level logger.LogLevel) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.config.Logger = logger.Default.LogMode(level)

	// 更新所有已存在的连接
	for _, db := range m.dbs {
		db.Logger = m.config.Logger
	}
}
