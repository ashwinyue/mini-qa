package chat

import (
	"context"
	"testing"
	"time"

	"eino-qa/internal/domain/entity"
	"eino-qa/internal/infrastructure/logger"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockSessionRepository 模拟会话仓储
type MockSessionRepository struct {
	mock.Mock
}

func (m *MockSessionRepository) Save(ctx context.Context, session *entity.Session) error {
	args := m.Called(ctx, session)
	return args.Error(0)
}

func (m *MockSessionRepository) Load(ctx context.Context, sessionID string) (*entity.Session, error) {
	args := m.Called(ctx, sessionID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Session), args.Error(1)
}

func (m *MockSessionRepository) Delete(ctx context.Context, sessionID string) error {
	args := m.Called(ctx, sessionID)
	return args.Error(0)
}

func (m *MockSessionRepository) Exists(ctx context.Context, sessionID string) (bool, error) {
	args := m.Called(ctx, sessionID)
	return args.Bool(0), args.Error(1)
}

func (m *MockSessionRepository) AddMessage(ctx context.Context, sessionID string, message *entity.Message) error {
	args := m.Called(ctx, sessionID, message)
	return args.Error(0)
}

func (m *MockSessionRepository) GetMessages(ctx context.Context, sessionID string) ([]*entity.Message, error) {
	args := m.Called(ctx, sessionID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.Message), args.Error(1)
}

func (m *MockSessionRepository) UpdateExpiration(ctx context.Context, sessionID string, expiresAt time.Time) error {
	args := m.Called(ctx, sessionID, expiresAt)
	return args.Error(0)
}

func (m *MockSessionRepository) DeleteExpired(ctx context.Context) (int, error) {
	args := m.Called(ctx)
	return args.Int(0), args.Error(1)
}

func (m *MockSessionRepository) ListByTenant(ctx context.Context, tenantID string) ([]*entity.Session, error) {
	args := m.Called(ctx, tenantID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.Session), args.Error(1)
}

func (m *MockSessionRepository) Count(ctx context.Context) (int64, error) {
	args := m.Called(ctx)
	return args.Get(0).(int64), args.Error(1)
}

// TestChatRequest_Validate 测试请求验证
func TestChatRequest_Validate(t *testing.T) {
	tests := []struct {
		name    string
		req     *ChatRequest
		wantErr bool
	}{
		{
			name: "valid request",
			req: &ChatRequest{
				Query:    "Hello",
				TenantID: "tenant1",
			},
			wantErr: false,
		},
		{
			name: "empty query",
			req: &ChatRequest{
				Query:    "",
				TenantID: "tenant1",
			},
			wantErr: true,
		},
		{
			name: "empty tenant id",
			req: &ChatRequest{
				Query:    "Hello",
				TenantID: "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.req.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// TestChatUseCase_loadOrCreateSession 测试会话加载或创建
func TestChatUseCase_loadOrCreateSession(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockSessionRepository)
	log, _ := logger.New(logger.Config{
		Level:  "info",
		Format: "text",
		Output: "stdout",
	})

	uc := &ChatUseCase{
		sessionRepo: mockRepo,
		sessionTTL:  30 * time.Minute,
		logger:      log,
	}

	t.Run("create new session when session id is empty", func(t *testing.T) {
		session, err := uc.loadOrCreateSession(ctx, "tenant1", "")
		assert.NoError(t, err)
		assert.NotNil(t, session)
		assert.Equal(t, "tenant1", session.TenantID)
		assert.False(t, session.IsExpired())
	})

	t.Run("load existing session", func(t *testing.T) {
		existingSession := entity.NewSession("tenant1", 30*time.Minute)
		existingSession.ID = "sess_123"

		mockRepo.On("Load", ctx, "sess_123").Return(existingSession, nil).Once()

		session, err := uc.loadOrCreateSession(ctx, "tenant1", "sess_123")
		assert.NoError(t, err)
		assert.NotNil(t, session)
		assert.Equal(t, "sess_123", session.ID)
		assert.Equal(t, "tenant1", session.TenantID)

		mockRepo.AssertExpectations(t)
	})

	t.Run("create new session when load fails", func(t *testing.T) {
		mockRepo.On("Load", ctx, "sess_invalid").Return(nil, entity.ErrSessionExpired).Once()

		session, err := uc.loadOrCreateSession(ctx, "tenant1", "sess_invalid")
		assert.NoError(t, err)
		assert.NotNil(t, session)
		assert.Equal(t, "tenant1", session.TenantID)
		assert.NotEqual(t, "sess_invalid", session.ID) // 应该是新的会话 ID

		mockRepo.AssertExpectations(t)
	})
}

// TestNewChatUseCase 测试创建 ChatUseCase
func TestNewChatUseCase(t *testing.T) {
	mockRepo := new(MockSessionRepository)
	log, _ := logger.New(logger.Config{
		Level:  "info",
		Format: "text",
		Output: "stdout",
	})

	t.Run("with default session ttl", func(t *testing.T) {
		uc := NewChatUseCase(nil, nil, nil, nil, mockRepo, 0, log)
		assert.NotNil(t, uc)
		assert.Equal(t, 30*time.Minute, uc.sessionTTL)
	})

	t.Run("with custom session ttl", func(t *testing.T) {
		customTTL := 1 * time.Hour
		uc := NewChatUseCase(nil, nil, nil, nil, mockRepo, customTTL, log)
		assert.NotNil(t, uc)
		assert.Equal(t, customTTL, uc.sessionTTL)
	})
}

// 注意：完整的集成测试需要实际的 AI 组件和数据库连接
// 这里只提供了基本的单元测试示例
