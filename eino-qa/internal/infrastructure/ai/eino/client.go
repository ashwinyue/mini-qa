package eino

import (
	"context"
	"fmt"
	"time"

	arkEmbed "github.com/cloudwego/eino-ext/components/embedding/ark"
	arkModel "github.com/cloudwego/eino-ext/components/model/ark"
	"github.com/cloudwego/eino/components/embedding"
	"github.com/cloudwego/eino/components/model"
)

// ClientConfig DashScope 客户端配置
type ClientConfig struct {
	APIKey     string
	ChatModel  string
	EmbedModel string
	MaxRetries int
	Timeout    time.Duration
}

// Client DashScope 客户端
type Client struct {
	chatModel  model.ChatModel
	embedModel embedding.Embedder
	config     ClientConfig
}

// NewClient 创建新的 DashScope 客户端
func NewClient(config ClientConfig) (*Client, error) {
	if config.APIKey == "" {
		return nil, fmt.Errorf("dashscope api_key is required")
	}

	if config.ChatModel == "" {
		config.ChatModel = "qwen-turbo"
	}

	if config.EmbedModel == "" {
		config.EmbedModel = "text-embedding-v2"
	}

	if config.MaxRetries <= 0 {
		config.MaxRetries = 3
	}

	if config.Timeout <= 0 {
		config.Timeout = 30 * time.Second
	}

	// 初始化聊天模型（使用 Ark，兼容 DashScope）
	chatModel, err := arkModel.NewChatModel(
		context.Background(),
		&arkModel.ChatModelConfig{
			APIKey: config.APIKey,
			Model:  config.ChatModel,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize chat model: %w", err)
	}

	// 初始化嵌入模型
	embedModel, err := arkEmbed.NewEmbedder(
		context.Background(),
		&arkEmbed.EmbeddingConfig{
			APIKey: config.APIKey,
			Model:  config.EmbedModel,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize embedding model: %w", err)
	}

	return &Client{
		chatModel:  chatModel,
		embedModel: embedModel,
		config:     config,
	}, nil
}

// GetChatModel 获取聊天模型
func (c *Client) GetChatModel() model.ChatModel {
	return c.chatModel
}

// GetEmbedModel 获取嵌入模型
func (c *Client) GetEmbedModel() embedding.Embedder {
	return c.embedModel
}

// Close 关闭客户端
func (c *Client) Close() error {
	// DashScope 客户端目前不需要显式关闭
	return nil
}
