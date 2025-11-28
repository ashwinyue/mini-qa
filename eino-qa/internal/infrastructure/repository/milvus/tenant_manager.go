package milvus

import (
	"context"
	"fmt"
	"sync"

	"github.com/sirupsen/logrus"
)

// TenantManager 管理多租户的 Collection 映射
type TenantManager struct {
	collectionManager *CollectionManager
	collections       map[string]string // tenantID -> collectionName
	dimension         int
	mu                sync.RWMutex
	logger            *logrus.Logger
}

// NewTenantManager 创建租户管理器
func NewTenantManager(collectionManager *CollectionManager, dimension int, logger *logrus.Logger) *TenantManager {
	if logger == nil {
		logger = logrus.New()
	}

	return &TenantManager{
		collectionManager: collectionManager,
		collections:       make(map[string]string),
		dimension:         dimension,
		logger:            logger,
	}
}

// GetCollection 获取租户对应的 Collection 名称
// 如果 Collection 不存在，会自动创建
func (tm *TenantManager) GetCollection(ctx context.Context, tenantID string) (string, error) {
	// 先尝试从缓存读取
	tm.mu.RLock()
	collectionName, exists := tm.collections[tenantID]
	tm.mu.RUnlock()

	if exists {
		return collectionName, nil
	}

	// 需要创建新的 Collection
	tm.mu.Lock()
	defer tm.mu.Unlock()

	// 双重检查，防止并发创建
	if collectionName, exists := tm.collections[tenantID]; exists {
		return collectionName, nil
	}

	// 生成 Collection 名称
	collectionName = tm.generateCollectionName(tenantID)

	tm.logger.WithFields(logrus.Fields{
		"tenant_id":  tenantID,
		"collection": collectionName,
	}).Info("creating collection for tenant")

	// 创建 Collection
	err := tm.collectionManager.CreateCollection(ctx, collectionName, tm.dimension)
	if err != nil {
		return "", fmt.Errorf("failed to create collection for tenant %s: %w", tenantID, err)
	}

	// 缓存映射关系
	tm.collections[tenantID] = collectionName

	tm.logger.WithFields(logrus.Fields{
		"tenant_id":  tenantID,
		"collection": collectionName,
	}).Info("collection created and cached for tenant")

	return collectionName, nil
}

// CollectionExists 检查租户的 Collection 是否存在
func (tm *TenantManager) CollectionExists(ctx context.Context, tenantID string) (bool, error) {
	tm.mu.RLock()
	collectionName, cached := tm.collections[tenantID]
	tm.mu.RUnlock()

	if cached {
		// 验证 Collection 是否真实存在
		exists, err := tm.collectionManager.CollectionExists(ctx, collectionName)
		if err != nil {
			return false, err
		}
		if !exists {
			// 缓存失效，清除
			tm.mu.Lock()
			delete(tm.collections, tenantID)
			tm.mu.Unlock()
		}
		return exists, nil
	}

	// 检查是否存在对应的 Collection
	collectionName = tm.generateCollectionName(tenantID)
	exists, err := tm.collectionManager.CollectionExists(ctx, collectionName)
	if err != nil {
		return false, err
	}

	if exists {
		// 更新缓存
		tm.mu.Lock()
		tm.collections[tenantID] = collectionName
		tm.mu.Unlock()
	}

	return exists, nil
}

// DropTenantCollection 删除租户的 Collection
func (tm *TenantManager) DropTenantCollection(ctx context.Context, tenantID string) error {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	collectionName, exists := tm.collections[tenantID]
	if !exists {
		collectionName = tm.generateCollectionName(tenantID)
	}

	err := tm.collectionManager.DropCollection(ctx, collectionName)
	if err != nil {
		return err
	}

	// 从缓存中删除
	delete(tm.collections, tenantID)

	tm.logger.WithFields(logrus.Fields{
		"tenant_id":  tenantID,
		"collection": collectionName,
	}).Info("tenant collection dropped")

	return nil
}

// GetAllTenants 获取所有已缓存的租户 ID
func (tm *TenantManager) GetAllTenants() []string {
	tm.mu.RLock()
	defer tm.mu.RUnlock()

	tenants := make([]string, 0, len(tm.collections))
	for tenantID := range tm.collections {
		tenants = append(tenants, tenantID)
	}
	return tenants
}

// ClearCache 清除缓存
func (tm *TenantManager) ClearCache() {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	tm.collections = make(map[string]string)
	tm.logger.Info("tenant collection cache cleared")
}

// generateCollectionName 生成 Collection 名称
func (tm *TenantManager) generateCollectionName(tenantID string) string {
	// 使用前缀 + 租户 ID 作为 Collection 名称
	// 确保名称符合 Milvus 命名规范
	if tenantID == "" || tenantID == "default" {
		return "kb_default"
	}
	return fmt.Sprintf("kb_%s", tenantID)
}
