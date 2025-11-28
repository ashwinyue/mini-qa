package errors

import (
	"errors"
	"fmt"
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
)

// AppError 应用错误
type AppError struct {
	Code    int
	Message string
	Err     error
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
		Code:    code,
		Message: message,
		Err:     err,
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

	// 检查是否是超时或服务不可用错误
	return errors.Is(err, ErrTimeout) || errors.Is(err, ErrServiceUnavailable)
}
