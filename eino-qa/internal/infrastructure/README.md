# Infrastructure Layer

基础设施层实现了系统的外部依赖和技术细节，包括配置管理、日志系统和 AI 客户端。

## 组件

### 1. 配置管理 (config)

负责加载和验证系统配置。

#### 使用示例

```go
import "eino-qa/internal/infrastructure/config"

// 加载配置文件
cfg, err := config.Load("config/config.yaml")
if err != nil {
    log.Fatalf("failed to load config: %v", err)
}

// 访问配置
fmt.Println("Server Port:", cfg.Server.Port)
fmt.Println("Chat Model:", cfg.DashScope.ChatModel)
```

#### 配置文件格式

配置文件使用 YAML 格式，支持环境变量展开：

```yaml
server:
  port: 8080
  mode: debug

dashscope:
  api_key: ${DASHSCOPE_API_KEY}  # 从环境变量读取
  chat_model: qwen-turbo
  embed_model: text-embedding-v2
```

#### 配置验证

配置加载时会自动验证必需字段：
- Server Port 必须在 1-65535 范围内
- DashScope API Key 不能为空
- Milvus Host 不能为空
- Database BasePath 不能为空

### 2. 日志系统 (logger)

提供结构化日志记录功能，支持 JSON 和文本格式。

#### 使用示例

```go
import "eino-qa/internal/infrastructure/logger"

// 创建日志实例
lgr, err := logger.New(logger.Config{
    Level:  "info",
    Format: "json",
    Output: "stdout",
})
if err != nil {
    log.Fatalf("failed to create logger: %v", err)
}

// 创建带有上下文的日志
ctx := context.Background()
ctx = context.WithValue(ctx, "trace_id", "trace-123")
ctx = context.WithValue(ctx, "tenant_id", "tenant-456")

// 记录日志
lgr.Info(ctx, "用户请求", map[string]interface{}{
    "user_id": "user-789",
    "action":  "query",
})
```

#### 日志级别

- `debug`: 调试信息
- `info`: 一般信息
- `warn`: 警告信息
- `error`: 错误信息

#### 上下文字段

日志系统会自动从 context 中提取以下字段：
- `trace_id`: 请求追踪 ID
- `tenant_id`: 租户 ID
- `session_id`: 会话 ID
- `request_id`: 请求 ID

#### 预设字段

可以创建带有预设字段的 logger：

```go
serviceLogger := lgr.WithFields(map[string]interface{}{
    "service": "chat-service",
    "version": "1.0.0",
})

serviceLogger.Info(ctx, "服务启动", nil)
```

### 3. DashScope 客户端 (ai/eino)

封装了 Eino 框架的 DashScope 客户端，提供聊天模型和嵌入模型。

#### 使用示例

```go
import "eino-qa/internal/infrastructure/ai/eino"

// 创建客户端
client, err := eino.NewClient(eino.ClientConfig{
    APIKey:     "your-api-key",
    ChatModel:  "qwen-turbo",
    EmbedModel: "text-embedding-v2",
    MaxRetries: 3,
    Timeout:    30 * time.Second,
})
if err != nil {
    log.Fatalf("failed to create client: %v", err)
}
defer client.Close()

// 获取聊天模型
chatModel := client.GetChatModel()

// 获取嵌入模型
embedModel := client.GetEmbedModel()
```

#### 默认值

- `ChatModel`: "qwen-turbo"
- `EmbedModel`: "text-embedding-v2"
- `MaxRetries`: 3
- `Timeout`: 30 秒

#### 支持的模型

**聊天模型**:
- `qwen-turbo`: 快速响应，适合一般对话
- `qwen-plus`: 平衡性能和质量
- `qwen-max`: 最高质量，适合复杂任务

**嵌入模型**:
- `text-embedding-v2`: 通用文本嵌入模型

## 完整示例

参见 `examples/infrastructure_usage.go` 文件，展示了如何组合使用这些组件。

运行示例：

```bash
cd examples
go run infrastructure_usage.go
```

## 测试

运行所有基础设施层测试：

```bash
go test ./internal/infrastructure/... -v
```

运行特定组件测试：

```bash
go test ./internal/infrastructure/config -v
go test ./internal/infrastructure/logger -v
go test ./internal/infrastructure/ai/eino -v
```

## 依赖

- `github.com/cloudwego/eino`: Eino AI 框架
- `github.com/cloudwego/eino-ext/components/model/ark`: Ark 模型组件
- `github.com/cloudwego/eino-ext/components/embedding/ark`: Ark 嵌入组件
- `github.com/sirupsen/logrus`: 结构化日志库
- `gopkg.in/yaml.v3`: YAML 解析库

## 架构说明

基础设施层遵循 Clean Architecture 原则：
- 不依赖于业务逻辑层
- 实现领域层定义的接口
- 可以被轻松替换或模拟

## 下一步

完成基础设施层后，可以继续实现：
1. Domain Layer - 定义核心实体和仓储接口
2. Use Case Layer - 实现业务逻辑
3. Interface Adapter Layer - 实现 HTTP 接口
