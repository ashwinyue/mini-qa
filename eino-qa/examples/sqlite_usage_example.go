package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"eino-qa/internal/domain/entity"
	"eino-qa/internal/infrastructure/repository/sqlite"
)

func main() {
	// 创建仓储工厂
	factory := sqlite.NewRepositoryFactory("./data/db")
	defer factory.Close()

	ctx := context.Background()
	tenantID := "tenant1"

	// 示例 1: 订单操作
	fmt.Println("=== 订单操作示例 ===")
	orderRepo := factory.GetOrderRepository(tenantID)

	// 创建订单
	order := entity.NewOrder("user123", "Go 高级编程课程", 299.00, tenantID)
	if err := orderRepo.Create(ctx, order); err != nil {
		log.Fatalf("创建订单失败: %v", err)
	}
	fmt.Printf("创建订单成功: %s\n", order.ID)

	// 查询订单
	foundOrder, err := orderRepo.FindByID(ctx, order.ID)
	if err != nil {
		log.Fatalf("查询订单失败: %v", err)
	}
	fmt.Printf("查询订单: ID=%s, 课程=%s, 金额=%.2f, 状态=%s\n",
		foundOrder.ID, foundOrder.CourseName, foundOrder.Amount, foundOrder.Status)

	// 更新订单状态
	if err := foundOrder.UpdateStatus(entity.OrderStatusPaid); err != nil {
		log.Fatalf("更新订单状态失败: %v", err)
	}
	if err := orderRepo.Update(ctx, foundOrder); err != nil {
		log.Fatalf("保存订单失败: %v", err)
	}
	fmt.Printf("订单状态已更新为: %s\n", foundOrder.Status)

	// 按用户查询订单
	userOrders, err := orderRepo.FindByUserID(ctx, "user123")
	if err != nil {
		log.Fatalf("查询用户订单失败: %v", err)
	}
	fmt.Printf("用户 user123 的订单数量: %d\n", len(userOrders))

	// 示例 2: 会话操作
	fmt.Println("\n=== 会话操作示例 ===")
	sessionRepo := factory.GetSessionRepository(tenantID)

	// 创建会话
	session := entity.NewSession(tenantID, 24*time.Hour)
	fmt.Printf("创建会话: %s\n", session.ID)

	// 添加消息
	userMsg := entity.NewMessage("你好，我想了解 Go 课程", "user")
	if err := session.AddMessage(userMsg); err != nil {
		log.Fatalf("添加用户消息失败: %v", err)
	}

	assistantMsg := entity.NewMessage("您好！我们的 Go 课程包含基础和高级内容...", "assistant")
	if err := session.AddMessage(assistantMsg); err != nil {
		log.Fatalf("添加助手消息失败: %v", err)
	}

	// 保存会话
	if err := sessionRepo.Save(ctx, session); err != nil {
		log.Fatalf("保存会话失败: %v", err)
	}
	fmt.Printf("会话已保存，消息数量: %d\n", session.GetMessageCount())

	// 加载会话
	loadedSession, err := sessionRepo.Load(ctx, session.ID)
	if err != nil {
		log.Fatalf("加载会话失败: %v", err)
	}
	fmt.Printf("加载会话: ID=%s, 消息数量=%d\n", loadedSession.ID, loadedSession.GetMessageCount())

	// 打印消息
	for i, msg := range loadedSession.GetMessages() {
		fmt.Printf("  消息 %d: [%s] %s\n", i+1, msg.Role, msg.Content)
	}

	// 示例 3: 未命中查询记录
	fmt.Println("\n=== 未命中查询记录示例 ===")
	missedRepo := factory.GetMissedQueryRepository(tenantID)

	// 记录未命中查询
	if err := missedRepo.Create(ctx, "这个问题没有答案", "course"); err != nil {
		log.Fatalf("记录未命中查询失败: %v", err)
	}
	fmt.Println("未命中查询已记录")

	// 查询未命中记录
	missedQueries, err := missedRepo.List(ctx, 0, 10)
	if err != nil {
		log.Fatalf("查询未命中记录失败: %v", err)
	}
	fmt.Printf("未命中查询数量: %d\n", len(missedQueries))

	// 统计
	orderCount, _ := orderRepo.Count(ctx)
	sessionCount, _ := sessionRepo.Count(ctx)
	missedCount, _ := missedRepo.Count(ctx)

	fmt.Println("\n=== 统计信息 ===")
	fmt.Printf("订单总数: %d\n", orderCount)
	fmt.Printf("会话总数: %d\n", sessionCount)
	fmt.Printf("未命中查询总数: %d\n", missedCount)

	fmt.Println("\n所有操作完成！")
}
