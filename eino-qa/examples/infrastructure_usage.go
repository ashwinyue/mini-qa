package main

import (
	"context"
	"fmt"
	"log"

	"eino-qa/internal/infrastructure/ai/eino"
	"eino-qa/internal/infrastructure/config"
	"eino-qa/internal/infrastructure/logger"
)

func main() {
	// 1. 加载配置
	cfg, err := config.Load("../config/config.yaml")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	fmt.Println("✓ 配置加载成功")
	fmt.Printf("  - Server Port: %d\n", cfg.Server.Port)
	fmt.Printf("  - Chat Model: %s\n", cfg.DashScope.ChatModel)
	fmt.Printf("  - Embed Model: %s\n", cfg.DashScope.EmbedModel)

	// 2. 初始化日志系统
	lgr, err := logger.New(logger.Config{
		Level:  cfg.Logging.Level,
		Format: cfg.Logging.Format,
		Output: cfg.Logging.Output,
	})
	if err != nil {
		log.Fatalf("failed to create logger: %v", err)
	}

	fmt.Println("\n✓ 日志系统初始化成功")

	// 3. 创建带有上下文的日志
	ctx := context.Background()
	ctx = context.WithValue(ctx, "trace_id", "example-trace-001")
	ctx = context.WithValue(ctx, "tenant_id", "default")

	lgr.Info(ctx, "系统启动", map[string]interface{}{
		"version": "1.0.0",
		"env":     cfg.Server.Mode,
	})

	// 4. 初始化 DashScope 客户端
	client, err := eino.NewClient(eino.ClientConfig{
		APIKey:     cfg.DashScope.APIKey,
		ChatModel:  cfg.DashScope.ChatModel,
		EmbedModel: cfg.DashScope.EmbedModel,
		MaxRetries: cfg.DashScope.MaxRetries,
		Timeout:    cfg.DashScope.Timeout,
	})
	if err != nil {
		lgr.Error(ctx, "DashScope 客户端初始化失败", map[string]interface{}{
			"error": err.Error(),
		})
		log.Fatalf("failed to create dashscope client: %v", err)
	}
	defer client.Close()

	fmt.Println("\n✓ DashScope 客户端初始化成功")
	fmt.Printf("  - Chat Model: %v\n", client.GetChatModel() != nil)
	fmt.Printf("  - Embed Model: %v\n", client.GetEmbedModel() != nil)

	lgr.Info(ctx, "所有基础设施组件初始化完成", map[string]interface{}{
		"components": []string{"config", "logger", "dashscope"},
	})

	fmt.Println("\n✓ 所有基础设施组件初始化完成")
	fmt.Println("\n示例运行成功！")
}
