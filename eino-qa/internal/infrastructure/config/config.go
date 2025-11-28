package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// Config 系统配置
type Config struct {
	Server    ServerConfig    `yaml:"server"`
	DashScope DashScopeConfig `yaml:"dashscope"`
	Milvus    MilvusConfig    `yaml:"milvus"`
	Database  DatabaseConfig  `yaml:"database"`
	RAG       RAGConfig       `yaml:"rag"`
	Intent    IntentConfig    `yaml:"intent"`
	Session   SessionConfig   `yaml:"session"`
	Security  SecurityConfig  `yaml:"security"`
	Logging   LoggingConfig   `yaml:"logging"`
}

// ServerConfig HTTP 服务器配置
type ServerConfig struct {
	Port int    `yaml:"port"`
	Mode string `yaml:"mode"` // debug, release
}

// DashScopeConfig DashScope API 配置
type DashScopeConfig struct {
	APIKey             string        `yaml:"api_key"`
	ChatModel          string        `yaml:"chat_model"`
	EmbedModel         string        `yaml:"embed_model"`
	EmbeddingDimension int           `yaml:"embedding_dimension"`
	MaxRetries         int           `yaml:"max_retries"`
	Timeout            time.Duration `yaml:"timeout"`
}

// MilvusConfig Milvus 向量数据库配置
type MilvusConfig struct {
	Host     string        `yaml:"host"`
	Port     int           `yaml:"port"`
	Username string        `yaml:"username"`
	Password string        `yaml:"password"`
	Timeout  time.Duration `yaml:"timeout"`
}

// DatabaseConfig SQLite 数据库配置
type DatabaseConfig struct {
	BasePath string `yaml:"base_path"`
}

// RAGConfig RAG 检索配置
type RAGConfig struct {
	TopK           int     `yaml:"top_k"`
	ScoreThreshold float64 `yaml:"score_threshold"`
}

// IntentConfig 意图识别配置
type IntentConfig struct {
	ConfidenceThreshold float64 `yaml:"confidence_threshold"`
}

// SessionConfig 会话管理配置
type SessionConfig struct {
	MaxHistory int           `yaml:"max_history"`
	Timeout    time.Duration `yaml:"timeout"`
}

// SecurityConfig 安全配置
type SecurityConfig struct {
	APIKeys         []string `yaml:"api_keys"`
	SensitiveFields []string `yaml:"sensitive_fields"`
}

// LoggingConfig 日志配置
type LoggingConfig struct {
	Level    string `yaml:"level"`
	Format   string `yaml:"format"`
	Output   string `yaml:"output"`
	FilePath string `yaml:"file_path"`
}

// Load 从文件加载配置
func Load(configPath string) (*Config, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	// 展开环境变量
	expanded := os.ExpandEnv(string(data))

	var config Config
	if err := yaml.Unmarshal([]byte(expanded), &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// 验证配置
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return &config, nil
}

// Validate 验证配置的有效性
func (c *Config) Validate() error {
	if c.Server.Port <= 0 || c.Server.Port > 65535 {
		return fmt.Errorf("invalid server port: %d", c.Server.Port)
	}

	if c.DashScope.APIKey == "" {
		return fmt.Errorf("dashscope api_key is required")
	}

	if c.Milvus.Host == "" {
		return fmt.Errorf("milvus host is required")
	}

	if c.Database.BasePath == "" {
		return fmt.Errorf("database base_path is required")
	}

	return nil
}
