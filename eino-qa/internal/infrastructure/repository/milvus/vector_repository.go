package milvus

import (
	"context"
	"encoding/json"
	"fmt"

	"eino-qa/internal/domain/entity"
	"eino-qa/internal/domain/repository"

	milvusEntity "github.com/milvus-io/milvus-sdk-go/v2/entity"
	"github.com/sirupsen/logrus"
)

// VectorRepository Milvus 向量仓储实现
type VectorRepository struct {
	client        *Client
	tenantManager *TenantManager
	logger        *logrus.Logger
}

// NewVectorRepository 创建向量仓储
func NewVectorRepository(client *Client, tenantManager *TenantManager, logger *logrus.Logger) repository.VectorRepository {
	if logger == nil {
		logger = logrus.New()
	}

	return &VectorRepository{
		client:        client,
		tenantManager: tenantManager,
		logger:        logger,
	}
}

// Search 执行向量相似度搜索
func (r *VectorRepository) Search(ctx context.Context, vector []float32, topK int) ([]*entity.Document, error) {
	// 从上下文获取租户 ID
	tenantID, ok := ctx.Value("tenant_id").(string)
	if !ok || tenantID == "" {
		tenantID = "default"
	}

	// 获取租户的 Collection
	collectionName, err := r.tenantManager.GetCollection(ctx, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get collection for tenant %s: %w", tenantID, err)
	}

	r.logger.WithFields(logrus.Fields{
		"tenant_id":  tenantID,
		"collection": collectionName,
		"top_k":      topK,
	}).Debug("searching vectors")

	// 构建搜索向量
	searchVectors := []milvusEntity.Vector{
		milvusEntity.FloatVector(vector),
	}

	// 执行搜索
	sp, _ := milvusEntity.NewIndexHNSWSearchParam(16)
	searchResult, err := r.client.GetClient().Search(
		ctx,
		collectionName,
		nil, // partitions
		"",  // expr
		[]string{"id", "content", "metadata", "tenant_id", "created_at"},
		searchVectors,
		"vector",
		milvusEntity.L2,
		topK,
		sp,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to search vectors: %w", err)
	}

	if len(searchResult) == 0 {
		return []*entity.Document{}, nil
	}

	// 转换结果
	results := searchResult[0]
	documents := make([]*entity.Document, 0, results.ResultCount)

	for i := 0; i < results.ResultCount; i++ {
		doc := &entity.Document{}

		// ID
		if idField := results.Fields.GetColumn("id"); idField != nil {
			if idData, ok := idField.(*milvusEntity.ColumnVarChar); ok {
				doc.ID = idData.Data()[i]
			}
		}

		// Content
		if contentField := results.Fields.GetColumn("content"); contentField != nil {
			if contentData, ok := contentField.(*milvusEntity.ColumnVarChar); ok {
				doc.Content = contentData.Data()[i]
			}
		}

		// Metadata
		if metadataField := results.Fields.GetColumn("metadata"); metadataField != nil {
			if metadataData, ok := metadataField.(*milvusEntity.ColumnJSONBytes); ok {
				var metadata map[string]any
				if err := json.Unmarshal(metadataData.Data()[i], &metadata); err == nil {
					doc.Metadata = metadata
				}
			}
		}

		// TenantID
		if tenantField := results.Fields.GetColumn("tenant_id"); tenantField != nil {
			if tenantData, ok := tenantField.(*milvusEntity.ColumnVarChar); ok {
				doc.TenantID = tenantData.Data()[i]
			}
		}

		// Score
		doc.Score = float64(results.Scores[i])

		documents = append(documents, doc)
	}

	r.logger.WithFields(logrus.Fields{
		"tenant_id":  tenantID,
		"collection": collectionName,
		"found":      len(documents),
	}).Debug("search completed")

	return documents, nil
}

// Insert 插入文档向量
func (r *VectorRepository) Insert(ctx context.Context, docs []*entity.Document) error {
	if len(docs) == 0 {
		return nil
	}

	// 从上下文获取租户 ID
	tenantID, ok := ctx.Value("tenant_id").(string)
	if !ok || tenantID == "" {
		tenantID = "default"
	}

	// 获取租户的 Collection
	collectionName, err := r.tenantManager.GetCollection(ctx, tenantID)
	if err != nil {
		return fmt.Errorf("failed to get collection for tenant %s: %w", tenantID, err)
	}

	r.logger.WithFields(logrus.Fields{
		"tenant_id":  tenantID,
		"collection": collectionName,
		"count":      len(docs),
	}).Info("inserting documents")

	// 准备数据列
	ids := make([]string, len(docs))
	vectors := make([][]float32, len(docs))
	contents := make([]string, len(docs))
	metadatas := make([][]byte, len(docs))
	tenantIDs := make([]string, len(docs))
	createdAts := make([]int64, len(docs))

	for i, doc := range docs {
		ids[i] = doc.ID
		vectors[i] = doc.Vector
		contents[i] = doc.Content
		tenantIDs[i] = tenantID

		// 序列化 metadata
		metadataBytes, err := json.Marshal(doc.Metadata)
		if err != nil {
			return fmt.Errorf("failed to marshal metadata for doc %s: %w", doc.ID, err)
		}
		metadatas[i] = metadataBytes

		createdAts[i] = doc.CreatedAt.Unix()
	}

	// 构建列数据
	idColumn := milvusEntity.NewColumnVarChar("id", ids)
	vectorColumn := milvusEntity.NewColumnFloatVector("vector", len(vectors[0]), vectors)
	contentColumn := milvusEntity.NewColumnVarChar("content", contents)
	metadataColumn := milvusEntity.NewColumnJSONBytes("metadata", metadatas)
	tenantIDColumn := milvusEntity.NewColumnVarChar("tenant_id", tenantIDs)
	createdAtColumn := milvusEntity.NewColumnInt64("created_at", createdAts)

	// 插入数据
	_, err = r.client.GetClient().Insert(
		ctx,
		collectionName,
		"", // partition
		idColumn,
		vectorColumn,
		contentColumn,
		metadataColumn,
		tenantIDColumn,
		createdAtColumn,
	)
	if err != nil {
		return fmt.Errorf("failed to insert documents: %w", err)
	}

	// 刷新以确保数据持久化
	err = r.client.GetClient().Flush(ctx, collectionName, false)
	if err != nil {
		r.logger.WithError(err).Warn("failed to flush collection")
	}

	r.logger.WithFields(logrus.Fields{
		"tenant_id":  tenantID,
		"collection": collectionName,
		"count":      len(docs),
	}).Info("documents inserted successfully")

	return nil
}

// Delete 删除文档向量
func (r *VectorRepository) Delete(ctx context.Context, ids []string) (int, error) {
	if len(ids) == 0 {
		return 0, nil
	}

	// 从上下文获取租户 ID
	tenantID, ok := ctx.Value("tenant_id").(string)
	if !ok || tenantID == "" {
		tenantID = "default"
	}

	// 获取租户的 Collection
	collectionName, err := r.tenantManager.GetCollection(ctx, tenantID)
	if err != nil {
		return 0, fmt.Errorf("failed to get collection for tenant %s: %w", tenantID, err)
	}

	r.logger.WithFields(logrus.Fields{
		"tenant_id":  tenantID,
		"collection": collectionName,
		"count":      len(ids),
	}).Info("deleting documents")

	// 构建删除表达式
	expr := fmt.Sprintf("id in [%s]", r.buildIDList(ids))

	// 执行删除
	err = r.client.GetClient().Delete(ctx, collectionName, "", expr)
	if err != nil {
		return 0, fmt.Errorf("failed to delete documents: %w", err)
	}

	// 刷新以确保删除生效
	err = r.client.GetClient().Flush(ctx, collectionName, false)
	if err != nil {
		r.logger.WithError(err).Warn("failed to flush collection after delete")
	}

	r.logger.WithFields(logrus.Fields{
		"tenant_id":  tenantID,
		"collection": collectionName,
		"count":      len(ids),
	}).Info("documents deleted successfully")

	return len(ids), nil
}

// GetByID 根据 ID 获取文档
func (r *VectorRepository) GetByID(ctx context.Context, id string) (*entity.Document, error) {
	// 从上下文获取租户 ID
	tenantID, ok := ctx.Value("tenant_id").(string)
	if !ok || tenantID == "" {
		tenantID = "default"
	}

	// 获取租户的 Collection
	collectionName, err := r.tenantManager.GetCollection(ctx, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get collection for tenant %s: %w", tenantID, err)
	}

	// 构建查询表达式
	expr := fmt.Sprintf("id == \"%s\"", id)

	// 执行查询
	queryResult, err := r.client.GetClient().Query(
		ctx,
		collectionName,
		nil, // partitions
		expr,
		[]string{"id", "content", "metadata", "tenant_id", "created_at"},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query document: %w", err)
	}

	if queryResult == nil || len(queryResult) == 0 {
		return nil, fmt.Errorf("document not found: %s", id)
	}

	// 转换结果
	doc := &entity.Document{}

	if idField := queryResult.GetColumn("id"); idField != nil {
		if idData, ok := idField.(*milvusEntity.ColumnVarChar); ok && len(idData.Data()) > 0 {
			doc.ID = idData.Data()[0]
		}
	}

	if contentField := queryResult.GetColumn("content"); contentField != nil {
		if contentData, ok := contentField.(*milvusEntity.ColumnVarChar); ok && len(contentData.Data()) > 0 {
			doc.Content = contentData.Data()[0]
		}
	}

	if metadataField := queryResult.GetColumn("metadata"); metadataField != nil {
		if metadataData, ok := metadataField.(*milvusEntity.ColumnJSONBytes); ok && len(metadataData.Data()) > 0 {
			var metadata map[string]any
			if err := json.Unmarshal(metadataData.Data()[0], &metadata); err == nil {
				doc.Metadata = metadata
			}
		}
	}

	if tenantField := queryResult.GetColumn("tenant_id"); tenantField != nil {
		if tenantData, ok := tenantField.(*milvusEntity.ColumnVarChar); ok && len(tenantData.Data()) > 0 {
			doc.TenantID = tenantData.Data()[0]
		}
	}

	return doc, nil
}

// Count 获取文档总数
func (r *VectorRepository) Count(ctx context.Context) (int64, error) {
	// 从上下文获取租户 ID
	tenantID, ok := ctx.Value("tenant_id").(string)
	if !ok || tenantID == "" {
		tenantID = "default"
	}

	// 获取租户的 Collection
	collectionName, err := r.tenantManager.GetCollection(ctx, tenantID)
	if err != nil {
		return 0, fmt.Errorf("failed to get collection for tenant %s: %w", tenantID, err)
	}

	// 获取统计信息
	stats, err := r.client.GetClient().GetCollectionStatistics(ctx, collectionName)
	if err != nil {
		return 0, fmt.Errorf("failed to get collection statistics: %w", err)
	}

	// 解析行数
	var count int64
	if rowCountStr, ok := stats["row_count"]; ok {
		fmt.Sscanf(rowCountStr, "%d", &count)
	}

	return count, nil
}

// CreateCollection 创建向量集合
func (r *VectorRepository) CreateCollection(ctx context.Context, collectionName string, dimension int) error {
	collectionManager := NewCollectionManager(r.client, r.logger)
	return collectionManager.CreateCollection(ctx, collectionName, dimension)
}

// CollectionExists 检查集合是否存在
func (r *VectorRepository) CollectionExists(ctx context.Context, collectionName string) (bool, error) {
	collectionManager := NewCollectionManager(r.client, r.logger)
	return collectionManager.CollectionExists(ctx, collectionName)
}

// DropCollection 删除向量集合
func (r *VectorRepository) DropCollection(ctx context.Context, collectionName string) error {
	collectionManager := NewCollectionManager(r.client, r.logger)
	return collectionManager.DropCollection(ctx, collectionName)
}

// buildIDList 构建 ID 列表字符串
func (r *VectorRepository) buildIDList(ids []string) string {
	result := ""
	for i, id := range ids {
		if i > 0 {
			result += ", "
		}
		result += fmt.Sprintf("\"%s\"", id)
	}
	return result
}
