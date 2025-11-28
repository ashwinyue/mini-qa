package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"eino-qa/internal/domain/entity"
	"eino-qa/internal/usecase/chat"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockChatUseCase 模拟对话用例
type MockChatUseCase struct {
	mock.Mock
}

func (m *MockChatUseCase) Execute(ctx context.Context, req *chat.ChatRequest) (*chat.ChatResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*chat.ChatResponse), args.Error(1)
}

func (m *MockChatUseCase) ExecuteStream(ctx context.Context, req *chat.ChatRequest) (<-chan *chat.StreamChunk, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(<-chan *chat.StreamChunk), args.Error(1)
}

func TestChatHandler_HandleChat_Success(t *testing.T) {
	// 设置 Gin 为测试模式
	gin.SetMode(gin.TestMode)

	// 创建模拟用例
	mockUseCase := new(MockChatUseCase)
	handler := NewChatHandler(mockUseCase)

	// 准备测试数据
	requestBody := ChatRequestDTO{
		Query:     "Python课程包含哪些内容？",
		TenantID:  "tenant1",
		SessionID: "session123",
		Stream:    false,
	}

	expectedResponse := &chat.ChatResponse{
		Answer:    "Python课程包含基础语法、数据结构等内容",
		Route:     "course",
		SessionID: "session123",
		Sources: []*entity.Document{
			{
				ID:      "doc1",
				Content: "Python基础课程",
				Score:   0.95,
			},
		},
		Metadata: map[string]any{
			"intent":     "course",
			"confidence": 0.98,
		},
	}

	// 设置模拟期望
	mockUseCase.On("Execute", mock.Anything, mock.MatchedBy(func(req *chat.ChatRequest) bool {
		return req.Query == "Python课程包含哪些内容？" &&
			req.TenantID == "tenant1" &&
			req.SessionID == "session123"
	})).Return(expectedResponse, nil)

	// 创建请求
	body, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPost, "/chat", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	// 创建响应记录器
	w := httptest.NewRecorder()

	// 创建 Gin 上下文
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Set("tenant_id", "tenant1")

	// 执行处理器
	handler.HandleChat(c)

	// 验证响应
	assert.Equal(t, http.StatusOK, w.Code)

	var response ChatResponseDTO
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Python课程包含基础语法、数据结构等内容", response.Answer)
	assert.Equal(t, "course", response.Route)
	assert.Equal(t, "session123", response.SessionID)
	assert.Len(t, response.Sources, 1)
	assert.Equal(t, "Python基础课程", response.Sources[0].Content)

	// 验证模拟调用
	mockUseCase.AssertExpectations(t)
}

func TestChatHandler_HandleChat_InvalidRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockUseCase := new(MockChatUseCase)
	handler := NewChatHandler(mockUseCase)

	// 创建无效请求（缺少 query）
	requestBody := ChatRequestDTO{
		TenantID:  "tenant1",
		SessionID: "session123",
	}

	body, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPost, "/chat", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handler.HandleChat(c)

	// 验证有错误记录
	assert.NotEmpty(t, c.Errors)
}

func TestChatHandler_HandleChat_DefaultTenantID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockUseCase := new(MockChatUseCase)
	handler := NewChatHandler(mockUseCase)

	requestBody := ChatRequestDTO{
		Query:     "Hello",
		SessionID: "session123",
		// 不提供 TenantID
	}

	expectedResponse := &chat.ChatResponse{
		Answer:    "Hello!",
		Route:     "direct",
		SessionID: "session123",
		Metadata:  map[string]any{},
	}

	// 验证使用默认租户 ID
	mockUseCase.On("Execute", mock.Anything, mock.MatchedBy(func(req *chat.ChatRequest) bool {
		return req.TenantID == "default"
	})).Return(expectedResponse, nil)

	body, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPost, "/chat", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handler.HandleChat(c)

	assert.Equal(t, http.StatusOK, w.Code)
	mockUseCase.AssertExpectations(t)
}

func TestChatHandler_toChatResponseDTO(t *testing.T) {
	handler := NewChatHandler(nil)

	// 测试转换
	resp := &chat.ChatResponse{
		Answer:    "Test answer",
		Route:     "course",
		SessionID: "session123",
		Sources: []*entity.Document{
			{
				ID:      "doc1",
				Content: "Source 1",
				Score:   0.95,
				Metadata: map[string]any{
					"category": "test",
				},
			},
		},
		Metadata: map[string]any{
			"intent": "course",
		},
	}

	dto := handler.toChatResponseDTO(resp)

	assert.Equal(t, "Test answer", dto.Answer)
	assert.Equal(t, "course", dto.Route)
	assert.Equal(t, "session123", dto.SessionID)
	assert.Len(t, dto.Sources, 1)
	assert.Equal(t, "Source 1", dto.Sources[0].Content)
	assert.Equal(t, 0.95, dto.Sources[0].Score)
	assert.Equal(t, "test", dto.Sources[0].Metadata["category"])
}

func TestChatHandler_toChatResponseDTO_NoSources(t *testing.T) {
	handler := NewChatHandler(nil)

	resp := &chat.ChatResponse{
		Answer:    "Test answer",
		Route:     "direct",
		SessionID: "session123",
		Sources:   nil,
		Metadata:  map[string]any{},
	}

	dto := handler.toChatResponseDTO(resp)

	assert.Equal(t, "Test answer", dto.Answer)
	assert.Nil(t, dto.Sources)
}
