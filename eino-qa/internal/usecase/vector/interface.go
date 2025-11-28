package vector

import (
	"context"
	"eino-qa/internal/domain/entity"
)

// VectorUseCaseInterface 向量管理用例接口
type VectorUseCaseInterface interface {
	AddVectors(ctx context.Context, req *AddVectorRequest) (*AddVectorResponse, error)
	DeleteVectors(ctx context.Context, req *DeleteVectorRequest) (*DeleteVectorResponse, error)
	GetVectorCount(ctx context.Context, tenantID string) (int64, error)
	GetVectorByID(ctx context.Context, id string, tenantID string) (*entity.Document, error)
}
