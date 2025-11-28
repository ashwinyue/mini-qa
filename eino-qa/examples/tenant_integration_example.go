package main

import (
	"context"
	"fmt"
	"log"

	"eino-qa/internal/domain/entity"
	"eino-qa/internal/infrastructure/config"
	"eino-qa/internal/infrastructure/repository/sqlite"
	"eino-qa/internal/infrastructure/tenant"

	"github.com/sirupsen/logrus"
)

// 这个示例展示如何在实际应用中集成租户管理器
func main() {
	// 1. 加载配置
	cfg, err := config.Load("config/config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 2. 创建 logger
	logger := logrus.New()
	logger.SetLevel(logrus.InfoLevel)

	// 3. 创建租户管理器
	manager, err := tenant.NewManagerFromConfig(tenant.FactoryConfig{
		Config: cfg,
		Logger: logger,
	})
	if err != nil {
		log.Fatalf("Failed to create tenant manager: %v", err)
	}
	defer manager.Close()

	ctx := context.Background()

	// 4. 模拟多租户场景
	fmt.Println("=== Multi-Tenant Application Demo ===\n")

	// 场景 1: 租户 A 的用户查询订单
	fmt.Println("--- Scenario 1: Tenant A User Query ---")
	if err := handleTenantAQuery(ctx, manager); err != nil {
		log.Printf("Error in tenant A: %v", err)
	}

	// 场景 2: 租户 B 的用户查询订单
	fmt.Println("\n--- Scenario 2: Tenant B User Query ---")
	if err := handleTenantBQuery(ctx, manager); err != nil {
		log.Printf("Error in tenant B: %v", err)
	}

	// 场景 3: 验证租户隔离
	fmt.Println("\n--- Scenario 3: Verify Tenant Isolation ---")
	verifyTenantIsolation(ctx, manager)

	fmt.Println("\n=== Demo Complete ===")
}

// handleTenantAQuery 处理租户 A 的查询
func handleTenantAQuery(ctx context.Context, manager *tenant.Manager) error {
	tenantID := "tenant_a"

	// 获取租户（自动创建）
	t, err := manager.GetTenant(ctx, tenantID)
	if err != nil {
		return fmt.Errorf("failed to get tenant: %w", err)
	}

	fmt.Printf("Processing query for tenant: %s\n", t.ID)
	fmt.Printf("Using collection: %s\n", t.CollectionName)
	fmt.Printf("Using database: %s\n", t.DatabasePath)

	// 获取租户的数据库连接
	db, err := manager.GetDB(ctx, tenantID)
	if err != nil {
		return fmt.Errorf("failed to get database: %w", err)
	}

	// 创建测试订单
	order := &entity.Order{
		ID:         "order_a_001",
		UserID:     "user_a_001",
		TenantID:   tenantID,
		CourseName: "Python 基础课程",
		Amount:     299.00,
		Status:     entity.OrderStatusPaid,
	}

	// 使用 OrderRepository 保存订单
	orderRepo := sqlite.NewOrderRepository(
		&sqlite.DBManager{},
		tenantID,
	)

	// 直接使用数据库连接创建订单
	var orderModel sqlite.OrderModel
	if err := orderModel.FromEntity(order); err != nil {
		return fmt.Errorf("failed to convert order: %w", err)
	}

	if err := db.Create(&orderModel).Error; err != nil {
		return fmt.Errorf("failed to create order: %w", err)
	}

	fmt.Printf("Created order: %s for user: %s\n", order.ID, order.UserID)

	// 查询订单
	var foundOrder sqlite.OrderModel
	if err := db.Where("id = ?", order.ID).First(&foundOrder).Error; err != nil {
		return fmt.Errorf("failed to find order: %w", err)
	}

	fmt.Printf("Found order: %s, amount: %.2f\n", foundOrder.ID, foundOrder.Amount)

	return nil
}

// handleTenantBQuery 处理租户 B 的查询
func handleTenantBQuery(ctx context.Context, manager *tenant.Manager) error {
	tenantID := "tenant_b"

	// 获取租户（自动创建）
	t, err := manager.GetTenant(ctx, tenantID)
	if err != nil {
		return fmt.Errorf("failed to get tenant: %w", err)
	}

	fmt.Printf("Processing query for tenant: %s\n", t.ID)
	fmt.Printf("Using collection: %s\n", t.CollectionName)
	fmt.Printf("Using database: %s\n", t.DatabasePath)

	// 获取租户的数据库连接
	db, err := manager.GetDB(ctx, tenantID)
	if err != nil {
		return fmt.Errorf("failed to get database: %w", err)
	}

	// 创建测试订单
	order := &entity.Order{
		ID:         "order_b_001",
		UserID:     "user_b_001",
		TenantID:   tenantID,
		CourseName: "Go 语言进阶",
		Amount:     399.00,
		Status:     entity.OrderStatusPaid,
	}

	// 直接使用数据库连接创建订单
	var orderModel sqlite.OrderModel
	if err := orderModel.FromEntity(order); err != nil {
		return fmt.Errorf("failed to convert order: %w", err)
	}

	if err := db.Create(&orderModel).Error; err != nil {
		return fmt.Errorf("failed to create order: %w", err)
	}

	fmt.Printf("Created order: %s for user: %s\n", order.ID, order.UserID)

	// 查询订单
	var foundOrder sqlite.OrderModel
	if err := db.Where("id = ?", order.ID).First(&foundOrder).Error; err != nil {
		return fmt.Errorf("failed to find order: %w", err)
	}

	fmt.Printf("Found order: %s, amount: %.2f\n", foundOrder.ID, foundOrder.Amount)

	return nil
}

// verifyTenantIsolation 验证租户隔离
func verifyTenantIsolation(ctx context.Context, manager *tenant.Manager) {
	// 获取租户 A 的数据库
	dbA, err := manager.GetDB(ctx, "tenant_a")
	if err != nil {
		log.Printf("Failed to get tenant A database: %v", err)
		return
	}

	// 获取租户 B 的数据库
	dbB, err := manager.GetDB(ctx, "tenant_b")
	if err != nil {
		log.Printf("Failed to get tenant B database: %v", err)
		return
	}

	// 统计租户 A 的订单数
	var countA int64
	dbA.Model(&sqlite.OrderModel{}).Count(&countA)
	fmt.Printf("Tenant A has %d orders\n", countA)

	// 统计租户 B 的订单数
	var countB int64
	dbB.Model(&sqlite.OrderModel{}).Count(&countB)
	fmt.Printf("Tenant B has %d orders\n", countB)

	// 尝试在租户 A 的数据库中查询租户 B 的订单（应该找不到）
	var orderB sqlite.OrderModel
	result := dbA.Where("id = ?", "order_b_001").First(&orderB)
	if result.Error != nil {
		fmt.Println("✓ Tenant isolation verified: Tenant A cannot access Tenant B's orders")
	} else {
		fmt.Println("✗ Tenant isolation failed: Tenant A can access Tenant B's orders")
	}

	// 获取租户信息
	infoA, _ := manager.GetTenantInfo(ctx, "tenant_a")
	infoB, _ := manager.GetTenantInfo(ctx, "tenant_b")

	fmt.Printf("\nTenant A Info:\n")
	fmt.Printf("  Collection: %s\n", infoA.CollectionName)
	fmt.Printf("  Database: %s\n", infoA.DatabasePath)

	fmt.Printf("\nTenant B Info:\n")
	fmt.Printf("  Collection: %s\n", infoB.CollectionName)
	fmt.Printf("  Database: %s\n", infoB.DatabasePath)

	fmt.Println("\n✓ Each tenant has isolated resources!")
}
