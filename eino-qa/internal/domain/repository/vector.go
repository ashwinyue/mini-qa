package repository

import (
	"context"
	"eino-qa/internal/domain/entity"
)

// VectorRepository 定义向量数据库操作接口
type VectorRepository interface {
	// Search 执行向量相似度搜索
	// vector: 查询向量
	// topK: 返回的最相似文档数量
	// 返回: 相似文档列表和错误
	Search(ctx context.Context, vector []float32, topK int) ([]*entity.Document, error)

	// Insert 插入文档向量
	// docs: 要插入的文档列表
	// 返回: 错误
	Insert(ctx context.Context, docs []*entity.Document) error

	// Delete 删除文档向量
	// ids: 要删除的文档 ID 列表
	// 返回: 删除的文档数量和错误
	Delete(ctx context.Context, ids []string) (int, error)

	// GetByID 根据 ID 获取文档
	// id: 文档 ID
	// 返回: 文档和错误
	GetByID(ctx context.Context, id string) (*entity.Document, error)

	// Count 获取文档总数
	// 返回: 文档数量和错误
	Count(ctx context.Context) (int64, error)

	// CreateCollection 创建向量集合
	// collectionName: 集合名称
	// dimension: 向量维度
	// 返回: 错误
	CreateCollection(ctx context.Context, collectionName string, dimension int) error

	// CollectionExists 检查集合是否存在
	// collectionName: 集合名称
	// 返回: 是否存在和错误
	CollectionExists(ctx context.Context, collectionName string) (bool, error)

	// DropCollection 删除向量集合
	// collectionName: 集合名称
	// 返回: 错误
	DropCollection(ctx context.Context, collectionName string) error
}
