package errors

import (
	"errors"
	"fmt"
)

// ErrorCategory 错误类别
type ErrorCategory string

const (
	// CategoryClient 客户端错误
	CategoryClient ErrorCategory = "client"
	// CategoryServer 服务端错误
	CategoryServer ErrorCategory = "server"
	// CategoryExternal 外部服务错误
	CategoryExternal ErrorCategory = "external"
	// CategoryValidation 验证错误
	CategoryValidation ErrorCategory = "validation"
)

// 预定义错误类型
var (
	// ErrNotFound 资源不存在
	ErrNotFound = errors.New("resource not found")

	// ErrInvalidInput 无效输入
	ErrInvalidInput = errors.New("invalid input")

	// ErrUnauthorized 未授权
	ErrUnauthorized = errors.New("unauthorized")

	// ErrInternal 内部错误
	ErrInternal = errors.New("internal error")

	// ErrServiceUnavailable 服务不可用
	ErrServiceUnavailable = errors.New("service unavailable")

	// ErrTimeout 超时
	ErrTimeout = errors.New("timeout")

	// ErrRateLimited 请求频率超限
	ErrRateLimited = errors.New("rate limited")

	// ErrDashScopeUnavailable DashScope 服务不可用
	ErrDashScopeUnavailable = errors.New("dashscope service unavailable")

	// ErrMilvusUnavailable Milvus 服务不可用
	ErrMilvusUnavailable = errors.New("milvus service unavailable")

	// ErrDatabaseUnavailable 数据库不可用
	ErrDatabaseUnavailable = errors.New("database unavailable")

	// ErrEmbeddingFailed 向量生成失败
	ErrEmbeddingFailed = errors.New("embedding generation failed")

	// ErrVectorSearchFailed 向量搜索失败
	ErrVectorSearchFailed = errors.New("vector search failed")

	// ErrLLMCallFailed LLM 调用失败
	ErrLLMCallFailed = errors.New("llm call failed")
)

// AppError 应用错误
type AppError struct {
	Code      int
	Message   string
	Err       error
	Category  ErrorCategory
	Retryable bool
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.Err
}

// New 创建新的应用错误
func New(code int, message string, err error) *AppError {
	return &AppError{
		Code:      code,
		Message:   message,
		Err:       err,
		Category:  CategoryServer,
		Retryable: false,
	}
}

// NewWithCategory 创建带类别的应用错误
func NewWithCategory(code int, message string, err error, category ErrorCategory, retryable bool) *AppError {
	return &AppError{
		Code:      code,
		Message:   message,
		Err:       err,
		Category:  category,
		Retryable: retryable,
	}
}

// Wrap 包装错误
func Wrap(err error, message string) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s: %w", message, err)
}

// IsRetryable 判断错误是否可重试
func IsRetryable(err error) bool {
	if err == nil {
		return false
	}

	// 检查是否是 AppError 且标记为可重试
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.Retryable
	}

	// 检查是否是超时或服务不可用错误
	return errors.Is(err, ErrTimeout) ||
		errors.Is(err, ErrServiceUnavailable) ||
		errors.Is(err, ErrDashScopeUnavailable) ||
		errors.Is(err, ErrMilvusUnavailable) ||
		errors.Is(err, ErrDatabaseUnavailable) ||
		errors.Is(err, ErrVectorSearchFailed) ||
		errors.Is(err, ErrLLMCallFailed)
}

// IsExternalError 判断是否是外部服务错误
func IsExternalError(err error) bool {
	if err == nil {
		return false
	}

	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.Category == CategoryExternal
	}

	return errors.Is(err, ErrDashScopeUnavailable) ||
		errors.Is(err, ErrMilvusUnavailable) ||
		errors.Is(err, ErrDatabaseUnavailable)
}

// GetErrorCategory 获取错误类别
func GetErrorCategory(err error) ErrorCategory {
	if err == nil {
		return CategoryServer
	}

	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.Category
	}

	return CategoryServer
}
