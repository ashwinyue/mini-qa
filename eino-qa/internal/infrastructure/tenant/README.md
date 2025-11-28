# 租户管理器 (Tenant Manager)

## 概述

租户管理器是一个统一的多租户资源管理组件，负责协调和管理：
- Milvus 向量数据库的 Collection（每个租户独立的 Collection）
- SQLite 关系数据库（每个租户独立的数据库文件）
- 租户元数据和缓存

## 架构

```
┌─────────────────────────────────────────────────────────┐
│              Tenant Manager (统一管理器)                  │
│  ┌───────────────────────────────────────────────────┐  │
│  │  租户缓存 (Tenant Cache)                           │  │
│  │  - 租户 ID -> 租户实体映射                          │  │
│  │  - 自动创建和初始化                                 │  │
│  └───────────────────────────────────────────────────┘  │
│                                                          │
│  ┌──────────────────────┐  ┌──────────────────────┐    │
│  │ Milvus Tenant Mgr    │  │ SQLite DB Manager    │    │
│  │ - Collection 映射     │  │ - 数据库文件映射      │    │
│  │ - 自动创建 Collection │  │ - 连接池管理          │    │
│  └──────────────────────┘  └──────────────────────┘    │
└─────────────────────────────────────────────────────────┘
                    │                    │
                    ▼                    ▼
        ┌──────────────────┐  ┌──────────────────┐
        │  Milvus          │  │  SQLite          │
        │  kb_tenant1      │  │  tenant1.db      │
        │  kb_tenant2      │  │  tenant2.db      │
        │  kb_default      │  │  default.db      │
        └──────────────────┘  └──────────────────┘
```

## 核心功能

### 1. 租户资源自动创建

当访问一个新租户时，管理器会自动：
1. 创建租户实体
2. 在 Milvus 中创建对应的 Collection（格式：`kb_{tenantID}`）
3. 创建 SQLite 数据库文件（格式：`{tenantID}.db`）
4. 初始化数据库表结构
5. 缓存租户信息

### 2. 租户隔离

- **数据隔离**：每个租户拥有独立的 Milvus Collection 和 SQLite 数据库
- **资源隔离**：租户之间的数据操作完全隔离，互不影响
- **默认租户**：系统提供 `default` 租户作为默认选项

### 3. 并发安全

- 使用读写锁保护租户缓存
- 双重检查锁定模式防止并发创建
- 线程安全的数据库连接管理

## 使用示例

### 基本使用

```go
import (
    "context"
    "eino-qa/internal/infrastructure/tenant"
    "eino-qa/internal/infrastructure/config"
)

// 1. 从配置创建租户管理器
cfg, _ := config.LoadConfig("config/config.yaml")
logger := logrus.New()

manager, err := tenant.NewManagerFromConfig(tenant.FactoryConfig{
    Config: cfg,
    Logger: logger,
})
if err != nil {
    log.Fatal(err)
}
defer manager.Close()

// 2. 获取租户（自动创建）
ctx := context.Background()
tenant, err := manager.GetTenant(ctx, "tenant1")
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Tenant: %s, Collection: %s\n", tenant.ID, tenant.CollectionName)

// 3. 获取租户的数据库连接
db, err := manager.GetDB(ctx, "tenant1")
if err != nil {
    log.Fatal(err)
}

// 使用数据库进行操作
var count int64
db.Model(&Order{}).Count(&count)

// 4. 获取租户的 Collection 名称
collectionName, err := manager.GetCollection(ctx, "tenant1")
if err != nil {
    log.Fatal(err)
}
```

### 显式创建租户

```go
// 创建具有自定义名称的租户
tenant, err := manager.CreateTenant(ctx, "company_abc", "ABC Company")
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Created tenant: %s (%s)\n", tenant.ID, tenant.Name)
```

### 检查租户是否存在

```go
exists, err := manager.TenantExists(ctx, "tenant1")
if err != nil {
    log.Fatal(err)
}

if exists {
    fmt.Println("Tenant exists")
} else {
    fmt.Println("Tenant does not exist")
}
```

### 获取租户信息

```go
info, err := manager.GetTenantInfo(ctx, "tenant1")
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Tenant Info:\n")
fmt.Printf("  ID: %s\n", info.TenantID)
fmt.Printf("  Name: %s\n", info.Name)
fmt.Printf("  Collection: %s (exists: %v)\n", info.CollectionName, info.CollectionExists)
fmt.Printf("  Database: %s\n", info.DatabasePath)
fmt.Printf("  DB Connections: %d\n", info.DBConnections)
```

### 列出所有租户

```go
tenants := manager.ListTenants()
fmt.Printf("Active tenants: %v\n", tenants)
```

### 删除租户

```go
// 注意：这会删除租户的所有数据！
err := manager.DeleteTenant(ctx, "tenant1")
if err != nil {
    log.Fatal(err)
}
```

## 在 HTTP 中间件中使用

```go
func TenantMiddleware(manager *tenant.Manager) gin.HandlerFunc {
    return func(c *gin.Context) {
        // 从请求头或查询参数提取租户 ID
        tenantID := c.GetHeader("X-Tenant-ID")
        if tenantID == "" {
            tenantID = c.Query("tenant")
        }
        if tenantID == "" {
            tenantID = "default"
        }

        // 确保租户存在（自动创建）
        tenant, err := manager.GetTenant(c.Request.Context(), tenantID)
        if err != nil {
            c.JSON(500, gin.H{"error": "failed to initialize tenant"})
            c.Abort()
            return
        }

        // 将租户信息存储到上下文
        c.Set("tenant_id", tenant.ID)
        c.Set("tenant", tenant)

        c.Next()
    }
}
```

## 在 Use Case 中使用

```go
type ChatUseCase struct {
    tenantManager *tenant.Manager
    // ... 其他依赖
}

func (uc *ChatUseCase) Execute(ctx context.Context, req ChatRequest) (*ChatResponse, error) {
    // 从上下文获取租户 ID
    tenantID, ok := ctx.Value("tenant_id").(string)
    if !ok {
        tenantID = "default"
    }

    // 获取租户的数据库连接
    db, err := uc.tenantManager.GetDB(ctx, tenantID)
    if err != nil {
        return nil, err
    }

    // 获取租户的 Collection 名称
    collectionName, err := uc.tenantManager.GetCollection(ctx, tenantID)
    if err != nil {
        return nil, err
    }

    // 使用租户特定的资源进行操作
    // ...
}
```

## 配置

租户管理器需要以下配置：

```yaml
# config.yaml
milvus:
  host: localhost
  port: 19530
  username: ""
  password: ""

database:
  base_path: "data/db"  # SQLite 数据库文件存储路径

dashscope:
  embedding_dimension: 1536  # 向量维度
```

## 注意事项

1. **默认租户**：系统始终提供 `default` 租户，不能删除
2. **资源清理**：使用完毕后应调用 `manager.Close()` 释放资源
3. **并发安全**：管理器是并发安全的，可以在多个 goroutine 中使用
4. **自动创建**：首次访问租户时会自动创建所有必需的资源
5. **错误处理**：应妥善处理租户创建失败的情况

## 测试

运行单元测试：

```bash
cd internal/infrastructure/tenant
go test -v
```

运行集成测试（需要 Milvus 服务）：

```bash
go test -v -tags=integration
```

## 性能考虑

1. **缓存机制**：租户信息会被缓存，避免重复创建
2. **连接池**：SQLite 使用 GORM 的连接池管理
3. **延迟初始化**：只在首次访问时创建租户资源
4. **并发控制**：使用读写锁优化并发访问性能

## 扩展性

租户管理器设计为可扩展的：

1. 可以添加更多的租户级资源（如 Redis、Kafka 等）
2. 可以实现租户配额和限流
3. 可以添加租户级别的监控和统计
4. 可以实现租户的热备份和迁移
