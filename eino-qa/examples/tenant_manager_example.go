package main

import (
	"context"
	"fmt"
	"log"

	"eino-qa/internal/infrastructure/config"
	"eino-qa/internal/infrastructure/tenant"

	"github.com/sirupsen/logrus"
)

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

	// 4. 获取默认租户（自动创建）
	fmt.Println("\n=== Getting Default Tenant ===")
	defaultTenant, err := manager.GetTenant(ctx, "default")
	if err != nil {
		log.Fatalf("Failed to get default tenant: %v", err)
	}
	fmt.Printf("Tenant ID: %s\n", defaultTenant.ID)
	fmt.Printf("Collection: %s\n", defaultTenant.CollectionName)
	fmt.Printf("Database: %s\n", defaultTenant.DatabasePath)

	// 5. 创建自定义租户
	fmt.Println("\n=== Creating Custom Tenant ===")
	customTenant, err := manager.CreateTenant(ctx, "company_abc", "ABC Company")
	if err != nil {
		log.Fatalf("Failed to create tenant: %v", err)
	}
	fmt.Printf("Created Tenant: %s (%s)\n", customTenant.ID, customTenant.Name)
	fmt.Printf("Collection: %s\n", customTenant.CollectionName)
	fmt.Printf("Database: %s\n", customTenant.DatabasePath)

	// 6. 检查租户是否存在
	fmt.Println("\n=== Checking Tenant Existence ===")
	exists, err := manager.TenantExists(ctx, "company_abc")
	if err != nil {
		log.Fatalf("Failed to check tenant: %v", err)
	}
	fmt.Printf("Tenant 'company_abc' exists: %v\n", exists)

	// 7. 获取租户信息
	fmt.Println("\n=== Getting Tenant Info ===")
	info, err := manager.GetTenantInfo(ctx, "company_abc")
	if err != nil {
		log.Fatalf("Failed to get tenant info: %v", err)
	}
	fmt.Printf("Tenant Info:\n")
	fmt.Printf("  ID: %s\n", info.TenantID)
	fmt.Printf("  Name: %s\n", info.Name)
	fmt.Printf("  Collection: %s (exists: %v)\n", info.CollectionName, info.CollectionExists)
	fmt.Printf("  Database: %s\n", info.DatabasePath)
	fmt.Printf("  DB Connections: %d\n", info.DBConnections)

	// 8. 获取租户的数据库连接
	fmt.Println("\n=== Getting Database Connection ===")
	db, err := manager.GetDB(ctx, "company_abc")
	if err != nil {
		log.Fatalf("Failed to get database: %v", err)
	}
	sqlDB, _ := db.DB()
	fmt.Printf("Database connection established\n")
	fmt.Printf("Open connections: %d\n", sqlDB.Stats().OpenConnections)

	// 9. 获取租户的 Collection 名称
	fmt.Println("\n=== Getting Collection Name ===")
	collectionName, err := manager.GetCollection(ctx, "company_abc")
	if err != nil {
		log.Fatalf("Failed to get collection: %v", err)
	}
	fmt.Printf("Collection name: %s\n", collectionName)

	// 10. 列出所有租户
	fmt.Println("\n=== Listing All Tenants ===")
	tenants := manager.ListTenants()
	fmt.Printf("Active tenants: %v\n", tenants)

	// 11. 演示多租户隔离
	fmt.Println("\n=== Multi-Tenant Isolation Demo ===")
	tenant1, _ := manager.GetTenant(ctx, "tenant1")
	tenant2, _ := manager.GetTenant(ctx, "tenant2")
	fmt.Printf("Tenant 1 - Collection: %s, DB: %s\n", tenant1.CollectionName, tenant1.DatabasePath)
	fmt.Printf("Tenant 2 - Collection: %s, DB: %s\n", tenant2.CollectionName, tenant2.DatabasePath)
	fmt.Println("Each tenant has isolated resources!")

	// 12. 清除缓存（可选）
	fmt.Println("\n=== Clearing Cache ===")
	manager.ClearCache()
	fmt.Println("Cache cleared")

	// 13. 再次获取租户（会重新加载）
	fmt.Println("\n=== Re-loading Tenant ===")
	reloadedTenant, err := manager.GetTenant(ctx, "company_abc")
	if err != nil {
		log.Fatalf("Failed to reload tenant: %v", err)
	}
	fmt.Printf("Reloaded Tenant: %s\n", reloadedTenant.ID)

	fmt.Println("\n=== Example Complete ===")
}
