# 错误处理和重试机制

本包提供了统一的错误处理、重试策略和降级机制。

## 功能特性

1. **统一错误类型定义** - 标准化的错误类型和分类
2. **指数退避重试** - 自动重试失败的操作
3. **熔断器** - 防止级联故障
4. **降级策略** - 服务不可用时的备选方案

## 使用示例

### 1. 基本错误处理

```go
import "eino-qa/pkg/errors"

// 创建应用错误
err := errors.NewWithCategory(
    500,
    "DashScope 服务调用失败",
    originalErr,
    errors.CategoryExternal,
    true, // 可重试
)

// 判断错误是否可重试
if errors.IsRetryable(err) {
    // 执行重试逻辑
}

// 判断是否是外部服务错误
if errors.IsExternalError(err) {
    // 记录外部服务错误
}
```

### 2. 重试策略

#### 基本重试

```go
import (
    "context"
    "time"
    "eino-qa/pkg/errors"
)

// 使用默认配置重试
err := errors.RetryWithBackoff(
    context.Background(),
    errors.DefaultRetryConfig,
    func() error {
        // 可能失败的操作
        return callExternalService()
    },
)

// 自定义重试配置
config := errors.RetryConfig{
    MaxAttempts:  5,
    InitialDelay: 200 * time.Millisecond,
    MaxDelay:     10 * time.Second,
    Multiplier:   2.0,
    ShouldRetry:  errors.IsRetryable,
}

err = errors.RetryWithBackoff(ctx, config, func() error {
    return callDashScope()
})
```

#### 带抖动的重试

```go
// 使用抖动避免惊群效应
err := errors.RetryWithJitter(
    ctx,
    errors.DefaultRetryConfig,
    func() error {
        return searchMilvus()
    },
)
```

### 3. 降级策略

#### 简单降级

```go
// 主函数失败时执行降级函数
err := errors.HandleWithFallback(
    func() error {
        // 主要逻辑：调用 LLM
        return callLLM()
    },
    func() error {
        // 降级逻辑：返回预设回复
        return returnDefaultResponse()
    },
)
```

#### 降级链

```go
chain := errors.NewFallbackChain()

// 添加多个降级处理器
chain.Add("primary", func() error {
    return callPrimaryService()
}, func() error {
    return callBackupService()
})

chain.Add("secondary", func() error {
    return callSecondaryService()
}, func() error {
    return returnCachedResult()
})

// 执行降级链
err := chain.Execute(ctx)
```

### 4. 熔断器

```go
// 创建熔断器（5次失败后熔断，30秒后尝试恢复）
cb := errors.NewCircuitBreaker(5, 30*time.Second)

// 执行函数（带熔断保护）
err := cb.Execute(func() error {
    return callUnstableService()
})

// 检查熔断器状态
state := cb.GetState()
switch state {
case errors.StateClosed:
    // 正常状态
case errors.StateOpen:
    // 熔断状态
case errors.StateHalfOpen:
    // 尝试恢复状态
}
```

### 5. 综合使用

```go
// 创建错误处理器（集成重试、熔断、降级）
handler := errors.NewErrorHandler(errors.RetryConfig{
    MaxAttempts:  3,
    InitialDelay: 100 * time.Millisecond,
    MaxDelay:     5 * time.Second,
    Multiplier:   2.0,
}).WithCircuitBreaker(5, 30*time.Second).
   WithFallback(func() error {
       return returnDefaultResponse()
   })

// 执行函数
err := handler.Execute(ctx, func() error {
    return callExternalAPI()
})
```

## 实际应用场景

### 场景 1: LLM 调用失败

```go
// LLM 调用失败时的处理
err := errors.HandleWithRetry(ctx, errors.DefaultRetryConfig, func() error {
    resp, err := chatModel.Generate(ctx, messages)
    if err != nil {
        return errors.NewWithCategory(
            502,
            "LLM 调用失败",
            err,
            errors.CategoryExternal,
            true,
        )
    }
    return nil
})

if err != nil {
    // 降级：返回预设回复
    return "抱歉，我现在无法回答您的问题，请稍后再试或联系人工客服。"
}
```

### 场景 2: Milvus 向量搜索失败

```go
// Milvus 不可用时降级到关键词匹配
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
// 数据库查询失败时的处理
var order *Order
err := errors.RetryWithBackoff(ctx, errors.RetryConfig{
    MaxAttempts:  3,
    InitialDelay: 50 * time.Millisecond,
    MaxDelay:     1 * time.Second,
    Multiplier:   2.0,
}, func() error {
    return db.Where("id = ?", orderID).First(&order).Error
})

if err != nil {
    return nil, errors.NewWithCategory(
        500,
        "订单查询失败",
        err,
        errors.CategoryServer,
        false,
    )
}
```

## 最佳实践

1. **明确错误类别** - 使用 `CategoryExternal` 标记外部服务错误
2. **合理设置重试次数** - 避免过度重试导致级联故障
3. **使用抖动** - 在高并发场景下使用 `RetryWithJitter` 避免惊群
4. **设置超时** - 始终使用带超时的 context
5. **记录日志** - 在重试和降级时记录详细日志
6. **监控指标** - 统计重试次数、熔断次数、降级次数

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
