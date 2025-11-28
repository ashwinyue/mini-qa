package vector

import (
	"context"
	"testing"

	"eino-qa/internal/domain/entity"

	"github.com/cloudwego/eino/components/embedding"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockEmbedder 模拟嵌入模型
type MockEmbedder struct {
	mock.Mock
}

func (m *MockEmbedder) EmbedStrings(ctx context.Context, texts []string, opts ...embedding.Option) ([][]float64, error) {
	args := m.Called(ctx, texts)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([][]float64), args.Error(1)
}

// MockVectorRepository 模拟向量仓储
type MockVectorRepository struct {
	mock.Mock
}

func (m *MockVectorRepository) Search(ctx context.Context, vector []float32, topK int) ([]*entity.Document, error) {
	args := m.Called(ctx, vector, topK)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.Document), args.Error(1)
}

func (m *MockVectorRepository) Insert(ctx context.Context, docs []*entity.Document) error {
	args := m.Called(ctx, docs)
	return args.Error(0)
}

func (m *MockVectorRepository) Delete(ctx context.Context, ids []string) (int, error) {
	args := m.Called(ctx, ids)
	return args.Int(0), args.Error(1)
}

func (m *MockVectorRepository) GetByID(ctx context.Context, id string) (*entity.Document, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Document), args.Error(1)
}

func (m *MockVectorRepository) Count(ctx context.Context) (int64, error) {
	args := m.Called(ctx)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockVectorRepository) CreateCollection(ctx context.Context, collectionName string, dimension int) error {
	args := m.Called(ctx, collectionName, dimension)
	return args.Error(0)
}

func (m *MockVectorRepository) CollectionExists(ctx context.Context, collectionName string) (bool, error) {
	args := m.Called(ctx, collectionName)
	return args.Bool(0), args.Error(1)
}

func (m *MockVectorRepository) DropCollection(ctx context.Context, collectionName string) error {
	args := m.Called(ctx, collectionName)
	return args.Error(0)
}

// TestAddVectors 测试添加向量
func TestAddVectors(t *testing.T) {
	mockEmbedder := new(MockEmbedder)
	mockVectorRepo := new(MockVectorRepository)

	uc := NewVectorManagementUseCase(mockEmbedder, mockVectorRepo, nil)

	ctx := context.Background()
	texts := []string{"Python 课程介绍", "Go 语言基础"}

	// 模拟嵌入模型返回
	mockEmbeddings := [][]float64{
		{0.1, 0.2, 0.3},
		{0.4, 0.5, 0.6},
	}
	mockEmbedder.On("EmbedStrings", mock.Anything, texts).Return(mockEmbeddings, nil)

	// 模拟向量仓储插入
	mockVectorRepo.On("Insert", mock.Anything, mock.MatchedBy(func(docs []*entity.Document) bool {
		return len(docs) == 2
	})).Return(nil)

	req := &AddVectorRequest{
		Texts:    texts,
		TenantID: "test",
	}

	resp, err := uc.AddVectors(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.True(t, resp.Success)
	assert.Equal(t, 2, resp.Count)
	assert.Len(t, resp.DocumentIDs, 2)

	mockEmbedder.AssertExpectations(t)
	mockVectorRepo.AssertExpectations(t)
}

// TestAddVectors_EmptyTexts 测试空文本列表
func TestAddVectors_EmptyTexts(t *testing.T) {
	mockEmbedder := new(MockEmbedder)
	mockVectorRepo := new(MockVectorRepository)

	uc := NewVectorManagementUseCase(mockEmbedder, mockVectorRepo, nil)

	ctx := context.Background()
	req := &AddVectorRequest{
		Texts:    []string{},
		TenantID: "test",
	}

	resp, err := uc.AddVectors(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "texts cannot be empty")
}

// TestAddVectors_WithMetadata 测试带元数据的添加
func TestAddVectors_WithMetadata(t *testing.T) {
	mockEmbedder := new(MockEmbedder)
	mockVectorRepo := new(MockVectorRepository)

	uc := NewVectorManagementUseCase(mockEmbedder, mockVectorRepo, nil)

	ctx := context.Background()
	texts := []string{"Python 课程介绍"}

	mockEmbeddings := [][]float64{
		{0.1, 0.2, 0.3},
	}
	mockEmbedder.On("EmbedStrings", mock.Anything, texts).Return(mockEmbeddings, nil)

	mockVectorRepo.On("Insert", mock.Anything, mock.MatchedBy(func(docs []*entity.Document) bool {
		if len(docs) != 1 {
			return false
		}
		// 验证元数据
		category, exists := docs[0].GetMetadata("category")
		return exists && category == "course"
	})).Return(nil)

	req := &AddVectorRequest{
		Texts:    texts,
		TenantID: "test",
		Metadata: map[string]any{
			"category": "course",
		},
	}

	resp, err := uc.AddVectors(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.True(t, resp.Success)

	mockEmbedder.AssertExpectations(t)
	mockVectorRepo.AssertExpectations(t)
}

// TestDeleteVectors 测试删除向量
func TestDeleteVectors(t *testing.T) {
	mockEmbedder := new(MockEmbedder)
	mockVectorRepo := new(MockVectorRepository)

	uc := NewVectorManagementUseCase(mockEmbedder, mockVectorRepo, nil)

	ctx := context.Background()
	ids := []string{"doc_001", "doc_002"}

	// 模拟删除操作
	mockVectorRepo.On("Delete", mock.Anything, ids).Return(2, nil)

	req := &DeleteVectorRequest{
		IDs:      ids,
		TenantID: "test",
	}

	resp, err := uc.DeleteVectors(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.True(t, resp.Success)
	assert.Equal(t, 2, resp.DeletedCount)

	mockVectorRepo.AssertExpectations(t)
}

// TestDeleteVectors_EmptyIDs 测试空 ID 列表
func TestDeleteVectors_EmptyIDs(t *testing.T) {
	mockEmbedder := new(MockEmbedder)
	mockVectorRepo := new(MockVectorRepository)

	uc := NewVectorManagementUseCase(mockEmbedder, mockVectorRepo, nil)

	ctx := context.Background()
	req := &DeleteVectorRequest{
		IDs:      []string{},
		TenantID: "test",
	}

	resp, err := uc.DeleteVectors(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "ids cannot be empty")
}

// TestGetVectorCount 测试获取向量数量
func TestGetVectorCount(t *testing.T) {
	mockEmbedder := new(MockEmbedder)
	mockVectorRepo := new(MockVectorRepository)

	uc := NewVectorManagementUseCase(mockEmbedder, mockVectorRepo, nil)

	ctx := context.Background()

	mockVectorRepo.On("Count", mock.Anything).Return(int64(100), nil)

	count, err := uc.GetVectorCount(ctx, "test")

	assert.NoError(t, err)
	assert.Equal(t, int64(100), count)

	mockVectorRepo.AssertExpectations(t)
}

// TestGetVectorByID 测试根据 ID 获取向量
func TestGetVectorByID(t *testing.T) {
	mockEmbedder := new(MockEmbedder)
	mockVectorRepo := new(MockVectorRepository)

	uc := NewVectorManagementUseCase(mockEmbedder, mockVectorRepo, nil)

	ctx := context.Background()
	docID := "doc_001"

	expectedDoc := &entity.Document{
		ID:       docID,
		Content:  "Python 课程介绍",
		TenantID: "test",
	}

	mockVectorRepo.On("GetByID", mock.Anything, docID).Return(expectedDoc, nil)

	doc, err := uc.GetVectorByID(ctx, docID, "test")

	assert.NoError(t, err)
	assert.NotNil(t, doc)
	assert.Equal(t, docID, doc.ID)
	assert.Equal(t, "Python 课程介绍", doc.Content)

	mockVectorRepo.AssertExpectations(t)
}

// TestGetVectorByID_EmptyID 测试空 ID
func TestGetVectorByID_EmptyID(t *testing.T) {
	mockEmbedder := new(MockEmbedder)
	mockVectorRepo := new(MockVectorRepository)

	uc := NewVectorManagementUseCase(mockEmbedder, mockVectorRepo, nil)

	ctx := context.Background()

	doc, err := uc.GetVectorByID(ctx, "", "test")

	assert.Error(t, err)
	assert.Nil(t, doc)
	assert.Contains(t, err.Error(), "id cannot be empty")
}

// TestAddVectors_BatchOperation 测试批量添加
func TestAddVectors_BatchOperation(t *testing.T) {
	mockEmbedder := new(MockEmbedder)
	mockVectorRepo := new(MockVectorRepository)

	uc := NewVectorManagementUseCase(mockEmbedder, mockVectorRepo, nil)

	ctx := context.Background()

	// 批量添加 10 个文档
	texts := make([]string, 10)
	mockEmbeddings := make([][]float64, 10)
	for i := 0; i < 10; i++ {
		texts[i] = "Document " + string(rune('A'+i))
		mockEmbeddings[i] = []float64{float64(i) * 0.1, float64(i) * 0.2, float64(i) * 0.3}
	}

	mockEmbedder.On("EmbedStrings", mock.Anything, texts).Return(mockEmbeddings, nil)
	mockVectorRepo.On("Insert", mock.Anything, mock.MatchedBy(func(docs []*entity.Document) bool {
		return len(docs) == 10
	})).Return(nil)

	req := &AddVectorRequest{
		Texts:    texts,
		TenantID: "test",
	}

	resp, err := uc.AddVectors(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.True(t, resp.Success)
	assert.Equal(t, 10, resp.Count)
	assert.Len(t, resp.DocumentIDs, 10)

	mockEmbedder.AssertExpectations(t)
	mockVectorRepo.AssertExpectations(t)
}

// TestDeleteVectors_BatchOperation 测试批量删除
func TestDeleteVectors_BatchOperation(t *testing.T) {
	mockEmbedder := new(MockEmbedder)
	mockVectorRepo := new(MockVectorRepository)

	uc := NewVectorManagementUseCase(mockEmbedder, mockVectorRepo, nil)

	ctx := context.Background()

	// 批量删除 10 个文档
	ids := make([]string, 10)
	for i := 0; i < 10; i++ {
		ids[i] = "doc_" + string(rune('0'+i))
	}

	mockVectorRepo.On("Delete", mock.Anything, ids).Return(10, nil)

	req := &DeleteVectorRequest{
		IDs:      ids,
		TenantID: "test",
	}

	resp, err := uc.DeleteVectors(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.True(t, resp.Success)
	assert.Equal(t, 10, resp.DeletedCount)

	mockVectorRepo.AssertExpectations(t)
}
