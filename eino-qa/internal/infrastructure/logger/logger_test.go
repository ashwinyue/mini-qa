package logger

import (
	"context"
	"testing"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name    string
		config  Config
		wantErr bool
	}{
		{
			name: "valid json stdout config",
			config: Config{
				Level:  "info",
				Format: "json",
				Output: "stdout",
			},
			wantErr: false,
		},
		{
			name: "valid text stdout config",
			config: Config{
				Level:  "debug",
				Format: "text",
				Output: "stdout",
			},
			wantErr: false,
		},
		{
			name: "invalid level",
			config: Config{
				Level:  "invalid",
				Format: "json",
				Output: "stdout",
			},
			wantErr: true,
		},
		{
			name: "invalid format",
			config: Config{
				Level:  "info",
				Format: "invalid",
				Output: "stdout",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := New(tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestLogger_Info(t *testing.T) {
	logger, err := New(Config{
		Level:  "info",
		Format: "json",
		Output: "stdout",
	})
	if err != nil {
		t.Fatalf("failed to create logger: %v", err)
	}

	// 创建带有上下文的 context
	ctx := context.Background()
	ctx = context.WithValue(ctx, "trace_id", "test-trace-123")
	ctx = context.WithValue(ctx, "tenant_id", "test-tenant")

	// 记录日志
	logger.Info(ctx, "test message", map[string]interface{}{
		"key1": "value1",
		"key2": 123,
	})

	// 注意：由于输出到 stdout，我们无法直接捕获，但可以验证不会 panic
}

func TestLogger_WithFields(t *testing.T) {
	logger, err := New(Config{
		Level:  "info",
		Format: "json",
		Output: "stdout",
	})
	if err != nil {
		t.Fatalf("failed to create logger: %v", err)
	}

	// 创建带有预设字段的 logger
	loggerWithFields := logger.WithFields(map[string]interface{}{
		"service": "test-service",
		"version": "1.0.0",
	})

	// 记录日志
	ctx := context.Background()
	loggerWithFields.Info(ctx, "test message", map[string]interface{}{
		"extra": "data",
	})

	// 验证不会 panic
}

func TestExtractContextFields(t *testing.T) {
	ctx := context.Background()
	ctx = context.WithValue(ctx, "trace_id", "trace-123")
	ctx = context.WithValue(ctx, "tenant_id", "tenant-456")
	ctx = context.WithValue(ctx, "session_id", "session-789")

	fields := extractContextFields(ctx)

	if fields["trace_id"] != "trace-123" {
		t.Errorf("expected trace_id 'trace-123', got '%v'", fields["trace_id"])
	}

	if fields["tenant_id"] != "tenant-456" {
		t.Errorf("expected tenant_id 'tenant-456', got '%v'", fields["tenant_id"])
	}

	if fields["session_id"] != "session-789" {
		t.Errorf("expected session_id 'session-789', got '%v'", fields["session_id"])
	}
}

func TestMergeFields(t *testing.T) {
	base := map[string]interface{}{
		"key1": "value1",
		"key2": "value2",
	}

	additional := map[string]interface{}{
		"key2": "new_value2",
		"key3": "value3",
	}

	result := mergeFields(base, additional)

	if result["key1"] != "value1" {
		t.Errorf("expected key1 'value1', got '%v'", result["key1"])
	}

	if result["key2"] != "new_value2" {
		t.Errorf("expected key2 'new_value2', got '%v'", result["key2"])
	}

	if result["key3"] != "value3" {
		t.Errorf("expected key3 'value3', got '%v'", result["key3"])
	}
}
