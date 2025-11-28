package errors

import (
	"context"
	"fmt"
	"time"
)

// FallbackFunc 降级函数类型
type FallbackFunc func() error

// FallbackChain 降级链
type FallbackChain struct {
	handlers []FallbackHandler
}

// FallbackHandler 降级处理器
type FallbackHandler struct {
	Name     string
	Fn       FallbackFunc
	Fallback FallbackFunc
}

// NewFallbackChain 创建降级链
func NewFallbackChain() *FallbackChain {
	return &FallbackChain{
		handlers: make([]FallbackHandler, 0),
	}
}

// Add 添加降级处理器
func (fc *FallbackChain) Add(name string, fn FallbackFunc, fallback FallbackFunc) *FallbackChain {
	fc.handlers = append(fc.handlers, FallbackHandler{
		Name:     name,
		Fn:       fn,
		Fallback: fallback,
	})
	return fc
}

// Execute 执行降级链
func (fc *FallbackChain) Execute(ctx context.Context) error {
	for _, handler := range fc.handlers {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			err := handler.Fn()
			if err == nil {
				return nil
			}

			// 如果有降级函数，执行降级
			if handler.Fallback != nil {
				fallbackErr := handler.Fallback()
				if fallbackErr == nil {
					return nil
				}
				// 降级也失败，继续下一个处理器
			}
		}
	}

	return fmt.Errorf("all fallback handlers failed")
}

// WithFallback 为函数添加降级处理
func WithFallback(primary FallbackFunc, fallback FallbackFunc) FallbackFunc {
	return func() error {
		err := primary()
		if err != nil && fallback != nil {
			return fallback()
		}
		return err
	}
}

// CircuitBreaker 熔断器
type CircuitBreaker struct {
	maxFailures  int
	resetTimeout time.Duration
	failures     int
	lastFailTime time.Time
	state        CircuitState
}

// CircuitState 熔断器状态
type CircuitState int

const (
	// StateClosed 关闭状态（正常）
	StateClosed CircuitState = iota
	// StateOpen 打开状态（熔断）
	StateOpen
	// StateHalfOpen 半开状态（尝试恢复）
	StateHalfOpen
)

// NewCircuitBreaker 创建熔断器
func NewCircuitBreaker(maxFailures int, resetTimeout time.Duration) *CircuitBreaker {
	return &CircuitBreaker{
		maxFailures:  maxFailures,
		resetTimeout: resetTimeout,
		state:        StateClosed,
	}
}

// Execute 执行函数（带熔断保护）
func (cb *CircuitBreaker) Execute(fn FallbackFunc) error {
	// 检查熔断器状态
	if cb.state == StateOpen {
		// 检查是否可以尝试恢复
		if time.Since(cb.lastFailTime) > cb.resetTimeout {
			cb.state = StateHalfOpen
		} else {
			return fmt.Errorf("circuit breaker is open")
		}
	}

	// 执行函数
	err := fn()
	if err != nil {
		cb.recordFailure()
		return err
	}

	// 成功，重置失败计数
	cb.recordSuccess()
	return nil
}

// recordFailure 记录失败
func (cb *CircuitBreaker) recordFailure() {
	cb.failures++
	cb.lastFailTime = time.Now()

	if cb.failures >= cb.maxFailures {
		cb.state = StateOpen
	}
}

// recordSuccess 记录成功
func (cb *CircuitBreaker) recordSuccess() {
	cb.failures = 0
	cb.state = StateClosed
}

// GetState 获取熔断器状态
func (cb *CircuitBreaker) GetState() CircuitState {
	return cb.state
}

// Reset 重置熔断器
func (cb *CircuitBreaker) Reset() {
	cb.failures = 0
	cb.state = StateClosed
}
