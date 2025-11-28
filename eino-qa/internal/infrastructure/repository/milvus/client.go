package milvus

import (
	"context"
	"fmt"
	"time"

	"github.com/milvus-io/milvus-sdk-go/v2/client"
	"github.com/sirupsen/logrus"
)

// Client Milvus 客户端封装
type Client struct {
	client client.Client
	logger *logrus.Logger
}

// ClientConfig Milvus 客户端配置
type ClientConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	Timeout  time.Duration
}

// NewClient 创建新的 Milvus 客户端
func NewClient(config ClientConfig, logger *logrus.Logger) (*Client, error) {
	if logger == nil {
		logger = logrus.New()
	}

	// 创建 Milvus 连接
	ctx, cancel := context.WithTimeout(context.Background(), config.Timeout)
	defer cancel()

	address := fmt.Sprintf("%s:%d", config.Host, config.Port)

	logger.WithFields(logrus.Fields{
		"address": address,
		"timeout": config.Timeout,
	}).Info("connecting to Milvus")

	milvusClient, err := client.NewClient(ctx, client.Config{
		Address:  address,
		Username: config.Username,
		Password: config.Password,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Milvus: %w", err)
	}

	logger.Info("successfully connected to Milvus")

	return &Client{
		client: milvusClient,
		logger: logger,
	}, nil
}

// GetClient 获取底层 Milvus 客户端
func (c *Client) GetClient() client.Client {
	return c.client
}

// Close 关闭 Milvus 连接
func (c *Client) Close() error {
	if c.client != nil {
		c.logger.Info("closing Milvus connection")
		return c.client.Close()
	}
	return nil
}

// Ping 检查 Milvus 连接状态
func (c *Client) Ping(ctx context.Context) error {
	// 通过列出数据库来检查连接
	_, err := c.client.ListDatabases(ctx)
	if err != nil {
		return fmt.Errorf("Milvus ping failed: %w", err)
	}
	return nil
}
