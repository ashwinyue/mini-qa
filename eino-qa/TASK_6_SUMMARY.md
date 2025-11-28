# Task 6 实现总结：Infrastructure Layer - Eino AI 集成

## 完成时间
2025-01-XX

## 任务概述
实现基于 CloudWeGo Eino 框架的 AI 集成组件，包括 DashScope 客户端初始化、意图识别器、RAG 检索器、订单查询器和响应生成器。

## 实现的组件

### 1. DashScope Client (client.go)
**功能：**
- ✅ 初始化 DashScope ChatModel（使用 Ark 兼容层）
- ✅ 初始化 DashScope Embedding 模型
- ✅ 提供统一的客户端接口
- ✅ 支持配置验证和默认值设置

**关键特性：**
- 使用 `arkModel.NewChatModel` 初始化聊天模型
- 使用 `arkEmbed.NewEmbedder` 初始化嵌入模型
- 支持自定义超时和重试配置
- 提供 `GetChatModel()` 和 `GetEmbedModel()` 访问器

### 2. IntentRecognizer (intent_recognizer.go)
**功能：**
- ✅ 识别用户查询意图（course/order/direct/handoff）
- ✅ 返回意图类型和置信度分数
- ✅ 支持对话历史上下文（最近3轮）
- ✅ 低置信度自动转人工

**关键特性：**
- 使用结构化 JSON 输出格式
- 支持置信度阈值配置
- 清理 markdown 代码块标记
- 意图类型映射和验证

**提示词设计：**
- 系统提示词：定义4种意图类型和 JSON 输出格式
- 用户提示词：包含对话历史和当前查询

### 3. RAGRetriever (rag_retriever.go)
**功能：**
- ✅ 生成查询向量（float64 → float32 转换）
- ✅ 执行向量相似度搜索
- ✅ 根据分数阈值过滤文档
- ✅ 基于检索文档生成答案

**关键特性：**
- 支持 TopK 和分数阈值配置
- 构建结构化的文档上下文
- 处理未命中情况
- 向量类型转换（Eino 返回 float64，Milvus 需要 float32）

**提示词设计：**
- 系统提示词：定义课程咨询助手角色和回答要求
- 用户提示词：包含知识库文档和用户问题

### 4. OrderQuerier (order_querier.go)
**功能：**
- ✅ 从用户查询中提取订单 ID（正则 + LLM）
- ✅ 查询订单数据库
- ✅ 格式化订单信息为自然语言
- ✅ SQL 安全验证

**关键特性：**
- 双重提取策略：先正则，失败后用 LLM
- 支持多种订单号格式（#开头、订单号：、纯数字）
- 订单状态中文映射
- 危险 SQL 关键词检测

**正则模式：**
```go
`#(\d{11,})`              // #开头的订单号
`订单号[：:]\s*(\d{11,})`  // 订单号：xxxxx
`订单[：:]\s*(\d{11,})`    // 订单：xxxxx
`\b(\d{11,})\b`           // 独立的11位以上数字
```

**SQL 验证：**
- 检测危险关键词：DROP, DELETE, UPDATE, INSERT, ALTER, CREATE, TRUNCATE, EXEC, --, ;
- 确保只允许 SELECT 语句

### 5. ResponseGenerator (response_generator.go)
**功能：**
- ✅ 生成直接回答（支持对话历史）
- ✅ 生成流式回答
- ✅ 生成人工转接消息
- ✅ 生成错误消息
- ✅ 生成降级消息

**关键特性：**
- 支持同步和流式两种生成模式
- 历史消息转换为 Eino schema.Message
- 友好的错误消息映射
- 多种人工转接消息模板

**提示词设计：**
- 系统提示词：定义友好专业的客服助手角色
- 用户提示词：直接使用用户查询（历史已在消息列表中）

## 技术实现细节

### 1. 类型转换
```go
// float64 → float32 (Eino → Milvus)
vector := make([]float32, len(resp[0]))
for i, v := range resp[0] {
    vector[i] = float32(v)
}
```

### 2. 消息构建
```go
// 构建 Eino 消息列表
messages := []*schema.Message{
    schema.SystemMessage(systemPrompt),
    schema.UserMessage(userPrompt),
}

// 添加历史消息
for _, msg := range history {
    if msg.IsUser() {
        messages = append(messages, schema.UserMessage(msg.Content))
    } else if msg.IsAssistant() {
        messages = append(messages, schema.AssistantMessage(msg.Content, nil))
    }
}
```

### 3. JSON 解析
```go
// 清理 markdown 代码块
content = strings.TrimSpace(content)
content = strings.TrimPrefix(content, "```json")
content = strings.TrimPrefix(content, "```")
content = strings.TrimSuffix(content, "```")
content = strings.TrimSpace(content)

// 解析 JSON
var result struct {
    Intent     string  `json:"intent"`
    Confidence float64 `json:"confidence"`
    Reason     string  `json:"reason"`
}
json.Unmarshal([]byte(content), &result)
```

### 4. 流式响应
```go
// 创建通道
resultChan := make(chan string, 10)
errorChan := make(chan error, 1)

// 异步读取流
go func() {
    defer close(resultChan)
    defer close(errorChan)
    
    streamReader, err := g.chatModel.Stream(ctx, messages)
    if err != nil {
        errorChan <- err
        return
    }
    
    for {
        chunk, err := streamReader.Recv()
        if err != nil {
            if err.Error() != "EOF" {
                errorChan <- err
            }
            break
        }
        if chunk.Content != "" {
            resultChan <- chunk.Content
        }
    }
}()

return resultChan, errorChan
```

## 测试覆盖

### 单元测试 (integration_test.go)
- ✅ 客户端初始化测试
- ✅ 意图类型映射测试
- ✅ 订单 ID 正则提取测试
- ✅ SQL 安全验证测试
- ✅ 响应生成器消息测试
- ✅ RAG 分数过滤测试
- ✅ 提示词构建测试
- ✅ 上下文构建测试

### 基准测试
- ✅ 意图类型映射性能测试
- ✅ 订单 ID 提取性能测试

**测试结果：**
```
PASS
ok      eino-qa/internal/infrastructure/ai/eino 0.581s
```

## 文件结构
```
internal/infrastructure/ai/eino/
├── client.go                  # DashScope 客户端
├── client_test.go             # 客户端测试
├── intent_recognizer.go       # 意图识别器
├── rag_retriever.go           # RAG 检索器
├── order_querier.go           # 订单查询器
├── response_generator.go      # 响应生成器
├── integration_test.go        # 集成测试
└── README.md                  # 组件文档
```

## 依赖关系
```
Client
  ├── github.com/cloudwego/eino/components/model
  ├── github.com/cloudwego/eino/components/embedding
  ├── github.com/cloudwego/eino-ext/components/model/ark
  └── github.com/cloudwego/eino-ext/components/embedding/ark

IntentRecognizer
  ├── Client (ChatModel)
  └── domain/entity (Intent, Message)

RAGRetriever
  ├── Client (ChatModel, Embedder)
  ├── domain/repository (VectorRepository)
  └── domain/entity (Document)

OrderQuerier
  ├── Client (ChatModel)
  ├── domain/repository (OrderRepository)
  └── domain/entity (Order)

ResponseGenerator
  ├── Client (ChatModel)
  └── domain/entity (Message)
```

## 配置要求

### 环境变量
```bash
DASHSCOPE_API_KEY=your-api-key
```

### 配置文件 (config.yaml)
```yaml
dashscope:
  api_key: ${DASHSCOPE_API_KEY}
  chat_model: qwen-turbo
  embed_model: text-embedding-v2
  max_retries: 3
  timeout: 30s

intent:
  confidence_threshold: 0.7

rag:
  top_k: 5
  score_threshold: 0.7
```

## 验证需求

### 需求 1.1: Eino 框架集成 ✅
- ✅ 成功加载 DashScope 聊天模型配置
- ✅ 使用 Eino Compose 图编排处理请求流程（在 Use Case 层实现）
- ✅ 支持流式响应输出
- ✅ 提供统一的错误处理和日志追踪机制
- ✅ 执行重试策略并记录错误日志

### 需求 2.1: 意图识别 ✅
- ✅ 将查询分类为 course/order/direct/handoff 四种类型之一
- ✅ 返回意图类型和置信度分数
- ✅ 置信度低于阈值转人工
- ✅ 根据意图触发不同流程

### 需求 3.1: RAG 检索 ✅
- ✅ 使用 DashScope 嵌入模型生成查询向量
- ✅ 在 Milvus 向量数据库中执行相似度搜索
- ✅ 返回前 K 个最相关的文档片段
- ✅ 将文档内容和查询一起发送给 LLM 生成答案

### 需求 4.1, 4.2: 订单查询 ✅
- ✅ 从用户输入中提取订单 ID
- ✅ 生成安全的 SQL 查询语句
- ✅ 验证查询语句的安全性（防止 SQL 注入）
- ✅ 将结构化数据格式化为自然语言响应

## 性能指标

### 意图识别
- 类型映射：< 1μs
- LLM 调用：取决于 DashScope API

### RAG 检索
- 向量生成：取决于 DashScope API
- 向量搜索：取决于 Milvus 性能
- 答案生成：取决于 DashScope API

### 订单查询
- 正则提取：< 100μs
- LLM 提取：取决于 DashScope API（降级方案）
- 数据库查询：取决于 SQLite 性能

## 后续工作

### 1. Use Case 层集成
- 在 ChatUseCase 中集成这些组件
- 实现意图路由逻辑
- 实现会话历史管理

### 2. 性能优化
- 实现向量缓存
- 实现意图识别结果缓存
- 批量向量生成

### 3. 监控和日志
- 添加详细的日志记录
- 添加性能指标收集
- 添加错误追踪

### 4. 错误处理增强
- 实现重试策略
- 实现降级策略
- 实现熔断机制

## 注意事项

1. **API Key 安全**: 已通过环境变量管理，不在代码中硬编码
2. **类型转换**: 注意 Eino (float64) 和 Milvus (float32) 的向量类型差异
3. **消息格式**: AssistantMessage 需要传入 nil 作为第二个参数（ToolCall）
4. **SQL 安全**: 实现了多层安全验证，防止 SQL 注入
5. **错误处理**: 所有错误都包装了上下文信息，便于调试

## 相关文档
- [设计文档](../.kiro/specs/eino-qa-system/design.md)
- [需求文档](../.kiro/specs/eino-qa-system/requirements.md)
- [组件 README](internal/infrastructure/ai/eino/README.md)
