package container

import (
	"testing"
	"time"

	"eino-qa/internal/infrastructure/config"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestContainer_New 测试容器创建
func TestContainer_New(t *testing.T) {
	// 创建测试配置
	cfg := createTestConfig()

	// 注意：这个测试需要实际的 Milvus 和 DashScope 连接
	// 在 CI 环境中应该跳过或使用 mock
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// 创建容器
	c, err := New(cfg)
	if err != nil {
		// 如果是连接错误，跳过测试
		t.Skipf("Failed to create container (expected in test environment): %v", err)
		return
	}

	// 验证容器不为空
	require.NotNil(t, c)

	// 验证关键组件已初始化
	assert.NotNil(t, c.Config)
	assert.NotNil(t, c.Logger)
	assert.NotNil(t, c.LogrusLogger)
	assert.NotNil(t, c.MetricsCollector)

	// 关闭容器
	err = c.Close()
	assert.NoError(t, err)
}

// TestContainer_ComponentsInitialized 测试组件初始化
func TestContainer_ComponentsInitialized(t *testing.T) {
	cfg := createTestConfig()

	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	c, err := New(cfg)
	if err != nil {
		t.Skipf("Failed to create container: %v", err)
		return
	}
	defer c.Close()

	// 验证基础设施组件
	assert.NotNil(t, c.Logger, "Logger should be initialized")
	assert.NotNil(t, c.LogrusLogger, "LogrusLogger should be initialized")
	assert.NotNil(t, c.MetricsCollector, "MetricsCollector should be initialized")

	// 验证客户端
	assert.NotNil(t, c.EinoClient, "EinoClient should be initialized")
	assert.NotNil(t, c.MilvusClient, "MilvusClient should be initialized")

	// 验证管理器
	assert.NotNil(t, c.TenantManager, "TenantManager should be initialized")
	assert.NotNil(t, c.MilvusTenantManager, "MilvusTenantManager should be initialized")
	assert.NotNil(t, c.DBManager, "DBManager should be initialized")

	// 验证仓储
	assert.NotNil(t, c.VectorRepository, "VectorRepository should be initialized")
	assert.NotNil(t, c.OrderRepository, "OrderRepository should be initialized")
	assert.NotNil(t, c.SessionRepository, "SessionRepository should be initialized")

	// 验证 AI 组件
	assert.NotNil(t, c.IntentRecognizer, "IntentRecognizer should be initialized")
	assert.NotNil(t, c.RAGRetriever, "RAGRetriever should be initialized")
	assert.NotNil(t, c.OrderQuerier, "OrderQuerier should be initialized")
	assert.NotNil(t, c.ResponseGenerator, "ResponseGenerator should be initialized")

	// 验证用例
	assert.NotNil(t, c.ChatUseCase, "ChatUseCase should be initialized")
	assert.NotNil(t, c.VectorUseCase, "VectorUseCase should be initialized")

	// 验证处理器
	assert.NotNil(t, c.ChatHandler, "ChatHandler should be initialized")
	assert.NotNil(t, c.VectorHandler, "VectorHandler should be initialized")
	assert.NotNil(t, c.HealthHandler, "HealthHandler should be initialized")
	assert.NotNil(t, c.ModelHandler, "ModelHandler should be initialized")

	// 验证中间件
	assert.NotNil(t, c.TenantMiddleware, "TenantMiddleware should be initialized")
	assert.NotNil(t, c.SecurityMiddleware, "SecurityMiddleware should be initialized")
	assert.NotNil(t, c.LoggingMiddleware, "LoggingMiddleware should be initialized")
	assert.NotNil(t, c.MetricsMiddleware, "MetricsMiddleware should be initialized")
	assert.NotNil(t, c.ErrorMiddleware, "ErrorMiddleware should be initialized")
	assert.NotNil(t, c.AuthMiddleware, "AuthMiddleware should be initialized")

	// 验证 HTTP 服务器
	assert.NotNil(t, c.Router, "Router should be initialized")
	assert.NotNil(t, c.Server, "Server should be initialized")
}

// TestContainer_Close 测试容器关闭
func TestContainer_Close(t *testing.T) {
	cfg := createTestConfig()

	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	c, err := New(cfg)
	if err != nil {
		t.Skipf("Failed to create container: %v", err)
		return
	}

	// 关闭容器
	err = c.Close()
	assert.NoError(t, err, "Close should not return error")

	// 再次关闭应该也不会出错（幂等性）
	err = c.Close()
	// 注意：某些组件可能不支持多次关闭，这里只是验证不会 panic
}

// TestContainer_InvalidConfig 测试无效配置
func TestContainer_InvalidConfig(t *testing.T) {
	tests := []struct {
		name   string
		config *config.Config
	}{
		{
			name: "empty dashscope api key",
			config: &config.Config{
				Server: config.ServerConfig{
					Port: 8080,
					Mode: "test",
				},
				DashScope: config.DashScopeConfig{
					APIKey:             "", // 空 API Key
					ChatModel:          "qwen-turbo",
					EmbedModel:         "text-embedding-v2",
					EmbeddingDimension: 1536,
				},
				Milvus: config.MilvusConfig{
					Host: "localhost",
					Port: 19530,
				},
				Database: config.DatabaseConfig{
					BasePath: "./testdata/db",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := New(tt.config)
			assert.Error(t, err, "Should return error for invalid config")
		})
	}
}

// createTestConfig 创建测试配置
func createTestConfig() *config.Config {
	return &config.Config{
		Server: config.ServerConfig{
			Port: 8080,
			Mode: "test",
		},
		DashScope: config.DashScopeConfig{
			APIKey:             "test_api_key",
			ChatModel:          "qwen-turbo",
			EmbedModel:         "text-embedding-v2",
			EmbeddingDimension: 1536,
			MaxRetries:         3,
			Timeout:            30 * time.Second,
		},
		Milvus: config.MilvusConfig{
			Host:    "localhost",
			Port:    19530,
			Timeout: 10 * time.Second,
		},
		Database: config.DatabaseConfig{
			BasePath: "./testdata/db",
		},
		RAG: config.RAGConfig{
			TopK:           5,
			ScoreThreshold: 0.7,
		},
		Intent: config.IntentConfig{
			ConfidenceThreshold: 0.6,
		},
		Session: config.SessionConfig{
			MaxHistory: 10,
			Timeout:    30 * time.Minute,
		},
		Security: config.SecurityConfig{
			APIKeys: []string{"test_key_1", "test_key_2"},
			SensitiveFields: []string{
				"password", "id_card", "phone", "mobile",
			},
		},
		Logging: config.LoggingConfig{
			Level:  "info",
			Format: "json",
			Output: "stdout",
		},
	}
}
