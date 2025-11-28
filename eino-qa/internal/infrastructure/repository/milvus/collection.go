package milvus

import (
	"context"
	"fmt"

	"github.com/milvus-io/milvus-sdk-go/v2/entity"
	"github.com/sirupsen/logrus"
)

// CollectionManager 管理 Milvus Collection
type CollectionManager struct {
	client *Client
	logger *logrus.Logger
}

// NewCollectionManager 创建 Collection 管理器
func NewCollectionManager(client *Client, logger *logrus.Logger) *CollectionManager {
	if logger == nil {
		logger = logrus.New()
	}

	return &CollectionManager{
		client: client,
		logger: logger,
	}
}

// CreateCollection 创建向量集合
func (cm *CollectionManager) CreateCollection(ctx context.Context, collectionName string, dimension int) error {
	cm.logger.WithFields(logrus.Fields{
		"collection": collectionName,
		"dimension":  dimension,
	}).Info("creating Milvus collection")

	// 检查集合是否已存在
	exists, err := cm.CollectionExists(ctx, collectionName)
	if err != nil {
		return fmt.Errorf("failed to check collection existence: %w", err)
	}

	if exists {
		cm.logger.WithField("collection", collectionName).Info("collection already exists")
		return nil
	}

	// 定义 Schema
	schema := &entity.Schema{
		CollectionName: collectionName,
		Description:    "Knowledge base documents for tenant",
		AutoID:         false,
		Fields: []*entity.Field{
			{
				Name:       "id",
				DataType:   entity.FieldTypeVarChar,
				PrimaryKey: true,
				AutoID:     false,
				TypeParams: map[string]string{
					"max_length": "256",
				},
			},
			{
				Name:     "vector",
				DataType: entity.FieldTypeFloatVector,
				TypeParams: map[string]string{
					"dim": fmt.Sprintf("%d", dimension),
				},
			},
			{
				Name:     "content",
				DataType: entity.FieldTypeVarChar,
				TypeParams: map[string]string{
					"max_length": "65535",
				},
			},
			{
				Name:     "metadata",
				DataType: entity.FieldTypeJSON,
			},
			{
				Name:     "tenant_id",
				DataType: entity.FieldTypeVarChar,
				TypeParams: map[string]string{
					"max_length": "128",
				},
			},
			{
				Name:     "created_at",
				DataType: entity.FieldTypeInt64,
			},
		},
	}

	// 创建集合
	err = cm.client.GetClient().CreateCollection(ctx, schema, entity.DefaultShardNumber)
	if err != nil {
		return fmt.Errorf("failed to create collection: %w", err)
	}

	// 创建索引
	idx, err := entity.NewIndexHNSW(entity.L2, 16, 256)
	if err != nil {
		return fmt.Errorf("failed to create index config: %w", err)
	}

	err = cm.client.GetClient().CreateIndex(ctx, collectionName, "vector", idx, false)
	if err != nil {
		return fmt.Errorf("failed to create index: %w", err)
	}

	// 加载集合到内存
	err = cm.client.GetClient().LoadCollection(ctx, collectionName, false)
	if err != nil {
		return fmt.Errorf("failed to load collection: %w", err)
	}

	cm.logger.WithField("collection", collectionName).Info("collection created successfully")
	return nil
}

// CollectionExists 检查集合是否存在
func (cm *CollectionManager) CollectionExists(ctx context.Context, collectionName string) (bool, error) {
	exists, err := cm.client.GetClient().HasCollection(ctx, collectionName)
	if err != nil {
		return false, fmt.Errorf("failed to check collection: %w", err)
	}
	return exists, nil
}

// DropCollection 删除集合
func (cm *CollectionManager) DropCollection(ctx context.Context, collectionName string) error {
	cm.logger.WithField("collection", collectionName).Info("dropping collection")

	err := cm.client.GetClient().DropCollection(ctx, collectionName)
	if err != nil {
		return fmt.Errorf("failed to drop collection: %w", err)
	}

	cm.logger.WithField("collection", collectionName).Info("collection dropped successfully")
	return nil
}

// LoadCollection 加载集合到内存
func (cm *CollectionManager) LoadCollection(ctx context.Context, collectionName string) error {
	err := cm.client.GetClient().LoadCollection(ctx, collectionName, false)
	if err != nil {
		return fmt.Errorf("failed to load collection: %w", err)
	}
	return nil
}

// ReleaseCollection 从内存释放集合
func (cm *CollectionManager) ReleaseCollection(ctx context.Context, collectionName string) error {
	err := cm.client.GetClient().ReleaseCollection(ctx, collectionName)
	if err != nil {
		return fmt.Errorf("failed to release collection: %w", err)
	}
	return nil
}

// GetCollectionStats 获取集合统计信息
func (cm *CollectionManager) GetCollectionStats(ctx context.Context, collectionName string) (map[string]string, error) {
	stats, err := cm.client.GetClient().GetCollectionStatistics(ctx, collectionName)
	if err != nil {
		return nil, fmt.Errorf("failed to get collection stats: %w", err)
	}
	return stats, nil
}
