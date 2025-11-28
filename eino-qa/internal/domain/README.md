# Domain Layer

Domain Layer 是系统的核心层，包含业务实体、值对象和仓储接口定义。这一层完全独立于外部框架和技术实现。

## 目录结构

```
domain/
├── entity/          # 核心实体和值对象
│   ├── message.go       # 消息实体
│   ├── intent.go        # 意图实体
│   ├── document.go      # 文档实体
│   ├── order.go         # 订单实体
│   ├── session.go       # 会话实体
│   ├── tenant.go        # 租户值对象
│   ├── query_result.go  # 查询结果值对象
│   ├── validator.go     # 业务规则验证器
│   ├── errors.go        # 错误定义
│   └── entity_test.go   # 单元测试
└── repository/      # 仓储接口
    ├── vector.go        # 向量仓储接口
    ├── order.go         # 订单仓储接口
    └── session.go       # 会话仓储接口
```

## 核心实体

### Message（消息）
表示对话中的一条消息，包含内容、角色（user/assistant/system）和时间戳。

**主要方法：**
- `NewMessage(content, role)` - 创建新消息
- `Validate()` - 验证消息有效性
- `IsUser()`, `IsAssistant()`, `IsSystem()` - 角色判断

### Intent（意图）
表示用户查询的意图类型和置信度。

**意图类型：**
- `IntentCourse` - 课程咨询
- `IntentOrder` - 订单查询
- `IntentDirect` - 直接回答
- `IntentHandoff` - 人工转接

**主要方法：**
- `NewIntent(type, confidence)` - 创建新意图
- `Validate()` - 验证意图有效性
- `IsHighConfidence(threshold)` - 判断置信度

### Document（文档）
表示知识库中的文档，包含内容、向量表示和元数据。

**主要方法：**
- `NewDocument(content, tenantID)` - 创建新文档
- `SetVector(vector)` - 设置向量
- `SetScore(score)` - 设置相似度分数
- `AddMetadata(key, value)` - 添加元数据

### Order（订单）
表示订单实体，包含用户信息、课程名称、金额和状态。

**订单状态：**
- `OrderStatusPending` - 待支付
- `OrderStatusPaid` - 已支付
- `OrderStatusRefunded` - 已退款
- `OrderStatusCancelled` - 已取消

**主要方法：**
- `NewOrder(userID, courseName, amount, tenantID)` - 创建新订单
- `UpdateStatus(status)` - 更新订单状态
- `IsPending()`, `IsPaid()`, `IsRefunded()`, `IsCancelled()` - 状态判断

### Session（会话）
表示对话会话，包含消息历史和过期时间。

**主要方法：**
- `NewSession(tenantID, ttl)` - 创建新会话
- `AddMessage(message)` - 添加消息
- `GetMessages()` - 获取所有消息
- `IsExpired()` - 判断是否过期
- `ExtendExpiration(duration)` - 延长过期时间

## 值对象

### Tenant（租户）
表示租户信息，包含 ID、名称、Collection 名称和数据库路径。

**主要方法：**
- `NewTenant(id, name)` - 创建新租户
- `IsDefault()` - 判断是否为默认租户

### QueryResult（查询结果）
表示查询结果，包含答案、路由类型、来源文档和元数据。

**主要方法：**
- `NewQueryResult(answer, route)` - 创建新查询结果
- `AddSource(doc)` - 添加来源文档
- `SetIntent(intent)` - 设置意图

## 业务规则验证器

`validator.go` 提供了一系列业务规则验证函数：

- `ValidateOrderID(orderID)` - 验证订单 ID 格式（#YYYYMMDDXXX）
- `ValidateSessionID(sessionID)` - 验证会话 ID 格式
- `ValidateTenantID(tenantID)` - 验证租户 ID
- `ValidateConfidence(confidence)` - 验证置信度分数（0-1）
- `ExtractOrderID(text)` - 从文本中提取订单 ID
- `SanitizeSQL(sql)` - 清理 SQL 查询，防止注入
- `ValidateVector(vector, expectedDim)` - 验证向量维度

## 仓储接口

### VectorRepository
定义向量数据库操作接口，用于 Milvus 集成。

**主要方法：**
- `Search(ctx, vector, topK)` - 向量相似度搜索
- `Insert(ctx, docs)` - 插入文档向量
- `Delete(ctx, ids)` - 删除文档向量
- `CreateCollection(ctx, name, dimension)` - 创建集合
- `CollectionExists(ctx, name)` - 检查集合是否存在

### OrderRepository
定义订单数据库操作接口，用于 SQLite 集成。

**主要方法：**
- `FindByID(ctx, orderID)` - 根据 ID 查询订单
- `FindByUserID(ctx, userID)` - 根据用户 ID 查询订单
- `FindByStatus(ctx, status)` - 根据状态查询订单
- `Create(ctx, order)` - 创建订单
- `Update(ctx, order)` - 更新订单
- `Delete(ctx, orderID)` - 删除订单

### SessionRepository
定义会话存储操作接口，用于会话管理。

**主要方法：**
- `Save(ctx, session)` - 保存会话
- `Load(ctx, sessionID)` - 加载会话
- `Delete(ctx, sessionID)` - 删除会话
- `AddMessage(ctx, sessionID, message)` - 添加消息
- `GetMessages(ctx, sessionID)` - 获取消息列表
- `DeleteExpired(ctx)` - 删除过期会话

## 设计原则

1. **依赖倒置**：Domain Layer 不依赖任何外部框架，所有依赖都通过接口定义
2. **单一职责**：每个实体只负责自己的业务逻辑
3. **不变性**：值对象（如 Tenant）创建后不可变
4. **验证规则**：所有实体都提供 `Validate()` 方法进行自我验证
5. **业务规则封装**：业务规则集中在 `validator.go` 中，便于维护和测试

## 测试

运行 Domain Layer 测试：

```bash
go test ./internal/domain/entity/... -v
go test ./internal/domain/repository/... -v
```

所有实体和验证器都有对应的单元测试，确保业务逻辑的正确性。
