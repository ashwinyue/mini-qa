package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"eino-qa/internal/domain/entity"
	"eino-qa/internal/infrastructure/ai/eino"
	"eino-qa/internal/infrastructure/config"
	"eino-qa/internal/infrastructure/logger"
	"eino-qa/internal/infrastructure/repository/sqlite"
	"eino-qa/internal/usecase/chat"
)

// 这是一个示例程序，展示如何使用 ChatUseCase
// 注意：这个示例需要实际的配置和数据库连接才能运行

func main() {
	// 1. 初始化配置
	cfg, err := config.Load("config/config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 2. 初始化日志
	logger, err := logger.New(logger.Config{
		Level:  "info",
		Format: "json",
		Output: "stdout",
	})
	if err != nil {
		log.Fatalf("Failed to create logger: %v", err)
	}

	// 3. 初始化 Eino 客户端
	einoClient, err := eino.NewClient(&cfg.DashScope)
	if err != nil {
		log.Fatalf("Failed to create Eino client: %v", err)
	}

	// 4. 初始化 AI 组件
	intentRecognizer := eino.NewIntentRecognizer(einoClient, &cfg.Intent)
	responseGenerator := eino.NewResponseGenerator(einoClient)

	// 注意：这里需要实际的 VectorRepository 和 OrderRepository 实现
	// 为了示例简化，这里使用 nil，实际使用时需要替换
	var ragRetriever *eino.RAGRetriever = nil
	var orderQuerier *eino.OrderQuerier = nil

	// 5. 初始化会话仓储
	dbManager := sqlite.NewDBManager(cfg.Database.BasePath)
	sessionRepo := sqlite.NewSessionRepository(dbManager)

	// 6. 创建 ChatUseCase
	chatUseCase := chat.NewChatUseCase(
		intentRecognizer,
		ragRetriever,
		orderQuerier,
		responseGenerator,
		sessionRepo,
		30*time.Minute, // 会话 TTL
		logger,
	)

	// 7. 执行对话
	ctx := context.Background()

	// 示例 1: 简单问候
	fmt.Println("=== 示例 1: 简单问候 ===")
	req1 := &chat.ChatRequest{
		Query:     "你好",
		TenantID:  "tenant1",
		SessionID: "",
		Stream:    false,
	}

	resp1, err := chatUseCase.Execute(ctx, req1)
	if err != nil {
		log.Printf("Error: %v", err)
	} else {
		fmt.Printf("Answer: %s\n", resp1.Answer)
		fmt.Printf("Route: %s\n", resp1.Route)
		fmt.Printf("Session ID: %s\n", resp1.SessionID)
	}

	// 示例 2: 课程咨询（需要 RAG 组件）
	fmt.Println("\n=== 示例 2: 课程咨询 ===")
	req2 := &chat.ChatRequest{
		Query:     "Python课程包含哪些内容？",
		TenantID:  "tenant1",
		SessionID: resp1.SessionID, // 使用上一个会话
		Stream:    false,
	}

	resp2, err := chatUseCase.Execute(ctx, req2)
	if err != nil {
		log.Printf("Error: %v", err)
	} else {
		fmt.Printf("Answer: %s\n", resp2.Answer)
		fmt.Printf("Route: %s\n", resp2.Route)
		if len(resp2.Sources) > 0 {
			fmt.Printf("Sources: %d documents\n", len(resp2.Sources))
		}
	}

	// 示例 3: 流式对话
	fmt.Println("\n=== 示例 3: 流式对话 ===")
	req3 := &chat.ChatRequest{
		Query:     "介绍一下Go语言的特点",
		TenantID:  "tenant1",
		SessionID: resp1.SessionID,
		Stream:    true,
	}

	chunkChan, err := chatUseCase.ExecuteStream(ctx, req3)
	if err != nil {
		log.Printf("Error: %v", err)
	} else {
		fmt.Print("Answer: ")
		for chunk := range chunkChan {
			if chunk.Error != nil {
				log.Printf("Stream error: %v", chunk.Error)
				break
			}

			if chunk.Done {
				fmt.Println("\n[Done]")
				if metadata, ok := chunk.Metadata["intent"]; ok {
					fmt.Printf("Intent: %v\n", metadata)
				}
				break
			}

			fmt.Print(chunk.Content)
		}
	}

	// 示例 4: 并行查询（需要完整的组件）
	fmt.Println("\n=== 示例 4: 并行查询 ===")
	req4 := &chat.ChatRequest{
		Query:    "查询我的课程和订单",
		TenantID: "tenant1",
	}

	parallelResult, err := chatUseCase.ExecuteParallel(ctx, req4)
	if err != nil {
		log.Printf("Error: %v", err)
	} else {
		fmt.Printf("Course Answer: %s\n", parallelResult.CourseAnswer)
		fmt.Printf("Order Answer: %s\n", parallelResult.OrderAnswer)
		fmt.Printf("Duration: %v\n", parallelResult.Duration)
	}
}

// 辅助函数：打印会话历史
func printSessionHistory(session *entity.Session) {
	fmt.Println("\n=== 会话历史 ===")
	fmt.Printf("Session ID: %s\n", session.ID)
	fmt.Printf("Tenant ID: %s\n", session.TenantID)
	fmt.Printf("Message Count: %d\n", session.GetMessageCount())
	fmt.Println("\nMessages:")
	for i, msg := range session.GetMessages() {
		fmt.Printf("%d. [%s] %s\n", i+1, msg.Role, msg.Content)
	}
}
