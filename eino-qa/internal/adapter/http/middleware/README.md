# HTTP 中间件

本目录包含了系统的所有 HTTP 中间件实现，用于处理跨切面关注点。

## 中间件列表

### 1. TenantMiddleware - 租户识别中间件

**功能**: 从请求中提取租户 ID 并设置到 context 中

**提取顺序**:
1. 请求头 `X-Tenant-ID`
2. 查询参数 `tenant`
3. 默认值 `default`

**使用示例**:
```go
router := gin.Default()
router.Use(middleware.TenantMiddleware())
```

**验证需求**: 5.1, 5.2

### 2. SecurityMiddleware - 安全脱敏中间件

**功能**: 对请求和响应中的敏感信息进行脱敏处理

**脱敏内容**:
- 密码字段（password, pwd, passwd）
- 身份证号（18位）
- 手机号（11位）
- 邮箱地址
- 信用卡号
- API Key 和 Token

**使用示例**:
```go
securityMW := middleware.NewSecurityMiddleware([]string{"password", "id_card", "phone"})
router.Use(securityMW.Handler())
```

**验证需求**: 8.1, 8.2

### 3. LoggingMiddleware - 日志记录中间件

**功能**: 记录请求的详细信息

**记录内容**:
- 请求 ID（自动生成 UUID）
- 租户 ID
- 请求方法和路径
- 查询参数
- 客户端 IP
- 处理时长
- 响应状态码
- 请求和响应大小

**使用示例**:
```go
loggingMW := middleware.NewLoggingMiddleware(logger)
router.Use(loggingMW.Handler())
```

**验证需求**: 7.1

### 4. ErrorHandler - 错误处理中间件

**功能**: 统一处理所有错误并返回标准化响应

**错误类型**:
- `ValidationError` - 400 验证错误
- `NotFoundError` - 404 资源不存在
- `UnauthorizedError` - 401 未授权
- `ForbiddenError` - 403 禁止访问
- `BadRequestError` - 400 错误请求
- `ServiceError` - 502 外部服务错误

**响应格式**:
```json
{
  "code": 400,
  "message": "Validation failed",
  "details": {
    "field": "error message"
  },
  "trace_id": "uuid"
}
```

**使用示例**:
```go
router.Use(middleware.ErrorHandler())
```

**验证需求**: 1.4, 1.5

### 5. AuthMiddleware - API Key 验证中间件

**功能**: 验证请求头中的 API Key

**验证方式**:
1. 请求头 `X-API-Key`
2. Authorization 头 `Bearer <token>`

**使用示例**:
```go
// 强制验证
authMW := middleware.NewAuthMiddleware([]string{"key1", "key2"})
router.Use(authMW.Handler())

// 可选验证
optionalAuthMW := middleware.NewOptionalAuthMiddleware([]string{"key1", "key2"})
router.Use(optionalAuthMW.Handler())
```

**验证需求**: 8.3, 8.4, 9.1

## 中间件使用顺序

推荐的中间件使用顺序（从上到下）：

```go
router := gin.Default()

// 1. 错误恢复（最外层）
router.Use(middleware.ErrorHandler())

// 2. 日志记录
loggingMW := middleware.NewLoggingMiddleware(logger)
router.Use(loggingMW.Handler())

// 3. 租户识别
router.Use(middleware.TenantMiddleware())

// 4. 安全脱敏
securityMW := middleware.NewSecurityMiddleware(config.Security.SensitiveFields)
router.Use(securityMW.Handler())

// 5. API Key 验证（针对特定路由组）
apiGroup := router.Group("/api/v1")
authMW := middleware.NewAuthMiddleware(config.Security.APIKeys)
apiGroup.Use(authMW.Handler())
```

## 测试

每个中间件都应该有对应的单元测试，测试以下场景：

1. **正常流程**: 中间件正确处理请求
2. **边界情况**: 缺少必需参数、空值等
3. **错误处理**: 验证失败、异常情况等
4. **性能**: 确保中间件不会显著影响性能

## 注意事项

1. **中间件顺序很重要**: 错误处理应该在最外层，日志记录应该在租户识别之后
2. **Context 传递**: 使用 `c.Set()` 和 `c.Get()` 在中间件之间传递数据
3. **性能考虑**: 避免在中间件中进行耗时操作
4. **错误处理**: 使用 `c.Error()` 记录错误，由 ErrorHandler 统一处理
5. **安全性**: 敏感信息必须脱敏后才能记录到日志
