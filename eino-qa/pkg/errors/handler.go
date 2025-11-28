package errors

import (
	"context"
	"fmt"
	"time"
)

// ErrorHandler 错误处理器
type ErrorHandler struct {
	retryConfig    RetryConfig
	circuitBreaker *CircuitBreaker
	fallbackFn     FallbackFunc
}

// NewErrorHandler 创建错误处理器
func NewErrorHandler(retryConfig RetryConfig) *ErrorHandler {
	return &ErrorHandler{
		retryConfig: retryConfig,
	}
}

// WithCircuitBreaker 设置熔断器
func (eh *ErrorHandler) WithCircuitBreaker(maxFailures int, resetTimeout time.Duration) *ErrorHandler {
	eh.circuitBreaker = NewCircuitBreaker(maxFailures, resetTimeout)
	return eh
}

// WithFallback 设置降级函数
func (eh *ErrorHandler) WithFallback(fallback FallbackFunc) *ErrorHandler {
	eh.fallbackFn = fallback
	return eh
}

// Execute 执行函数（带重试、熔断、降级）
func (eh *ErrorHandler) Execute(ctx context.Context, fn RetryableFunc) error {
	// 如果有熔断器，先检查熔断器状态
	if eh.circuitBreaker != nil {
		if eh.circuitBreaker.GetState() == StateOpen {
			// 熔断器打开，直接执行降级
			if eh.fallbackFn != nil {
				return eh.fallbackFn()
			}
			return fmt.Errorf("circuit breaker is open")
		}
	}

	// 执行重试
	err := RetryWithBackoff(ctx, eh.retryConfig, fn)

	// 如果有熔断器，记录结果
	if eh.circuitBreaker != nil {
		if err != nil {
			eh.circuitBreaker.recordFailure()
		} else {
			eh.circuitBreaker.recordSuccess()
		}
	}

	// 如果失败且有降级函数，执行降级
	if err != nil && eh.fallbackFn != nil {
		return eh.fallbackFn()
	}

	return err
}

// HandleWithRetry 处理函数（仅重试）
func HandleWithRetry(ctx context.Context, config RetryConfig, fn RetryableFunc) error {
	return RetryWithBackoff(ctx, config, fn)
}

// HandleWithFallback 处理函数（仅降级）
func HandleWithFallback(fn FallbackFunc, fallback FallbackFunc) error {
	err := fn()
	if err != nil && fallback != nil {
		return fallback()
	}
	return err
}

// HandleWithCircuitBreaker 处理函数（仅熔断）
func HandleWithCircuitBreaker(cb *CircuitBreaker, fn FallbackFunc) error {
	return cb.Execute(fn)
}
