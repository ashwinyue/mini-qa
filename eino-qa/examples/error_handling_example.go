package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"eino-qa/pkg/errors"
)

// 模拟外部服务调用
var callCount int

func simulateExternalService() error {
	callCount++
	if callCount < 3 {
		return errors.NewWithCategory(
			502,
			"外部服务暂时不可用",
			fmt.Errorf("connection timeout"),
			errors.CategoryExternal,
			true,
		)
	}
	return nil
}

func simulateUnstableService() error {
	callCount++
	if callCount%2 == 0 {
		return fmt.Errorf("service error")
	}
	return nil
}

// Example1_BasicRetry 基本重试示例
func Example1_BasicRetry() {
	fmt.Println("=== Example 1: 基本重试 ===")
	callCount = 0

	ctx := context.Background()
	err := errors.RetryWithBackoff(
		ctx,
		errors.DefaultRetryConfig,
		func() error {
			fmt.Printf("尝试调用外部服务 (第 %d 次)\n", callCount+1)
			return simulateExternalService()
		},
	)

	if err != nil {
		log.Printf("调用失败: %v\n", err)
	} else {
		fmt.Println("调用成功!")
	}
	fmt.Println()
}

// Example2_CustomRetryConfig 自定义重试配置示例
func Example2_CustomRetryConfig() {
	fmt.Println("=== Example 2: 自定义重试配置 ===")
	callCount = 0

	config := errors.RetryConfig{
		MaxAttempts:  5,
		InitialDelay: 50 * time.Millisecond,
		MaxDelay:     2 * time.Second,
		Multiplier:   2.0,
		ShouldRetry:  errors.IsRetryable,
	}

	ctx := context.Background()
	err := errors.RetryWithBackoff(
		ctx,
		config,
		func() error {
			fmt.Printf("尝试调用服务 (第 %d 次)\n", callCount+1)
			return simulateExternalService()
		},
	)

	if err != nil {
		log.Printf("调用失败: %v\n", err)
	} else {
		fmt.Println("调用成功!")
	}
	fmt.Println()
}

// Example3_Fallback 降级策略示例
func Example3_Fallback() {
	fmt.Println("=== Example 3: 降级策略 ===")

	err := errors.HandleWithFallback(
		func() error {
			fmt.Println("尝试主要服务...")
			return fmt.Errorf("主要服务不可用")
		},
		func() error {
			fmt.Println("执行降级逻辑...")
			fmt.Println("返回默认响应")
			return nil
		},
	)

	if err != nil {
		log.Printf("处理失败: %v\n", err)
	} else {
		fmt.Println("处理成功（使用降级）")
	}
	fmt.Println()
}

// Example4_FallbackChain 降级链示例
func Example4_FallbackChain() {
	fmt.Println("=== Example 4: 降级链 ===")

	chain := errors.NewFallbackChain()

	chain.Add("primary", func() error {
		fmt.Println("尝试主要服务...")
		return fmt.Errorf("主要服务失败")
	}, func() error {
		fmt.Println("尝试备份服务...")
		return fmt.Errorf("备份服务失败")
	})

	chain.Add("secondary", func() error {
		fmt.Println("尝试次要服务...")
		return fmt.Errorf("次要服务失败")
	}, func() error {
		fmt.Println("返回缓存结果...")
		return nil
	})

	ctx := context.Background()
	err := chain.Execute(ctx)

	if err != nil {
		log.Printf("所有降级处理失败: %v\n", err)
	} else {
		fmt.Println("处理成功（使用降级）")
	}
	fmt.Println()
}

// Example5_CircuitBreaker 熔断器示例
func Example5_CircuitBreaker() {
	fmt.Println("=== Example 5: 熔断器 ===")
	callCount = 0

	cb := errors.NewCircuitBreaker(3, 5*time.Second)

	// 模拟多次调用
	for i := 0; i < 10; i++ {
		err := cb.Execute(func() error {
			fmt.Printf("调用服务 (第 %d 次)\n", i+1)
			return simulateUnstableService()
		})

		if err != nil {
			fmt.Printf("  失败: %v\n", err)
		} else {
			fmt.Println("  成功")
		}

		// 检查熔断器状态
		state := cb.GetState()
		switch state {
		case errors.StateClosed:
			fmt.Println("  熔断器状态: 关闭（正常）")
		case errors.StateOpen:
			fmt.Println("  熔断器状态: 打开（熔断）")
		case errors.StateHalfOpen:
			fmt.Println("  熔断器状态: 半开（尝试恢复）")
		}

		time.Sleep(100 * time.Millisecond)
	}
	fmt.Println()
}

// Example6_ErrorHandler 综合错误处理器示例
func Example6_ErrorHandler() {
	fmt.Println("=== Example 6: 综合错误处理器 ===")
	callCount = 0

	handler := errors.NewErrorHandler(errors.RetryConfig{
		MaxAttempts:  3,
		InitialDelay: 100 * time.Millisecond,
		MaxDelay:     1 * time.Second,
		Multiplier:   2.0,
	}).WithCircuitBreaker(5, 10*time.Second).
		WithFallback(func() error {
			fmt.Println("执行降级逻辑...")
			return nil
		})

	ctx := context.Background()
	err := handler.Execute(ctx, func() error {
		fmt.Printf("尝试调用服务 (第 %d 次)\n", callCount+1)
		return simulateExternalService()
	})

	if err != nil {
		log.Printf("处理失败: %v\n", err)
	} else {
		fmt.Println("处理成功!")
	}
	fmt.Println()
}

// Example7_LLMCallWithRetry LLM 调用重试示例
func Example7_LLMCallWithRetry() {
	fmt.Println("=== Example 7: LLM 调用重试 ===")

	// 模拟 LLM 调用
	callLLM := func() error {
		fmt.Println("调用 LLM...")
		// 模拟偶尔失败
		if time.Now().Unix()%3 == 0 {
			return errors.NewWithCategory(
				502,
				"LLM 服务暂时不可用",
				fmt.Errorf("rate limit exceeded"),
				errors.CategoryExternal,
				true,
			)
		}
		fmt.Println("LLM 调用成功")
		return nil
	}

	ctx := context.Background()
	err := errors.RetryWithBackoff(ctx, errors.DefaultRetryConfig, callLLM)

	if err != nil {
		// 降级：返回预设回复
		fmt.Println("降级：返回预设回复")
		fmt.Println("回复: 抱歉，我现在无法回答您的问题，请稍后再试。")
	}
	fmt.Println()
}

// Example8_MilvusSearchWithFallback Milvus 搜索降级示例
func Example8_MilvusSearchWithFallback() {
	fmt.Println("=== Example 8: Milvus 搜索降级 ===")

	var results []string

	err := errors.HandleWithFallback(
		func() error {
			fmt.Println("尝试向量搜索...")
			// 模拟 Milvus 不可用
			return errors.ErrMilvusUnavailable
		},
		func() error {
			fmt.Println("降级到关键词搜索...")
			results = []string{"关键词匹配结果1", "关键词匹配结果2"}
			return nil
		},
	)

	if err != nil {
		log.Printf("搜索失败: %v\n", err)
	} else {
		fmt.Printf("搜索结果: %v\n", results)
	}
	fmt.Println()
}

func main() {
	fmt.Println("错误处理和重试机制示例")
	fmt.Println("================================\n")

	Example1_BasicRetry()
	Example2_CustomRetryConfig()
	Example3_Fallback()
	Example4_FallbackChain()
	Example5_CircuitBreaker()
	Example6_ErrorHandler()
	Example7_LLMCallWithRetry()
	Example8_MilvusSearchWithFallback()

	fmt.Println("所有示例执行完成!")
}
