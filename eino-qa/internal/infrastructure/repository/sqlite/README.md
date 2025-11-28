# SQLite Repository 实现

本包实现了基于 SQLite 的数据持久化层，支持多租户隔离。

## 功能特性

- ✅ 多租户数据库隔离（每个租户独立的 SQLite 文件）
- ✅ GORM ORM 支持
- ✅ 自动表结构迁移
- ✅ 订单管理（Order）
- ✅ 会话管理（Session）
- ✅ 未命中查询记录（MissedQuery）
- ✅ 线程安全的数据库连接管理

## 架构设计

### 数据库文件组织

```
data/db/
├── tenant1.db      # 租户1的数据库
├── tenant2.db      # 租户2的数据库
└── default.db      # 默认租户的数据库
```

每个租户拥有独立的 SQLite 数据库文件，实现完全的数据隔离。

### 表结构

#### orders 表
```sql
CREATE TABLE orders (
    id VARCHAR(50) PRIMARY KEY,
    user_id VARCHAR(100) NOT NULL,
    course_name VARCHAR(200) NOT NULL,
    amount DECIMAL(10,2) NOT NULL,
    status VARCHAR(20) NOT NULL,
    tenant_id VARCHAR(100) NOT NULL,
    metadata TEXT,
    created_at DATETIME,
    updated_at DATETIME
);
CREATE INDEX idx_orders_user_id ON orders(user_id);
CREATE INDEX idx_orders_status ON orders(status);
CREATE INDEX idx_orders_tenant_id ON orders(tenant_id);
```

#### sessions 表
```sql
CREATE TABLE sessions (
    id VARCHAR(100) PRIMARY KEY,
    tenant_id VARCHAR(100) NOT NULL,
    messages TEXT,
    metadata TEXT,
    created_at DATETIME,
    updated_at DATETIME,
    expires_at DATETIME NOT NULL
);
CREATE INDEX idx_sessions_tenant_id ON sessions(tenant_id);
CREATE INDEX idx_sessions_expires_at ON sessions(expires_at);
```

#### missed_queries 表
```sql
CREATE TABLE missed_queries (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    tenant_id VARCHAR(100) NOT NULL,
    query TEXT NOT NULL,
    intent VARCHAR(50),
    created_at DATETIME
);
CREATE INDEX idx_missed_queries_tenant_id ON missed_queries(tenant_id);
```

## 使用示例

### 1. 初始化仓储工厂

```go
import "eino-qa/internal/infrastructure/repository/sqlite"

// 创建仓储工厂
factory := sqlite.NewRepositoryFactory("./data/db")
defer factory.Close()
```

### 2. 订单操作

```go
// 获取订单仓储
orderRepo := factory.GetOrderRepository("tenant1")

// 创建订单
order := entity.NewOrder("user123", "Go 课程", 299.00, "tenant1")
err := orderRepo.Create(ctx, order)

// 查询订单
order, err := orderRepo.FindByID(ctx, orderID)

// 按用户查询
orders, err := orderRepo.FindByUserID(ctx, "user123")

// 按状态查询
orders, err := orderRepo.FindByStatus(ctx, entity.OrderStatusPaid)

// 更新订单
order.UpdateStatus(entity.OrderStatusPaid)
err = orderRepo.Update(ctx, order)

// 删除订单
err = orderRepo.Delete(ctx, orderID)

// 分页列表
orders, err := orderRepo.List(ctx, offset, limit)

// 统计数量
count, err := orderRepo.Count(ctx)
```

### 3. 会话操作

```go
// 获取会话仓储
sessionRepo := factory.GetSessionRepository("tenant1")

// 创建会话
session := entity.NewSession("tenant1", 24*time.Hour)
userMsg := entity.NewMessage("你好", "user")
session.AddMessage(userMsg)

// 保存会话
err := sessionRepo.Save(ctx, session)

// 加载会话
session, err := sessionRepo.Load(ctx, sessionID)

// 添加消息
msg := entity.NewMessage("回复内容", "assistant")
err = sessionRepo.AddMessage(ctx, sessionID, msg)

// 获取消息
messages, err := sessionRepo.GetMessages(ctx, sessionID)

// 更新过期时间
err = sessionRepo.UpdateExpiration(ctx, sessionID, time.Now().Add(48*time.Hour))

// 删除过期会话
count, err := sessionRepo.DeleteExpired(ctx)

// 列出租户会话
sessions, err := sessionRepo.ListByTenant(ctx, "tenant1")

// 删除会话
err = sessionRepo.Delete(ctx, sessionID)
```

### 4. 未命中查询记录

```go
// 获取未命中查询仓储
missedRepo := factory.GetMissedQueryRepository("tenant1")

// 记录未命中查询
err := missedRepo.Create(ctx, "无法回答的问题", "course")

// 列出记录
queries, err := missedRepo.List(ctx, offset, limit)

// 统计数量
count, err := missedRepo.Count(ctx)

// 删除旧记录
count, err := missedRepo.DeleteOlderThan(ctx, time.Now().AddDate(0, -1, 0))
```

## 多租户隔离

### DBManager

`DBManager` 负责管理多租户的数据库连接：

- 每个租户一个独立的 SQLite 数据库文件
- 连接池管理和缓存
- 自动表结构迁移
- 线程安全

```go
// 获取租户数据库连接
db, err := dbManager.GetDB("tenant1")

// 列出所有已连接的租户
tenants := dbManager.ListTenants()

// 移除租户连接
err = dbManager.RemoveDB("tenant1")

// 关闭所有连接
err = dbManager.Close()
```

## 数据模型转换

### OrderModel ↔ entity.Order

```go
// 领域实体 -> GORM 模型
var model OrderModel
err := model.FromEntity(order)

// GORM 模型 -> 领域实体
order, err := model.ToEntity()
```

### SessionModel ↔ entity.Session

```go
// 领域实体 -> GORM 模型
var model SessionModel
err := model.FromEntity(session)

// GORM 模型 -> 领域实体
session, err := model.ToEntity()
```

## 错误处理

所有仓储方法都返回详细的错误信息：

```go
order, err := orderRepo.FindByID(ctx, orderID)
if err != nil {
    if strings.Contains(err.Error(), "not found") {
        // 订单不存在
    } else {
        // 其他错误
    }
}
```

## 性能优化

1. **连接池**: 每个租户的数据库连接被缓存和重用
2. **索引**: 关键字段（user_id, status, tenant_id）都建立了索引
3. **批量操作**: 支持分页查询，避免一次加载大量数据
4. **JSON 序列化**: Metadata 和 Messages 使用 JSON 存储，灵活且高效

## 测试

运行示例程序：

```bash
cd eino-qa
go run examples/sqlite_usage_example.go
```

## 注意事项

1. **租户 ID 验证**: 所有操作都会验证租户 ID 是否匹配
2. **数据验证**: 创建和更新操作会调用实体的 `Validate()` 方法
3. **事务支持**: 可以通过 GORM 的事务 API 实现复杂的事务操作
4. **并发安全**: DBManager 使用读写锁保证并发安全
5. **资源清理**: 使用 `defer factory.Close()` 确保资源正确释放

## 依赖

- `gorm.io/gorm`: GORM ORM 框架
- `gorm.io/driver/sqlite`: SQLite 驱动
- `github.com/mattn/go-sqlite3`: SQLite C 库绑定

## 未来扩展

- [ ] 支持数据库备份和恢复
- [ ] 支持数据迁移工具
- [ ] 支持读写分离
- [ ] 支持分布式事务
- [ ] 支持数据加密
