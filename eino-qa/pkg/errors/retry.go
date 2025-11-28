package errors

import (
	"context"
	"fmt"
	"math"
	"time"
)

// RetryConfig 重试配置
type RetryConfig struct {
	// MaxAttempts 最大重试次数
	MaxAttempts int
	// InitialDelay 初始延迟时间
	InitialDelay time.Duration
	// MaxDelay 最大延迟时间
	MaxDelay time.Duration
	// Multiplier 延迟倍数
	Multiplier float64
	// ShouldRetry 自定义重试判断函数（可选）
	ShouldRetry func(error) bool
}

// DefaultRetryConfig 默认重试配置
var DefaultRetryConfig = RetryConfig{
	MaxAttempts:  3,
	InitialDelay: 100 * time.Millisecond,
	MaxDelay:     5 * time.Second,
	Multiplier:   2.0,
	ShouldRetry:  IsRetryable,
}

// RetryableFunc 可重试的函数类型
type RetryableFunc func() error

// RetryWithBackoff 使用指数退避策略重试函数
func RetryWithBackoff(ctx context.Context, config RetryConfig, fn RetryableFunc) error {
	var lastErr error
	delay := config.InitialDelay

	for attempt := 0; attempt < config.MaxAttempts; attempt++ {
		// 执行函数
		err := fn()
		if err == nil {
			return nil
		}

		lastErr = err

		// 判断是否可重试
		shouldRetry := config.ShouldRetry
		if shouldRetry == nil {
			shouldRetry = IsRetryable
		}

		if !shouldRetry(err) {
			return err
		}

		// 如果是最后一次尝试，不再等待
		if attempt == config.MaxAttempts-1 {
			break
		}

		// 等待后重试
		select {
		case <-time.After(delay):
			// 计算下次延迟时间（指数退避）
			delay = time.Duration(float64(delay) * config.Multiplier)
			if delay > config.MaxDelay {
				delay = config.MaxDelay
			}
		case <-ctx.Done():
			return fmt.Errorf("retry cancelled: %w", ctx.Err())
		}
	}

	return fmt.Errorf("max retry attempts (%d) exceeded: %w", config.MaxAttempts, lastErr)
}

// RetryWithJitter 使用带抖动的指数退避策略重试函数
// 抖动可以避免多个客户端同时重试造成的"惊群效应"
func RetryWithJitter(ctx context.Context, config RetryConfig, fn RetryableFunc) error {
	var lastErr error
	delay := config.InitialDelay

	for attempt := 0; attempt < config.MaxAttempts; attempt++ {
		// 执行函数
		err := fn()
		if err == nil {
			return nil
		}

		lastErr = err

		// 判断是否可重试
		shouldRetry := config.ShouldRetry
		if shouldRetry == nil {
			shouldRetry = IsRetryable
		}

		if !shouldRetry(err) {
			return err
		}

		// 如果是最后一次尝试，不再等待
		if attempt == config.MaxAttempts-1 {
			break
		}

		// 添加抖动（0.5 到 1.5 倍的随机延迟）
		jitter := 0.5 + (float64(time.Now().UnixNano()%1000) / 1000.0)
		actualDelay := time.Duration(float64(delay) * jitter)

		// 等待后重试
		select {
		case <-time.After(actualDelay):
			// 计算下次延迟时间（指数退避）
			delay = time.Duration(float64(delay) * config.Multiplier)
			if delay > config.MaxDelay {
				delay = config.MaxDelay
			}
		case <-ctx.Done():
			return fmt.Errorf("retry cancelled: %w", ctx.Err())
		}
	}

	return fmt.Errorf("max retry attempts (%d) exceeded: %w", config.MaxAttempts, lastErr)
}

// CalculateBackoffDelay 计算指数退避延迟时间
func CalculateBackoffDelay(attempt int, config RetryConfig) time.Duration {
	delay := time.Duration(float64(config.InitialDelay) * math.Pow(config.Multiplier, float64(attempt)))
	if delay > config.MaxDelay {
		delay = config.MaxDelay
	}
	return delay
}
