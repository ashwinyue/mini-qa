package vector

import (
	"context"
	"fmt"
	"time"

	"eino-qa/internal/domain/entity"
	"eino-qa/internal/domain/repository"

	"github.com/cloudwego/eino/components/embedding"
	"github.com/sirupsen/logrus"
)

// VectorManagementUseCase 向量管理用例
type VectorManagementUseCase struct {
	embedder   embedding.Embedder
	vectorRepo repository.VectorRepository
	logger     *logrus.Logger
}

// NewVectorManagementUseCase 创建向量管理用例
func NewVectorManagementUseCase(
	embedder embedding.Embedder,
	vectorRepo repository.VectorRepository,
	logger *logrus.Logger,
) *VectorManagementUseCase {
	if logger == nil {
		logger = logrus.New()
	}

	return &VectorManagementUseCase{
		embedder:   embedder,
		vectorRepo: vectorRepo,
		logger:     logger,
	}
}

// AddVectorRequest 添加向量请求
type AddVectorRequest struct {
	Texts    []string       `json:"texts" binding:"required"`
	TenantID string         `json:"tenant_id"`
	Metadata map[string]any `json:"metadata,omitempty"`
}

// AddVectorResponse 添加向量响应
type AddVectorResponse struct {
	Success     bool     `json:"success"`
	DocumentIDs []string `json:"document_ids"`
	Count       int      `json:"count"`
	Message     string   `json:"message"`
}

// DeleteVectorRequest 删除向量请求
type DeleteVectorRequest struct {
	IDs      []string `json:"ids" binding:"required"`
	TenantID string   `json:"tenant_id"`
}

// DeleteVectorResponse 删除向量响应
type DeleteVectorResponse struct {
	Success      bool   `json:"success"`
	DeletedCount int    `json:"deleted_count"`
	Message      string `json:"message"`
}

// AddVectors 添加向量
// 需求: 9.2, 9.3, 9.5
func (uc *VectorManagementUseCase) AddVectors(ctx context.Context, req *AddVectorRequest) (*AddVectorResponse, error) {
	if len(req.Texts) == 0 {
		return nil, fmt.Errorf("texts cannot be empty")
	}

	// 设置租户 ID
	tenantID := req.TenantID
	if tenantID == "" {
		tenantID = "default"
	}
	ctx = context.WithValue(ctx, "tenant_id", tenantID)

	uc.logger.WithFields(logrus.Fields{
		"tenant_id": tenantID,
		"count":     len(req.Texts),
	}).Info("adding vectors")

	// 1. 使用嵌入模型生成向量
	vectors, err := uc.generateVectors(ctx, req.Texts)
	if err != nil {
		uc.logger.WithError(err).Error("failed to generate vectors")
		return nil, fmt.Errorf("failed to generate vectors: %w", err)
	}

	// 2. 构建文档对象
	docs := make([]*entity.Document, len(req.Texts))
	documentIDs := make([]string, len(req.Texts))

	for i, text := range req.Texts {
		doc := entity.NewDocument(text, tenantID)
		doc.SetVector(vectors[i])

		// 添加元数据
		if req.Metadata != nil {
			for k, v := range req.Metadata {
				doc.AddMetadata(k, v)
			}
		}

		// 验证文档
		if err := doc.Validate(); err != nil {
			uc.logger.WithError(err).WithField("index", i).Error("document validation failed")
			return nil, fmt.Errorf("document validation failed at index %d: %w", i, err)
		}

		docs[i] = doc
		documentIDs[i] = doc.ID
	}

	// 3. 插入向量库
	err = uc.vectorRepo.Insert(ctx, docs)
	if err != nil {
		uc.logger.WithError(err).Error("failed to insert vectors")
		return nil, fmt.Errorf("failed to insert vectors: %w", err)
	}

	uc.logger.WithFields(logrus.Fields{
		"tenant_id": tenantID,
		"count":     len(docs),
	}).Info("vectors added successfully")

	return &AddVectorResponse{
		Success:     true,
		DocumentIDs: documentIDs,
		Count:       len(docs),
		Message:     fmt.Sprintf("successfully added %d vectors", len(docs)),
	}, nil
}

// DeleteVectors 删除向量
// 需求: 9.4, 9.5
func (uc *VectorManagementUseCase) DeleteVectors(ctx context.Context, req *DeleteVectorRequest) (*DeleteVectorResponse, error) {
	if len(req.IDs) == 0 {
		return nil, fmt.Errorf("ids cannot be empty")
	}

	// 设置租户 ID
	tenantID := req.TenantID
	if tenantID == "" {
		tenantID = "default"
	}
	ctx = context.WithValue(ctx, "tenant_id", tenantID)

	uc.logger.WithFields(logrus.Fields{
		"tenant_id": tenantID,
		"count":     len(req.IDs),
	}).Info("deleting vectors")

	// 执行删除
	deletedCount, err := uc.vectorRepo.Delete(ctx, req.IDs)
	if err != nil {
		uc.logger.WithError(err).Error("failed to delete vectors")
		return nil, fmt.Errorf("failed to delete vectors: %w", err)
	}

	uc.logger.WithFields(logrus.Fields{
		"tenant_id":     tenantID,
		"deleted_count": deletedCount,
	}).Info("vectors deleted successfully")

	return &DeleteVectorResponse{
		Success:      true,
		DeletedCount: deletedCount,
		Message:      fmt.Sprintf("successfully deleted %d vectors", deletedCount),
	}, nil
}

// generateVectors 生成文本向量
// 需求: 9.2
func (uc *VectorManagementUseCase) generateVectors(ctx context.Context, texts []string) ([][]float32, error) {
	startTime := time.Now()

	// 使用嵌入模型生成向量
	resp, err := uc.embedder.EmbedStrings(ctx, texts)
	if err != nil {
		return nil, fmt.Errorf("failed to embed texts: %w", err)
	}

	if len(resp) != len(texts) {
		return nil, fmt.Errorf("embedding result count mismatch: expected %d, got %d", len(texts), len(resp))
	}

	// 转换 float64 到 float32
	vectors := make([][]float32, len(resp))
	for i, embedding := range resp {
		if len(embedding) == 0 {
			return nil, fmt.Errorf("empty embedding result at index %d", i)
		}

		vector := make([]float32, len(embedding))
		for j, v := range embedding {
			vector[j] = float32(v)
		}
		vectors[i] = vector
	}

	duration := time.Since(startTime)
	uc.logger.WithFields(logrus.Fields{
		"count":    len(texts),
		"duration": duration.String(),
	}).Debug("vectors generated")

	return vectors, nil
}

// GetVectorCount 获取向量数量
func (uc *VectorManagementUseCase) GetVectorCount(ctx context.Context, tenantID string) (int64, error) {
	if tenantID == "" {
		tenantID = "default"
	}
	ctx = context.WithValue(ctx, "tenant_id", tenantID)

	count, err := uc.vectorRepo.Count(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to get vector count: %w", err)
	}

	return count, nil
}

// GetVectorByID 根据 ID 获取向量
func (uc *VectorManagementUseCase) GetVectorByID(ctx context.Context, id string, tenantID string) (*entity.Document, error) {
	if id == "" {
		return nil, fmt.Errorf("id cannot be empty")
	}

	if tenantID == "" {
		tenantID = "default"
	}
	ctx = context.WithValue(ctx, "tenant_id", tenantID)

	doc, err := uc.vectorRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get vector: %w", err)
	}

	return doc, nil
}
