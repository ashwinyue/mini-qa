package milvus

import (
	"eino-qa/internal/domain/repository"
	"eino-qa/internal/infrastructure/config"

	"github.com/sirupsen/logrus"
)

// Factory Milvus 仓储工厂
type Factory struct {
	client        *Client
	tenantManager *TenantManager
	logger        *logrus.Logger
}

// NewFactory 创建 Milvus 仓储工厂
func NewFactory(cfg config.MilvusConfig, dimension int, logger *logrus.Logger) (*Factory, error) {
	if logger == nil {
		logger = logrus.New()
	}

	// 创建客户端
	client, err := NewClient(ClientConfig{
		Host:     cfg.Host,
		Port:     cfg.Port,
		Username: cfg.Username,
		Password: cfg.Password,
		Timeout:  cfg.Timeout,
	}, logger)
	if err != nil {
		return nil, err
	}

	// 创建 Collection 管理器
	collectionManager := NewCollectionManager(client, logger)

	// 创建租户管理器
	tenantManager := NewTenantManager(collectionManager, dimension, logger)

	return &Factory{
		client:        client,
		tenantManager: tenantManager,
		logger:        logger,
	}, nil
}

// CreateVectorRepository 创建向量仓储
func (f *Factory) CreateVectorRepository() repository.VectorRepository {
	return NewVectorRepository(f.client, f.tenantManager, f.logger)
}

// GetClient 获取 Milvus 客户端
func (f *Factory) GetClient() *Client {
	return f.client
}

// GetTenantManager 获取租户管理器
func (f *Factory) GetTenantManager() *TenantManager {
	return f.tenantManager
}

// Close 关闭连接
func (f *Factory) Close() error {
	return f.client.Close()
}
