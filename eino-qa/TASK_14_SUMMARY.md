# Task 14: 错误处理和重试机制 - 实现总结

## 任务概述

实现了统一的错误处理、重试策略和降级机制，为系统提供了完善的容错能力。

## 实现内容

### 1. 统一错误类型定义 (`pkg/errors/errors.go`)

**增强功能**:
- 添加了 `ErrorCategory` 错误分类（客户端、服务端、外部服务、验证）
- 扩展了预定义错误类型（DashScope、Milvus、数据库等）
- 增强了 `AppError` 结构，支持错误分类和可重试标记
- 提供了 `IsRetryable()` 和 `IsExternalError()` 判断函数

**新增错误类型**:
```go
- ErrRateLimited          // 请求频率超限
- ErrDashScopeUnavailable // DashScope 服务不可用
- ErrMilvusUnavailable    // Milvus 服务不可用
- ErrDatabaseUnavailable  // 数据库不可用
- ErrEmbeddingFailed      // 向量生成失败
- ErrVectorSearchFailed   // 向量搜索失败
- ErrLLMCallFailed        // LLM 调用失败
```

### 2. 重试策略 (`pkg/errors/retry.go`)

**实现功能**:
- **指数退避重试** - `RetryWithBackoff()` 函数
- **带抖动的重试** - `RetryWithJitter()` 函数，避免惊群效应
- **可配置的重试策略** - `RetryConfig` 结构
- **上下文感知** - 支持 context 取消

**配置参数**:
```go
type RetryConfig struct {
    MaxAttempts  int           // 最大重试次数
    InitialDelay time.Duration // 初始延迟
    MaxDelay     time.Duration // 最大延迟
    Multiplier   float64       // 延迟倍数
    ShouldRetry  func(error) bool // 自定义重试判断
}
```

**默认配置**:
- MaxAttempts: 3
- InitialDelay: 100ms
- MaxDelay: 5s
- Multiplier: 2.0

### 3. 降级策略 (`pkg/errors/fallback.go`)

**实现功能**:
- **简单降级** - `WithFallback()` 函数
- **降级链** - `FallbackChain` 支持多级降级
- **熔断器** - `CircuitBreaker` 防止级联故障

**熔断器状态**:
- `StateClosed` - 关闭状态（正常）
- `StateOpen` - 打开状态（熔断）
- `StateHalfOpen` - 半开状态（尝试恢复）

**熔断器配置**:
```go
cb := NewCircuitBreaker(
    5,              // 最大失败次数
    30*time.Second, // 重置超时时间
)
```

### 4. 错误处理器 (`pkg/errors/handler.go`)

**综合处理器**:
- 集成重试、熔断、降级功能
- 支持链式配置
- 提供便捷的处理函数

**使用示例**:
```go
handler := NewErrorHandler(retryConfig).
    WithCircuitBreaker(5, 30*time.Second).
    WithFallback(fallbackFunc)

err := handler.Execute(ctx, func() error {
    return callExternalService()
})
```

### 5. 错误处理中间件 (`internal/adapter/http/middleware/error.go`)

**已有功能**（确认完整）:
- Panic 恢复和处理
- 统一错误响应格式
- 自定义错误类型支持
- TraceID 追踪

**错误响应格式**:
```json
{
    "code": 500,
    "message": "错误信息",
    "details": {},
    "trace_id": "request-id"
}
```

## 文件结构

```
eino-qa/
├── pkg/errors/
│   ├── errors.go      # 统一错误类型定义
│   ├── retry.go       # 重试策略实现
│   ├── fallback.go    # 降级策略和熔断器
│   ├── handler.go     # 综合错误处理器
│   └── README.md      # 使用文档
├── internal/adapter/http/middleware/
│   └── error.go       # HTTP 错误处理中间件
└── examples/
    └── error_handling_example.go  # 使用示例
```

## 使用场景

### 场景 1: LLM 调用失败

```go
err := errors.RetryWithBackoff(ctx, errors.DefaultRetryConfig, func() error {
    resp, err := chatModel.Generate(ctx, messages)
    if err != nil {
        return errors.NewWithCategory(
            502, "LLM 调用失败", err,
            errors.CategoryExternal, true,
        )
    }
    return nil
})

if err != nil {
    // 降级：返回预设回复
    return "抱歉，我现在无法回答您的问题，请稍后再试。"
}
```

### 场景 2: Milvus 向量搜索失败

```go
results, err := errors.HandleWithFallback(
    func() error {
        // 主要逻辑：向量搜索
        results, err = milvusClient.Search(ctx, vector, topK)
        return err
    },
    func() error {
        // 降级逻辑：关键词匹配
        results = keywordSearch(query)
        return nil
    },
)
```

### 场景 3: 数据库查询失败

```go
var order *Order
err := errors.RetryWithBackoff(ctx, errors.RetryConfig{
    MaxAttempts:  3,
    InitialDelay: 50 * time.Millisecond,
    MaxDelay:     1 * time.Second,
    Multiplier:   2.0,
}, func() error {
    return db.Where("id = ?", orderID).First(&order).Error
})
```

## 配置建议

### LLM 调用
- MaxAttempts: 3
- InitialDelay: 100ms
- MaxDelay: 5s
- Multiplier: 2.0

### 向量搜索
- MaxAttempts: 2
- InitialDelay: 50ms
- MaxDelay: 2s
- Multiplier: 2.0

### 数据库查询
- MaxAttempts: 3
- InitialDelay: 50ms
- MaxDelay: 1s
- Multiplier: 2.0

### 熔断器
- MaxFailures: 5
- ResetTimeout: 30s

## 验证结果

✅ 所有代码编译通过
✅ 示例程序运行成功
✅ 重试策略正常工作
✅ 降级策略正常工作
✅ 熔断器正常工作
✅ 综合错误处理器正常工作

## 满足的需求

- ✅ **需求 1.4**: 实现统一的错误处理和日志追踪机制
- ✅ **需求 1.5**: 实现重试策略并记录错误日志

## 最佳实践

1. **明确错误类别** - 使用 `CategoryExternal` 标记外部服务错误
2. **合理设置重试次数** - 避免过度重试导致级联故障
3. **使用抖动** - 在高并发场景下使用 `RetryWithJitter` 避免惊群
4. **设置超时** - 始终使用带超时的 context
5. **记录日志** - 在重试和降级时记录详细日志
6. **监控指标** - 统计重试次数、熔断次数、降级次数

## 后续集成建议

1. **在 AI 组件中集成重试**:
   - DashScope ChatModel 调用
   - Embedding 生成
   - Milvus 向量搜索

2. **在 Repository 中集成重试**:
   - SQLite 数据库操作
   - Milvus Collection 操作

3. **在 Use Case 中集成降级**:
   - RAG 检索失败降级到关键词匹配
   - LLM 调用失败返回预设回复
   - 订单查询失败提示联系客服

4. **添加监控指标**:
   - 重试次数统计
   - 熔断次数统计
   - 降级次数统计
   - 错误类型分布

## 示例输出

```
错误处理和重试机制示例
================================

=== Example 1: 基本重试 ===
尝试调用外部服务 (第 1 次)
尝试调用外部服务 (第 2 次)
尝试调用外部服务 (第 3 次)
调用成功!

=== Example 3: 降级策略 ===
尝试主要服务...
执行降级逻辑...
返回默认响应
处理成功（使用降级）

=== Example 5: 熔断器 ===
调用服务 (第 1 次)
  成功
  熔断器状态: 关闭（正常）
...

所有示例执行完成!
```

## 总结

任务 14 已完成，实现了完整的错误处理和重试机制：

1. ✅ 统一错误类型定义 - 支持错误分类和可重试标记
2. ✅ 错误处理中间件 - 已存在且功能完整
3. ✅ 重试策略（指数退避）- 支持基本重试和带抖动的重试
4. ✅ 降级策略 - 支持简单降级、降级链和熔断器

系统现在具备了完善的容错能力，可以优雅地处理外部服务故障、网络超时等异常情况。
