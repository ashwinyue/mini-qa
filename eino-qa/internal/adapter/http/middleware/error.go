package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
)

// ErrorResponse 统一错误响应格式
type ErrorResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
	TraceID string      `json:"trace_id,omitempty"`
}

// ErrorHandler 错误处理中间件
// 捕获并统一处理所有错误，返回标准化的错误响应
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 使用 defer 捕获 panic
		defer func() {
			if err := recover(); err != nil {
				// 记录 panic 堆栈
				stack := debug.Stack()

				// 获取 trace_id
				traceID := getTraceID(c)

				// 记录错误日志
				c.Error(fmt.Errorf("panic recovered: %v\nstack: %s", err, stack))

				// 返回 500 错误
				c.JSON(http.StatusInternalServerError, ErrorResponse{
					Code:    http.StatusInternalServerError,
					Message: "Internal server error",
					TraceID: traceID,
				})

				c.Abort()
			}
		}()

		// 处理请求
		c.Next()

		// 检查是否有错误
		if len(c.Errors) > 0 {
			// 获取最后一个错误
			err := c.Errors.Last().Err

			// 获取 trace_id
			traceID := getTraceID(c)

			// 确定状态码和错误消息
			statusCode, message, details := determineErrorResponse(err)

			// 如果还没有设置状态码，设置状态码
			if c.Writer.Status() == http.StatusOK {
				c.Status(statusCode)
			} else {
				statusCode = c.Writer.Status()
			}

			// 返回错误响应
			c.JSON(statusCode, ErrorResponse{
				Code:    statusCode,
				Message: message,
				Details: details,
				TraceID: traceID,
			})

			c.Abort()
		}
	}
}

// getTraceID 从 context 中获取 trace_id
func getTraceID(c *gin.Context) string {
	if traceID, exists := c.Get("trace_id"); exists {
		if id, ok := traceID.(string); ok {
			return id
		}
	}
	if requestID, exists := c.Get("request_id"); exists {
		if id, ok := requestID.(string); ok {
			return id
		}
	}
	return ""
}

// determineErrorResponse 根据错误类型确定响应
func determineErrorResponse(err error) (statusCode int, message string, details interface{}) {
	// 默认为 500 错误
	statusCode = http.StatusInternalServerError
	message = "Internal server error"
	details = nil

	// 检查是否是自定义错误类型
	var validationErr *ValidationError
	var notFoundErr *NotFoundError
	var unauthorizedErr *UnauthorizedError
	var forbiddenErr *ForbiddenError
	var badRequestErr *BadRequestError
	var serviceErr *ServiceError

	switch {
	case errors.As(err, &validationErr):
		statusCode = http.StatusBadRequest
		message = validationErr.Message
		details = validationErr.Fields
	case errors.As(err, &notFoundErr):
		statusCode = http.StatusNotFound
		message = notFoundErr.Message
	case errors.As(err, &unauthorizedErr):
		statusCode = http.StatusUnauthorized
		message = unauthorizedErr.Message
	case errors.As(err, &forbiddenErr):
		statusCode = http.StatusForbidden
		message = forbiddenErr.Message
	case errors.As(err, &badRequestErr):
		statusCode = http.StatusBadRequest
		message = badRequestErr.Message
	case errors.As(err, &serviceErr):
		statusCode = http.StatusBadGateway
		message = serviceErr.Message
		details = serviceErr.Service
	default:
		// 通用错误处理
		message = err.Error()
	}

	return statusCode, message, details
}

// 自定义错误类型

// ValidationError 验证错误
type ValidationError struct {
	Message string
	Fields  map[string]string
}

func (e *ValidationError) Error() string {
	return e.Message
}

// NewValidationError 创建验证错误
func NewValidationError(message string, fields map[string]string) *ValidationError {
	return &ValidationError{
		Message: message,
		Fields:  fields,
	}
}

// NotFoundError 资源不存在错误
type NotFoundError struct {
	Message string
}

func (e *NotFoundError) Error() string {
	return e.Message
}

// NewNotFoundError 创建资源不存在错误
func NewNotFoundError(message string) *NotFoundError {
	return &NotFoundError{Message: message}
}

// UnauthorizedError 未授权错误
type UnauthorizedError struct {
	Message string
}

func (e *UnauthorizedError) Error() string {
	return e.Message
}

// NewUnauthorizedError 创建未授权错误
func NewUnauthorizedError(message string) *UnauthorizedError {
	return &UnauthorizedError{Message: message}
}

// ForbiddenError 禁止访问错误
type ForbiddenError struct {
	Message string
}

func (e *ForbiddenError) Error() string {
	return e.Message
}

// NewForbiddenError 创建禁止访问错误
func NewForbiddenError(message string) *ForbiddenError {
	return &ForbiddenError{Message: message}
}

// BadRequestError 错误请求错误
type BadRequestError struct {
	Message string
}

func (e *BadRequestError) Error() string {
	return e.Message
}

// NewBadRequestError 创建错误请求错误
func NewBadRequestError(message string) *BadRequestError {
	return &BadRequestError{Message: message}
}

// ServiceError 外部服务错误
type ServiceError struct {
	Message string
	Service string
}

func (e *ServiceError) Error() string {
	return e.Message
}

// NewServiceError 创建外部服务错误
func NewServiceError(message string, service string) *ServiceError {
	return &ServiceError{
		Message: message,
		Service: service,
	}
}
