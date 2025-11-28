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

	// 为不同租户创建订单
	tenants := []string{"tenant1", "tenant2", "tenant3"}

	fmt.Println("=== 多租户隔离测试 ===\n")

	// 为每个租户创建订单
	for i, tenantID := range tenants {
		orderRepo := factory.GetOrderRepository(tenantID)

		order := entity.NewOrder(
			fmt.Sprintf("user%d", i+1),
			fmt.Sprintf("课程 %d", i+1),
			float64((i+1)*100),
			tenantID,
		)

		if err := orderRepo.Create(ctx, order); err != nil {
			log.Fatalf("创建订单失败 (租户 %s): %v", tenantID, err)
		}

		fmt.Printf("租户 %s: 创建订单 %s\n", tenantID, order.ID)
	}

	fmt.Println("\n=== 验证租户隔离 ===\n")

	// 验证每个租户只能看到自己的订单
	for _, tenantID := range tenants {
		orderRepo := factory.GetOrderRepository(tenantID)

		count, err := orderRepo.Count(ctx)
		if err != nil {
			log.Fatalf("统计订单失败 (租户 %s): %v", tenantID, err)
		}

		orders, err := orderRepo.List(ctx, 0, 100)
		if err != nil {
			log.Fatalf("列出订单失败 (租户 %s): %v", tenantID, err)
		}

		fmt.Printf("租户 %s: 订单数量 = %d\n", tenantID, count)
		for _, order := range orders {
			fmt.Printf("  - 订单 %s: 用户=%s, 课程=%s, 租户=%s\n",
				order.ID, order.UserID, order.CourseName, order.TenantID)
		}
		fmt.Println()
	}

	// 验证数据库文件
	fmt.Println("=== 数据库文件 ===\n")
	dbManager := factory.GetDBManager()
	connectedTenants := dbManager.ListTenants()
	fmt.Printf("已连接的租户: %v\n", connectedTenants)

	// 会话隔离测试
	fmt.Println("\n=== 会话隔离测试 ===\n")

	for i, tenantID := range tenants {
		sessionRepo := factory.GetSessionRepository(tenantID)

		session := entity.NewSession(tenantID, 24*time.Hour)
		msg := entity.NewMessage(fmt.Sprintf("租户 %d 的消息", i+1), "user")
		session.AddMessage(msg)

		if err := sessionRepo.Save(ctx, session); err != nil {
			log.Fatalf("保存会话失败 (租户 %s): %v", tenantID, err)
		}

		fmt.Printf("租户 %s: 创建会话 %s\n", tenantID, session.ID)
	}

	fmt.Println("\n=== 验证会话隔离 ===\n")

	for _, tenantID := range tenants {
		sessionRepo := factory.GetSessionRepository(tenantID)

		count, err := sessionRepo.Count(ctx)
		if err != nil {
			log.Fatalf("统计会话失败 (租户 %s): %v", tenantID, err)
		}

		sessions, err := sessionRepo.ListByTenant(ctx, tenantID)
		if err != nil {
			log.Fatalf("列出会话失败 (租户 %s): %v", tenantID, err)
		}

		fmt.Printf("租户 %s: 会话数量 = %d\n", tenantID, count)
		for _, session := range sessions {
			fmt.Printf("  - 会话 %s: 消息数=%d, 租户=%s\n",
				session.ID, session.GetMessageCount(), session.TenantID)
		}
		fmt.Println()
	}

	fmt.Println("多租户隔离测试完成！")
}
