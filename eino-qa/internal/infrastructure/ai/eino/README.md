# Eino AI 集成组件

本目录包含基于 CloudWeGo Eino 框架的 AI 集成组件实现。

## 组件概述

### 1. Client (client.go)
DashScope 客户端，负责初始化和管理聊天模型和嵌入模型。

**功能：**
- 初始化 DashScope ChatModel
- 初始化 DashScope Embedding 模型
- 提供统一的客户端接口

**使用示例：**
```go
client, err := eino.NewClient(eino.ClientConfig{
    APIKey:     "your-api-key",
    ChatModel:  "qwen-turbo",
    EmbedModel: "text-embedding-v2",
    MaxRetries: 3,
    Timeout:    30 * time.Second,
})
if err != nil {
    log.Fatal(err)
}
defer client.Close()
```

### 2. IntentRecognizer (intent_recognizer.go)
意图识别器，负责识别用户查询的意图类型。

**功能：**
- 识别用户查询意图（course/order/direct/handoff）
- 返回意图类型和置信度分数
- 支持对话历史上下文
- 低置信度自动转人工

**意图类型：**
- `course`: 课程咨询
- `order`: 订单查询
- `direct`: 直接回答
- `handoff`: 人工转接

**使用示例：**
```go
recognizer := eino.NewIntentRecognizer(client, &config.IntentConfig{
    ConfidenceThreshold: 0.7,
})

intent, err := recognizer.Recognize(ctx, "Python课程包含哪些内容？", history)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("意图: %s, 置信度: %.2f\n", intent.Type, intent.Confidence)
```

### 3. RAGRetriever (rag_retriever.go)
RAG 检索器，负责从知识库检索相关文档并生成答案。

**功能：**
- 生成查询向量
- 执行向量相似度搜索
- 过滤低分文档
- 基于检索文档生成答案

**使用示例：**
```go
retriever := eino.NewRAGRetriever(client, vectorRepo, &config.RAGConfig{
    TopK:           5,
    ScoreThreshold: 0.7,
})

answer, docs, err := retriever.Retrieve(ctx, "Python课程包含哪些内容？")
if err != nil {
    log.Fatal(err)
}

fmt.Printf("答案: %s\n", answer)
fmt.Printf("参考文档数: %d\n", len(docs))
```

### 4. OrderQuerier (order_querier.go)
订单查询器，负责查询订单信息并生成自然语言回复。

**功能：**
- 从用户查询中提取订单 ID（正则 + LLM）
- 查询订单数据库
- 格式化订单信息为自然语言
- SQL 安全验证

**使用示例：**
```go
querier := eino.NewOrderQuerier(client, orderRepo)

answer, err := querier.Query(ctx, "查询订单#20251114001")
if err != nil {
    log.Fatal(err)
}

fmt.Printf("答案: %s\n", answer)
```

### 5. ResponseGenerator (response_generator.go)
响应生成器，负责生成直接回答和各种系统消息。

**功能：**
- 生成直接回答（支持对话历史）
- 生成流式回答
- 生成人工转接消息
- 生成错误消息
- 生成降级消息

**使用示例：**
```go
generator := eino.NewResponseGenerator(client)

// 生成直接回答
answer, err := generator.Generate(ctx, "你好", history)
if err != nil {
    log.Fatal(err)
}

// 生成流式回答
resultChan, errorChan := generator.GenerateStream(ctx, "你好", history)
for chunk := range resultChan {
    fmt.Print(chunk)
}
if err := <-errorChan; err != nil {
    log.Fatal(err)
}
```

## 架构设计

### 依赖关系
```
Client (基础)
  ├── IntentRecognizer (意图识别)
  ├── RAGRetriever (RAG检索)
  │   └── VectorRepository (向量仓储)
  ├── OrderQuerier (订单查询)
  │   └── OrderRepository (订单仓储)
  └── ResponseGenerator (响应生成)
```

### 数据流
```
用户查询
  ↓
IntentRecognizer (识别意图)
  ↓
根据意图路由:
  ├─ course → RAGRetriever (检索知识库)
  ├─ order → OrderQuerier (查询订单)
  ├─ direct → ResponseGenerator (直接回答)
  └─ handoff → 人工转接消息
  ↓
返回答案
```

## 配置说明

### DashScope 配置
```yaml
dashscope:
  api_key: ${DASHSCOPE_API_KEY}
  chat_model: qwen-turbo
  embed_model: text-embedding-v2
  max_retries: 3
  timeout: 30s
```

### 意图识别配置
```yaml
intent:
  confidence_threshold: 0.7  # 置信度阈值，低于此值转人工
```

### RAG 配置
```yaml
rag:
  top_k: 5                   # 返回的最相似文档数
  score_threshold: 0.7       # 相似度分数阈值
```

## 错误处理

所有组件都遵循统一的错误处理策略：

1. **参数验证错误**: 返回明确的错误信息
2. **外部服务错误**: 包装原始错误，添加上下文
3. **业务逻辑错误**: 返回友好的错误消息

示例：
```go
if err != nil {
    return "", fmt.Errorf("failed to generate query vector: %w", err)
}
```

## 性能优化

### 1. 向量生成
- 批量生成向量以减少 API 调用
- 缓存常用查询的向量

### 2. 意图识别
- 使用轻量级模型（qwen-turbo）
- 限制历史消息数量（最近3轮）

### 3. RAG 检索
- 设置合理的 topK 值
- 使用分数阈值过滤低质量文档

## 测试

每个组件都应该有对应的单元测试：

```bash
# 运行所有测试
go test ./internal/infrastructure/ai/eino/...

# 运行特定测试
go test -v ./internal/infrastructure/ai/eino/ -run TestIntentRecognizer
```

## 注意事项

1. **API Key 安全**: 不要在代码中硬编码 API Key，使用环境变量
2. **超时设置**: 为所有外部调用设置合理的超时时间
3. **错误重试**: 对临时性错误实施重试策略
4. **日志记录**: 记录关键操作和错误信息
5. **资源清理**: 使用完毕后调用 Close() 方法

## 相关文档

- [Eino 官方文档](https://github.com/cloudwego/eino)
- [DashScope API 文档](https://help.aliyun.com/zh/dashscope/)
- [设计文档](../../../../.kiro/specs/eino-qa-system/design.md)
- [需求文档](../../../../.kiro/specs/eino-qa-system/requirements.md)
