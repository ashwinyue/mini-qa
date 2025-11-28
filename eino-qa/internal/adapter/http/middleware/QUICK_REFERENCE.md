# HTTP 中间件快速参考

## 快速开始

```go
import (
    "eino-qa/internal/adapter/http/middleware"
    "github.com/gin-gonic/gin"
)

func setupRouter(config *Config, logger Logger) *gin.Engine {
    router := gin.New()
    
    // 基础中间件（推荐顺序）
    router.Use(middleware.ErrorHandler())
    router.Use(middleware.NewLoggingMiddleware(logger).Handler())
    router.Use(middleware.TenantMiddleware())
    router.Use(middleware.NewSecurityMiddleware(config.Security.SensitiveFields).Handler())
    
    // 公开路由
    router.GET("/health", healthHandler)
    
    // 需要认证的路由
    api := router.Group("/api/v1")
    api.Use(middleware.NewAuthMiddleware(config.Security.APIKeys).Handler())
    {
        api.POST("/vectors/items", addVectorHandler)
        api.DELETE("/vectors/items", deleteVectorHandler)
    }
    
    return router
}
```

## 中间件速查表

| 中间件 | 用途 | 位置 | 必需 |
|--------|------|------|------|
| ErrorHandler | 统一错误处理 | 最外层 | ✓ |
| LoggingMiddleware | 请求日志 | 第二层 | ✓ |
| TenantMiddleware | 租户识别 | 第三层 | ✓ |
| SecurityMiddleware | 敏感信息脱敏 | 第四层 | ✓ |
| AuthMiddleware | API Key 验证 | 特定路由 | 可选 |

## 常见用法

### 1. 获取租户 ID

```go
func handler(c *gin.Context) {
    tenantID, _ := c.Get("tenant_id")
    // 或从 context 获取
    tenantID := c.Request.Context().Value("tenant_id")
}
```

### 2. 返回自定义错误

```go
func handler(c *gin.Context) {
    if err := validate(data); err != nil {
        c.Error(middleware.NewValidationError("Invalid data", fields))
        return
    }
}
```

### 3. 脱敏日志内容

```go
securityMW := middleware.NewSecurityMiddleware([]string{"password"})
sanitized := securityMW.SanitizeForLog(sensitiveData)
logger.Info("Data: " + sanitized)
```

### 4. 可选认证

```go
// 提供了 API Key 则验证，否则允许通过
optionalAuth := middleware.NewOptionalAuthMiddleware(apiKeys)
router.Use(optionalAuth.Handler())
```

## 错误类型

```go
// 400 Bad Request
middleware.NewValidationError("message", fields)
middleware.NewBadRequestError("message")

// 401 Unauthorized
middleware.NewUnauthorizedError("message")

// 403 Forbidden
middleware.NewForbiddenError("message")

// 404 Not Found
middleware.NewNotFoundError("message")

// 502 Bad Gateway
middleware.NewServiceError("message", "service_name")
```

## Context 键值

| 键 | 类型 | 设置者 | 用途 |
|----|------|--------|------|
| `tenant_id` | string | TenantMiddleware | 租户标识 |
| `request_id` | string | LoggingMiddleware | 请求追踪 |
| `trace_id` | string | 用户设置 | 分布式追踪 |
| `sanitized_request_body` | string | SecurityMiddleware | 脱敏后的请求体 |

## 配置示例

```yaml
security:
  api_keys:
    - "key1"
    - "key2"
  sensitive_fields:
    - "password"
    - "secret"
    - "token"

logging:
  level: "info"
  format: "json"
  output: "stdout"
```

## 测试示例

```go
func TestHandler(t *testing.T) {
    router := gin.New()
    router.Use(middleware.TenantMiddleware())
    router.GET("/test", handler)
    
    req := httptest.NewRequest("GET", "/test", nil)
    req.Header.Set("X-Tenant-ID", "tenant1")
    
    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)
    
    assert.Equal(t, 200, w.Code)
}
```

## 性能提示

1. **日志中间件**: 避免记录大型响应体
2. **安全中间件**: 仅在需要时启用脱敏
3. **认证中间件**: 使用 map 而非 slice 存储 API Key
4. **错误处理**: 避免在 panic 恢复中执行耗时操作

## 故障排查

### 问题：租户 ID 为空
- 检查请求头 `X-Tenant-ID` 或查询参数 `tenant`
- 确认 TenantMiddleware 已注册

### 问题：API Key 验证失败
- 检查请求头 `X-API-Key` 或 `Authorization: Bearer <token>`
- 确认 API Key 在配置列表中

### 问题：日志中包含敏感信息
- 确认 SecurityMiddleware 在 LoggingMiddleware 之后
- 检查 `sensitive_fields` 配置

### 问题：错误响应格式不正确
- 确认 ErrorHandler 在最外层
- 使用 `c.Error()` 而非直接 `c.JSON()`
