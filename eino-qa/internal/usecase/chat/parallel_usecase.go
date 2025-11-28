package chat

import (
	"context"
	"fmt"
	"sync"
	"time"

	"eino-qa/internal/domain/entity"
)

// ParallelQueryResult 并行查询结果
type ParallelQueryResult struct {
	CourseAnswer  string             // 课程咨询答案
	CourseSources []*entity.Document // 课程来源文档
	CourseError   error              // 课程查询错误

	OrderAnswer string // 订单查询答案
	OrderError  error  // 订单查询错误

	Duration time.Duration // 总耗时
}

// ExecuteParallel 并行执行多个查询
// 这个方法适用于需要同时查询多个数据源的场景
func (uc *ChatUseCase) ExecuteParallel(ctx context.Context, req *ChatRequest) (*ParallelQueryResult, error) {
	// 验证请求
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("invalid request: %w", err)
	}

	startTime := time.Now()
	uc.logger.Info(ctx, "parallel query started", map[string]interface{}{
		"tenant_id": req.TenantID,
		"query":     req.Query,
	})

	// 创建结果结构
	result := &ParallelQueryResult{}

	// 使用 WaitGroup 等待所有查询完成
	var wg sync.WaitGroup

	// 并行执行课程查询
	wg.Add(1)
	go func() {
		defer wg.Done()
		answer, sources, err := uc.ragRetriever.Retrieve(ctx, req.Query)
		result.CourseAnswer = answer
		result.CourseSources = sources
		result.CourseError = err

		if err != nil {
			uc.logger.Error(ctx, "parallel course query failed", map[string]interface{}{"error": err})
		}
	}()

	// 并行执行订单查询
	wg.Add(1)
	go func() {
		defer wg.Done()
		answer, err := uc.orderQuerier.Query(ctx, req.Query)
		result.OrderAnswer = answer
		result.OrderError = err

		if err != nil {
			uc.logger.Error(ctx, "parallel order query failed", map[string]interface{}{"error": err})
		}
	}()

	// 等待所有查询完成
	wg.Wait()

	result.Duration = time.Since(startTime)
	uc.logger.Info(ctx, "parallel query completed", map[string]interface{}{
		"duration_ms": result.Duration.Milliseconds(),
	})

	return result, nil
}

// ExecuteParallelWithTimeout 带超时的并行查询
func (uc *ChatUseCase) ExecuteParallelWithTimeout(ctx context.Context, req *ChatRequest, timeout time.Duration) (*ParallelQueryResult, error) {
	// 创建带超时的上下文
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	// 创建结果通道
	resultChan := make(chan *ParallelQueryResult, 1)
	errorChan := make(chan error, 1)

	// 在 goroutine 中执行并行查询
	go func() {
		result, err := uc.ExecuteParallel(ctx, req)
		if err != nil {
			errorChan <- err
			return
		}
		resultChan <- result
	}()

	// 等待结果或超时
	select {
	case result := <-resultChan:
		return result, nil
	case err := <-errorChan:
		return nil, err
	case <-ctx.Done():
		return nil, fmt.Errorf("parallel query timeout: %w", ctx.Err())
	}
}

// MergeParallelResults 合并并行查询结果
// 这个方法将多个数据源的结果合并成一个统一的答案
func (uc *ChatUseCase) MergeParallelResults(ctx context.Context, result *ParallelQueryResult, query string) (string, []*entity.Document, error) {
	// 如果两个查询都失败，返回错误
	if result.CourseError != nil && result.OrderError != nil {
		return "", nil, fmt.Errorf("all parallel queries failed")
	}

	// 如果只有课程查询成功
	if result.CourseError == nil && result.OrderError != nil {
		return result.CourseAnswer, result.CourseSources, nil
	}

	// 如果只有订单查询成功
	if result.OrderError == nil && result.CourseError != nil {
		return result.OrderAnswer, nil, nil
	}

	// 如果两个查询都成功，使用 LLM 合并结果
	mergedAnswer, err := uc.mergeAnswersWithLLM(ctx, query, result)
	if err != nil {
		uc.logger.Error(ctx, "failed to merge answers", map[string]interface{}{"error": err})
		// 降级：返回课程答案
		return result.CourseAnswer, result.CourseSources, nil
	}

	return mergedAnswer, result.CourseSources, nil
}

// mergeAnswersWithLLM 使用 LLM 合并多个答案
func (uc *ChatUseCase) mergeAnswersWithLLM(ctx context.Context, query string, result *ParallelQueryResult) (string, error) {
	// 构建合并提示词
	userPrompt := fmt.Sprintf(`用户问题：%s

课程信息：
%s

订单信息：
%s

请将上述信息整合成一个统一的回答。`, query, result.CourseAnswer, result.OrderAnswer)

	// 使用响应生成器生成合并后的答案
	messages := []*entity.Message{
		{Content: userPrompt, Role: "user"},
	}

	answer, err := uc.responseGenerator.Generate(ctx, userPrompt, messages)
	if err != nil {
		return "", fmt.Errorf("failed to merge with LLM: %w", err)
	}

	return answer, nil
}

// ExecuteWithParallelRetrieval 执行带并行检索的对话
// 这个方法会在识别意图后，根据需要并行查询多个数据源
func (uc *ChatUseCase) ExecuteWithParallelRetrieval(ctx context.Context, req *ChatRequest) (*ChatResponse, error) {
	// 验证请求
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("invalid request: %w", err)
	}

	startTime := time.Now()
	uc.logger.Info(ctx, "chat with parallel retrieval started", map[string]interface{}{
		"tenant_id":  req.TenantID,
		"session_id": req.SessionID,
		"query":      req.Query,
	})

	// 1. 加载或创建会话
	session, err := uc.loadOrCreateSession(ctx, req.TenantID, req.SessionID)
	if err != nil {
		return nil, fmt.Errorf("failed to load session: %w", err)
	}

	// 2. 添加用户消息
	userMessage := entity.NewMessage(req.Query, "user")
	if err := session.AddMessage(userMessage); err != nil {
		return nil, fmt.Errorf("failed to add user message: %w", err)
	}

	// 3. 识别意图
	intent, err := uc.intentRecognizer.Recognize(ctx, req.Query, session.GetMessages())
	if err != nil {
		return nil, fmt.Errorf("failed to recognize intent: %w", err)
	}

	uc.logger.Info(ctx, "intent recognized", map[string]interface{}{
		"intent":     intent.Type,
		"confidence": intent.Confidence,
	})

	var answer string
	var sources []*entity.Document

	// 4. 根据意图决定是否使用并行检索
	switch intent.Type {
	case entity.IntentCourse:
		// 单一数据源，不需要并行
		answer, sources, err = uc.ragRetriever.Retrieve(ctx, req.Query)
		if err != nil {
			answer = uc.responseGenerator.GenerateFallbackMessage()
		}

	case entity.IntentOrder:
		// 单一数据源，不需要并行
		answer, err = uc.orderQuerier.Query(ctx, req.Query)
		if err != nil {
			answer = uc.responseGenerator.GenerateErrorMessage(err)
		}

	case entity.IntentDirect:
		// 可能需要多个数据源，使用并行检索
		parallelResult, err := uc.ExecuteParallelWithTimeout(ctx, req, 5*time.Second)
		if err != nil {
			uc.logger.Error(ctx, "parallel retrieval failed", map[string]interface{}{"error": err})
			// 降级到普通响应生成
			answer, _ = uc.responseGenerator.Generate(ctx, req.Query, session.GetMessages())
		} else {
			// 合并并行查询结果
			answer, sources, err = uc.MergeParallelResults(ctx, parallelResult, req.Query)
			if err != nil {
				answer = uc.responseGenerator.GenerateFallbackMessage()
			}
		}

	case entity.IntentHandoff:
		answer = uc.handleHandoffIntent(ctx, intent)

	default:
		answer = uc.responseGenerator.GenerateFallbackMessage()
	}

	// 5. 添加助手消息
	assistantMessage := entity.NewMessage(answer, "assistant")
	if err := session.AddMessage(assistantMessage); err != nil {
		uc.logger.Error(ctx, "failed to add assistant message", map[string]interface{}{"error": err})
	}

	// 6. 保存会话
	if err := uc.sessionRepo.Save(ctx, session); err != nil {
		uc.logger.Error(ctx, "failed to save session", map[string]interface{}{"error": err})
	}

	duration := time.Since(startTime)
	uc.logger.Info(ctx, "chat with parallel retrieval completed", map[string]interface{}{
		"duration_ms": duration.Milliseconds(),
		"intent":      intent.Type,
	})

	// 7. 构建响应
	response := &ChatResponse{
		Answer:    answer,
		Route:     string(intent.Type),
		Sources:   sources,
		SessionID: session.ID,
		Metadata: map[string]any{
			"intent":      intent.Type,
			"confidence":  intent.Confidence,
			"duration_ms": duration.Milliseconds(),
		},
	}

	return response, nil
}
