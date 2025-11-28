package config

import (
	"os"
	"testing"
)

func TestLoad(t *testing.T) {
	// 创建临时配置文件
	tmpFile, err := os.CreateTemp("", "config-*.yaml")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	// 写入测试配置
	configContent := `
server:
  port: 8080
  mode: debug

dashscope:
  api_key: test_api_key
  chat_model: qwen-turbo
  embed_model: text-embedding-v2
  max_retries: 3
  timeout: 30s

milvus:
  host: localhost
  port: 19530
  username: ""
  password: ""
  timeout: 10s

database:
  base_path: ./data/sqlite

rag:
  top_k: 5
  score_threshold: 0.7

intent:
  confidence_threshold: 0.6

session:
  max_history: 10
  timeout: 30m

security:
  api_keys:
    - key1
    - key2
  sensitive_fields:
    - password
    - id_card

logging:
  level: info
  format: json
  output: stdout
  file_path: ./logs/app.log
`
	if _, err := tmpFile.WriteString(configContent); err != nil {
		t.Fatalf("failed to write config: %v", err)
	}
	tmpFile.Close()

	// 测试加载配置
	config, err := Load(tmpFile.Name())
	if err != nil {
		t.Fatalf("failed to load config: %v", err)
	}

	// 验证配置
	if config.Server.Port != 8080 {
		t.Errorf("expected port 8080, got %d", config.Server.Port)
	}

	if config.DashScope.APIKey != "test_api_key" {
		t.Errorf("expected api_key 'test_api_key', got '%s'", config.DashScope.APIKey)
	}

	if config.Milvus.Host != "localhost" {
		t.Errorf("expected milvus host 'localhost', got '%s'", config.Milvus.Host)
	}

	if config.Database.BasePath != "./data/sqlite" {
		t.Errorf("expected database base_path './data/sqlite', got '%s'", config.Database.BasePath)
	}
}

func TestValidate(t *testing.T) {
	tests := []struct {
		name    string
		config  Config
		wantErr bool
	}{
		{
			name: "valid config",
			config: Config{
				Server: ServerConfig{
					Port: 8080,
					Mode: "debug",
				},
				DashScope: DashScopeConfig{
					APIKey: "test_key",
				},
				Milvus: MilvusConfig{
					Host: "localhost",
				},
				Database: DatabaseConfig{
					BasePath: "./data",
				},
			},
			wantErr: false,
		},
		{
			name: "invalid port",
			config: Config{
				Server: ServerConfig{
					Port: 0,
				},
				DashScope: DashScopeConfig{
					APIKey: "test_key",
				},
				Milvus: MilvusConfig{
					Host: "localhost",
				},
				Database: DatabaseConfig{
					BasePath: "./data",
				},
			},
			wantErr: true,
		},
		{
			name: "missing api key",
			config: Config{
				Server: ServerConfig{
					Port: 8080,
				},
				DashScope: DashScopeConfig{
					APIKey: "",
				},
				Milvus: MilvusConfig{
					Host: "localhost",
				},
				Database: DatabaseConfig{
					BasePath: "./data",
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
