package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"eino-qa/internal/domain/entity"
	"eino-qa/internal/usecase/vector"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockVectorUseCase 模拟向量管理用例
type MockVectorUseCase struct {
	mock.Mock
}

func (m *MockVectorUseCase) AddVectors(ctx context.Context, req *vector.AddVectorRequest) (*vector.AddVectorResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*vector.AddVectorResponse), args.Error(1)
}

func (m *MockVectorUseCase) DeleteVectors(ctx context.Context, req *vector.DeleteVectorRequest) (*vector.DeleteVectorResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*vector.DeleteVectorResponse), args.Error(1)
}

func (m *MockVectorUseCase) GetVectorCount(ctx context.Context, tenantID string) (int64, error) {
	args := m.Called(ctx, tenantID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockVectorUseCase) GetVectorByID(ctx context.Context, id string, tenantID string) (*entity.Document, error) {
	args := m.Called(ctx, id, tenantID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Document), args.Error(1)
}

func TestVectorHandler_HandleAddVectors_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockUseCase := new(MockVectorUseCase)
	handler := NewVectorHandler(mockUseCase)

	requestBody := AddVectorRequestDTO{
		Texts:    []string{"Python是一门编程语言", "Go语言适合高并发"},
		TenantID: "tenant1",
		Metadata: map[string]any{"category": "programming"},
	}

	expectedResponse := &vector.AddVectorResponse{
		Success:     true,
		DocumentIDs: []string{"doc1", "doc2"},
		Count:       2,
		Message:     "successfully added 2 vectors",
	}

	mockUseCase.On("AddVectors", mock.Anything, mock.MatchedBy(func(req *vector.AddVectorRequest) bool {
		return len(req.Texts) == 2 && req.TenantID == "tenant1"
	})).Return(expectedResponse, nil)

	body, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/vectors/items", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Set("tenant_id", "tenant1")

	handler.HandleAddVectors(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var response AddVectorResponseDTO
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Success)
	assert.Equal(t, 2, response.Count)
	assert.Len(t, response.DocumentIDs, 2)

	mockUseCase.AssertExpectations(t)
}

func TestVectorHandler_HandleAddVectors_EmptyTexts(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockUseCase := new(MockVectorUseCase)
	handler := NewVectorHandler(mockUseCase)

	requestBody := AddVectorRequestDTO{
		Texts:    []string{},
		TenantID: "tenant1",
	}

	body, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/vectors/items", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handler.HandleAddVectors(c)

	// 验证有错误记录
	assert.NotEmpty(t, c.Errors)
}

func TestVectorHandler_HandleDeleteVectors_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockUseCase := new(MockVectorUseCase)
	handler := NewVectorHandler(mockUseCase)

	requestBody := DeleteVectorRequestDTO{
		IDs:      []string{"doc1", "doc2"},
		TenantID: "tenant1",
	}

	expectedResponse := &vector.DeleteVectorResponse{
		Success:      true,
		DeletedCount: 2,
		Message:      "successfully deleted 2 vectors",
	}

	mockUseCase.On("DeleteVectors", mock.Anything, mock.MatchedBy(func(req *vector.DeleteVectorRequest) bool {
		return len(req.IDs) == 2 && req.TenantID == "tenant1"
	})).Return(expectedResponse, nil)

	body, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodDelete, "/api/v1/vectors/items", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Set("tenant_id", "tenant1")

	handler.HandleDeleteVectors(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var response DeleteVectorResponseDTO
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Success)
	assert.Equal(t, 2, response.DeletedCount)

	mockUseCase.AssertExpectations(t)
}

func TestVectorHandler_HandleGetVectorCount_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockUseCase := new(MockVectorUseCase)
	handler := NewVectorHandler(mockUseCase)

	mockUseCase.On("GetVectorCount", mock.Anything, "tenant1").Return(int64(100), nil)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/vectors/count", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Set("tenant_id", "tenant1")

	handler.HandleGetVectorCount(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response["success"].(bool))
	assert.Equal(t, float64(100), response["count"].(float64))

	mockUseCase.AssertExpectations(t)
}

func TestVectorHandler_HandleAddVectors_DefaultTenantID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockUseCase := new(MockVectorUseCase)
	handler := NewVectorHandler(mockUseCase)

	requestBody := AddVectorRequestDTO{
		Texts: []string{"Test text"},
		// 不提供 TenantID
	}

	expectedResponse := &vector.AddVectorResponse{
		Success:     true,
		DocumentIDs: []string{"doc1"},
		Count:       1,
		Message:     "successfully added 1 vectors",
	}

	// 验证使用默认租户 ID
	mockUseCase.On("AddVectors", mock.Anything, mock.MatchedBy(func(req *vector.AddVectorRequest) bool {
		return req.TenantID == "default"
	})).Return(expectedResponse, nil)

	body, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/vectors/items", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handler.HandleAddVectors(c)

	assert.Equal(t, http.StatusOK, w.Code)
	mockUseCase.AssertExpectations(t)
}
