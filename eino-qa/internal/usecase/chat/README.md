# Chat UseCase

## 概述

Chat UseCase 实现了智能客服系统的核心对话逻辑，负责协调意图识别、信息检索、订单查询和响应生成等组件，完成完整的对话流程。

## 架构设计

### 核心组件

1. **ChatUseCase**: 主要的对话用例实现
   - 会话管理
   - 意图路由
   - 组件编排
   - 错误处理

2. **流式响应支持**: 
   - `ExecuteStream`: 支持流式对话
   - 实时推送响应内容
   - 适用于长文本生成场景

3. **并行信息收集**:
   - `ExecuteParallel`: 并行查询多个数据源
   - `ExecuteWithParallelRetrieval`: 智能并行检索
   - 提升响应速度

## 主要功能

### 1. 标准对话流程

```go
func (uc *ChatUseCase) Execute(ctx context.Context, req *ChatRequest) (*ChatResponse, error)
```

**流程说明**:
1. 验证请求参数
2. 加载或创建会话
3. 添加用户消息到会话历史
4. 识别用户意图
5. 根据意图路由到相应处理器：
   - `IntentCourse`: RAG 知识库检索
   - `IntentOrder`: 订单数据库查询
   - `IntentDirect`: 直接对话生成
   - `IntentHandoff`: 人工转接
6. 添加助手响应到会话历史
7. 保存会话状态
8. 返回响应结果

### 2. 流式对话

```go
func (uc *ChatUseCase) ExecuteStream(ctx context.Context, req *ChatRequest) (<-chan *StreamChunk, error)
```

**特点**:
- 实时推送响应内容
- 支持长文本生成
- 降低首字延迟
- 提升用户体验

### 3. 并行信息收集

```go
func (uc *ChatUseCase) ExecuteParallel(ctx context.Context, req *ChatRequest) (*ParallelQueryResult, error)
```

**应用场景**:
- 需要同时查询多个数据源
- 提升响应速度
- 信息整合需求

**实现方式**:
- 使用 Go 的 goroutine 并发执行
- WaitGroup 同步等待
- 支持超时控制

## 意图路由

### 意图类型

1. **课程咨询 (IntentCourse)**
   - 使用 RAG 检索器
   - 从向量数据库检索相关文档
   - 基于检索结果生成答案

2. **订单查询 (IntentOrder)**
   - 使用订单查询器
   - 从关系数据库查询订单信息
   - 格式化为自然语言响应

3. **直接回答 (IntentDirect)**
   - 使用响应生成器
   - 基于对话历史生成回答
   - 适用于简单问候和闲聊

4. **人工转接 (IntentHandoff)**
   - 低置信度查询
   - 复杂问题
   - 投诉处理

## 会话管理

### 会话生命周期

1. **创建**: 首次对话时创建新会话
2. **加载**: 后续对话加载已有会话
3. **更新**: 每次对话更新会话历史
4. **延期**: 活跃会话自动延长过期时间
5. **过期**: 超时会话自动清理

### 会话存储

- 使用 SessionRepository 接口
- 支持多种存储实现（SQLite、Redis 等）
- 租户级隔离

## 错误处理

### 错误类型

1. **验证错误**: 请求参数不合法
2. **会话错误**: 会话加载或保存失败
3. **意图识别错误**: LLM 调用失败
4. **路由错误**: 具体处理器执行失败

### 降级策略

1. **RAG 失败**: 返回降级消息，建议联系人工
2. **订单查询失败**: 返回友好错误提示
3. **响应生成失败**: 返回预设兜底回复
4. **会话保存失败**: 记录日志但不影响响应

## 性能优化

### 1. 并行执行

- 多数据源并行查询
- 减少总体响应时间
- 提升吞吐量

### 2. 流式响应

- 降低首字延迟
- 提升用户体验
- 适用于长文本生成

### 3. 会话缓存

- 减少数据库查询
- 快速加载历史记录
- 支持会话延期

## 日志记录

### 关键日志点

1. 请求开始/完成
2. 意图识别结果
3. 路由处理过程
4. 错误和异常
5. 性能指标（响应时间）

### 日志字段

- `tenant_id`: 租户标识
- `session_id`: 会话标识
- `query`: 用户查询
- `intent`: 识别的意图
- `confidence`: 置信度
- `duration_ms`: 处理时长

## 使用示例

### 标准对话

```go
useCase := NewChatUseCase(
    intentRecognizer,
    ragRetriever,
    orderQuerier,
    responseGenerator,
    sessionRepo,
    30*time.Minute,
    logger,
)

req := &ChatRequest{
    Query:     "Python课程包含哪些内容？",
    TenantID:  "tenant1",
    SessionID: "sess_123",
    Stream:    false,
}

resp, err := useCase.Execute(ctx, req)
if err != nil {
    log.Fatal(err)
}

fmt.Println("Answer:", resp.Answer)
fmt.Println("Route:", resp.Route)
```

### 流式对话

```go
req := &ChatRequest{
    Query:     "介绍一下Go语言的特点",
    TenantID:  "tenant1",
    SessionID: "sess_123",
    Stream:    true,
}

chunkChan, err := useCase.ExecuteStream(ctx, req)
if err != nil {
    log.Fatal(err)
}

for chunk := range chunkChan {
    if chunk.Error != nil {
        log.Printf("Error: %v\n", chunk.Error)
        break
    }
    
    if chunk.Done {
        fmt.Println("\nDone!")
        break
    }
    
    fmt.Print(chunk.Content)
}
```

### 并行查询

```go
req := &ChatRequest{
    Query:     "查询我的课程和订单",
    TenantID:  "tenant1",
}

result, err := useCase.ExecuteParallel(ctx, req)
if err != nil {
    log.Fatal(err)
}

fmt.Println("Course Answer:", result.CourseAnswer)
fmt.Println("Order Answer:", result.OrderAnswer)
fmt.Println("Duration:", result.Duration)
```

## 扩展性

### 添加新的意图类型

1. 在 `entity.IntentType` 中定义新类型
2. 在 `Execute` 方法的 switch 语句中添加处理分支
3. 实现对应的处理器

### 添加新的数据源

1. 实现数据源的查询接口
2. 在 `ExecuteParallel` 中添加并行查询逻辑
3. 在 `MergeParallelResults` 中添加结果合并逻辑

## 测试

### 单元测试

- 测试各个意图的路由逻辑
- 测试会话管理功能
- 测试错误处理机制

### 集成测试

- 测试完整的对话流程
- 测试多租户隔离
- 测试并发场景

### 性能测试

- 测试响应时间
- 测试并发处理能力
- 测试资源使用情况

## 依赖关系

```
ChatUseCase
├── IntentRecognizer (意图识别)
├── RAGRetriever (知识库检索)
├── OrderQuerier (订单查询)
├── ResponseGenerator (响应生成)
├── SessionRepository (会话存储)
└── Logger (日志记录)
```

## 注意事项

1. **会话管理**: 确保会话正确保存，避免数据丢失
2. **错误处理**: 所有错误都应该有友好的用户提示
3. **性能监控**: 记录关键性能指标，便于优化
4. **并发安全**: 注意并发场景下的数据一致性
5. **资源清理**: 及时清理过期会话，避免内存泄漏
