package tenant

import (
	"os"
	"path/filepath"
	"testing"

	"eino-qa/internal/domain/entity"
	"eino-qa/internal/infrastructure/repository/sqlite"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestManager(t *testing.T) (*Manager, func()) {
	// 创建临时目录用于测试数据库
	tempDir, err := os.MkdirTemp("", "tenant_test_*")
	require.NoError(t, err)

	// 创建 logger
	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)

	// 创建 DBManager
	dbManager := sqlite.NewDBManager(tempDir)

	// 注意：这里我们无法创建真实的 Milvus 连接用于单元测试
	// 在实际测试中，应该使用 mock 或集成测试
	// 这里我们只测试不依赖 Milvus 的功能

	cleanup := func() {
		dbManager.Close()
		os.RemoveAll(tempDir)
	}

	return &Manager{
		dbManager: dbManager,
		tenants:   make(map[string]*entity.Tenant),
		logger:    logger,
	}, cleanup
}

func TestManager_GetDB(t *testing.T) {
	manager, cleanup := setupTestManager(t)
	defer cleanup()

	tests := []struct {
		name     string
		tenantID string
		wantErr  bool
	}{
		{
			name:     "default tenant",
			tenantID: "default",
			wantErr:  false,
		},
		{
			name:     "custom tenant",
			tenantID: "tenant1",
			wantErr:  false,
		},
		{
			name:     "another tenant",
			tenantID: "tenant2",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 由于没有 Milvus，我们直接测试 DBManager
			db, err := manager.dbManager.GetDB(tt.tenantID)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.NotNil(t, db)

			// 验证数据库连接
			sqlDB, err := db.DB()
			require.NoError(t, err)
			assert.NoError(t, sqlDB.Ping())
		})
	}
}

func TestManager_ListTenants(t *testing.T) {
	manager, cleanup := setupTestManager(t)
	defer cleanup()

	// 初始应该为空
	tenants := manager.ListTenants()
	assert.Empty(t, tenants)

	// 手动添加一些租户到缓存
	manager.mu.Lock()
	manager.tenants["tenant1"] = &entity.Tenant{ID: "tenant1", Name: "Tenant 1"}
	manager.tenants["tenant2"] = &entity.Tenant{ID: "tenant2", Name: "Tenant 2"}
	manager.mu.Unlock()

	// 验证列表
	tenants = manager.ListTenants()
	assert.Len(t, tenants, 2)
	assert.Contains(t, tenants, "tenant1")
	assert.Contains(t, tenants, "tenant2")
}

func TestManager_ClearCache(t *testing.T) {
	manager, cleanup := setupTestManager(t)
	defer cleanup()

	// 添加一些租户到缓存
	manager.mu.Lock()
	manager.tenants["tenant1"] = &entity.Tenant{ID: "tenant1", Name: "Tenant 1"}
	manager.tenants["tenant2"] = &entity.Tenant{ID: "tenant2", Name: "Tenant 2"}
	manager.mu.Unlock()

	// 验证缓存不为空
	assert.Len(t, manager.ListTenants(), 2)

	// 清除缓存
	manager.ClearCache()

	// 验证缓存已清空
	assert.Empty(t, manager.ListTenants())
}

func TestManager_Close(t *testing.T) {
	manager, cleanup := setupTestManager(t)
	defer cleanup()

	// 创建一些数据库连接
	_, err := manager.dbManager.GetDB("tenant1")
	require.NoError(t, err)

	_, err = manager.dbManager.GetDB("tenant2")
	require.NoError(t, err)

	// 关闭管理器
	err = manager.Close()
	assert.NoError(t, err)

	// 验证缓存已清空
	assert.Empty(t, manager.ListTenants())
}

func TestDBManager_MultiTenant(t *testing.T) {
	// 创建临时目录
	tempDir, err := os.MkdirTemp("", "db_test_*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	dbManager := sqlite.NewDBManager(tempDir)
	defer dbManager.Close()

	// 测试多个租户的数据库隔离
	tenants := []string{"tenant1", "tenant2", "tenant3"}

	for _, tenantID := range tenants {
		db, err := dbManager.GetDB(tenantID)
		require.NoError(t, err)
		assert.NotNil(t, db)

		// 验证数据库文件存在
		dbPath := filepath.Join(tempDir, tenantID+".db")
		_, err = os.Stat(dbPath)
		assert.NoError(t, err, "database file should exist for tenant %s", tenantID)
	}

	// 验证租户列表
	listedTenants := dbManager.ListTenants()
	assert.Len(t, listedTenants, len(tenants))
	for _, tenantID := range tenants {
		assert.Contains(t, listedTenants, tenantID)
	}
}

func TestDBManager_ConcurrentAccess(t *testing.T) {
	// 创建临时目录
	tempDir, err := os.MkdirTemp("", "db_concurrent_*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	dbManager := sqlite.NewDBManager(tempDir)
	defer dbManager.Close()

	// 并发访问同一个租户的数据库
	tenantID := "concurrent_tenant"
	done := make(chan bool)

	for i := 0; i < 10; i++ {
		go func() {
			db, err := dbManager.GetDB(tenantID)
			assert.NoError(t, err)
			assert.NotNil(t, db)

			// 执行一个简单的查询
			sqlDB, err := db.DB()
			assert.NoError(t, err)
			assert.NoError(t, sqlDB.Ping())

			done <- true
		}()
	}

	// 等待所有 goroutine 完成
	for i := 0; i < 10; i++ {
		<-done
	}

	// 验证只创建了一个数据库连接
	listedTenants := dbManager.ListTenants()
	assert.Len(t, listedTenants, 1)
	assert.Contains(t, listedTenants, tenantID)
}
