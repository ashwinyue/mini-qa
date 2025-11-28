package eino

import (
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	tests := []struct {
		name    string
		config  ClientConfig
		wantErr bool
	}{
		{
			name: "valid config",
			config: ClientConfig{
				APIKey:     "test_api_key",
				ChatModel:  "qwen-turbo",
				EmbedModel: "text-embedding-v2",
				MaxRetries: 3,
				Timeout:    30 * time.Second,
			},
			wantErr: false,
		},
		{
			name: "missing api key",
			config: ClientConfig{
				APIKey:     "",
				ChatModel:  "qwen-turbo",
				EmbedModel: "text-embedding-v2",
			},
			wantErr: true,
		},
		{
			name: "default values",
			config: ClientConfig{
				APIKey: "test_api_key",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := NewClient(tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if client == nil {
					t.Error("expected non-nil client")
					return
				}

				if client.GetChatModel() == nil {
					t.Error("expected non-nil chat model")
				}

				if client.GetEmbedModel() == nil {
					t.Error("expected non-nil embed model")
				}

				// 测试默认值
				if tt.config.ChatModel == "" && client.config.ChatModel != "qwen-turbo" {
					t.Errorf("expected default chat model 'qwen-turbo', got '%s'", client.config.ChatModel)
				}

				if tt.config.EmbedModel == "" && client.config.EmbedModel != "text-embedding-v2" {
					t.Errorf("expected default embed model 'text-embedding-v2', got '%s'", client.config.EmbedModel)
				}

				if tt.config.MaxRetries == 0 && client.config.MaxRetries != 3 {
					t.Errorf("expected default max retries 3, got %d", client.config.MaxRetries)
				}

				if tt.config.Timeout == 0 && client.config.Timeout != 30*time.Second {
					t.Errorf("expected default timeout 30s, got %v", client.config.Timeout)
				}
			}
		})
	}
}

func TestClient_Close(t *testing.T) {
	client, err := NewClient(ClientConfig{
		APIKey: "test_api_key",
	})
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	if err := client.Close(); err != nil {
		t.Errorf("Close() error = %v", err)
	}
}
