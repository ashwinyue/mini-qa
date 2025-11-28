package tenant

import (
	"fmt"
	"time"

	"eino-qa/internal/infrastructure/config"
	"eino-qa/internal/infrastructure/repository/milvus"
	"eino-qa/internal/infrastructure/repository/sqlite"

	"github.com/sirupsen/logrus"
)

// FactoryConfig 租户管理器工厂配置
type FactoryConfig struct {
	Config *config.Config
	Logger *logrus.Logger
}

// NewManagerFromConfig 从配置创建租户管理器
func NewManagerFromConfig(cfg FactoryConfig) (*Manager, error) {
	if cfg.Logger == nil {
		cfg.Logger = logrus.New()
	}

	// 创建 Milvus 客户端
	milvusClient, err := milvus.NewClient(milvus.ClientConfig{
		Host:     cfg.Config.Milvus.Host,
		Port:     cfg.Config.Milvus.Port,
		Username: cfg.Config.Milvus.Username,
		Password: cfg.Config.Milvus.Password,
		Timeout:  30 * time.Second,
	}, cfg.Logger)
	if err != nil {
		return nil, fmt.Errorf("failed to create Milvus client: %w", err)
	}

	// 创建 Collection 管理器
	collectionManager := milvus.NewCollectionManager(milvusClient, cfg.Logger)

	// 创建 Milvus 租户管理器
	milvusTenantMgr := milvus.NewTenantManager(
		collectionManager,
		cfg.Config.DashScope.EmbeddingDimension,
		cfg.Logger,
	)

	// 创建 SQLite 数据库管理器
	dbManager := sqlite.NewDBManager(cfg.Config.Database.BasePath)

	// 创建统一的租户管理器
	manager := NewManager(Config{
		MilvusTenantManager: milvusTenantMgr,
		DBManager:           dbManager,
		Logger:              cfg.Logger,
	})

	cfg.Logger.Info("tenant manager created successfully")

	return manager, nil
}
