package chat

import (
	"context"
	"fmt"
	"time"

	"eino-qa/internal/domain/entity"
)

// ExecuteStream 执行流式对话用例
func (uc *ChatUseCase) ExecuteStream(ctx context.Context, req *ChatRequest) (<-chan *StreamChunk, error) {
	// 验证请求
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("invalid request: %w", err)
	}

	// 创建响应通道
	chunkChan := make(chan *StreamChunk, 10)

	// 在 goroutine 中处理流式响应
	go func() {
		defer close(chunkChan)

		// 记录请求开始
		startTime := time.Now()
		uc.logger.Info(ctx, "stream chat request started", map[string]interface{}{
			"tenant_id":  req.TenantID,
			"session_id": req.SessionID,
			"query":      req.Query,
		})

		// 1. 加载或创建会话
		session, err := uc.loadOrCreateSession(ctx, req.TenantID, req.SessionID)
		if err != nil {
			uc.logger.Error(ctx, "failed to load session", map[string]interface{}{"error": err})
			chunkChan <- &StreamChunk{
				Error: fmt.Errorf("failed to load session: %w", err),
				Done:  true,
			}
			return
		}

		// 2. 添加用户消息到会话
		userMessage := entity.NewMessage(req.Query, "user")
		if err := session.AddMessage(userMessage); err != nil {
			uc.logger.Error(ctx, "failed to add user message", map[string]interface{}{"error": err})
			chunkChan <- &StreamChunk{
				Error: fmt.Errorf("failed to add user message: %w", err),
				Done:  true,
			}
			return
		}

		// 3. 识别意图
		intent, err := uc.intentRecognizer.Recognize(ctx, req.Query, session.GetMessages())
		if err != nil {
			uc.logger.Error(ctx, "failed to recognize intent", map[string]interface{}{"error": err})
			chunkChan <- &StreamChunk{
				Error: fmt.Errorf("failed to recognize intent: %w", err),
				Done:  true,
			}
			return
		}

		uc.logger.Info(ctx, "intent recognized", map[string]interface{}{
			"intent":     intent.Type,
			"confidence": intent.Confidence,
		})

		// 4. 根据意图路由到不同的处理流程
		var fullAnswer string
		var sources []*entity.Document

		switch intent.Type {
		case entity.IntentCourse:
			fullAnswer, sources = uc.handleCourseIntentStream(ctx, req.Query, chunkChan)
		case entity.IntentOrder:
			fullAnswer = uc.handleOrderIntentStream(ctx, req.Query, chunkChan)
		case entity.IntentDirect:
			fullAnswer = uc.handleDirectIntentStream(ctx, req.Query, session.GetMessages(), chunkChan)
		case entity.IntentHandoff:
			fullAnswer = uc.handleHandoffIntent(ctx, intent)
			chunkChan <- &StreamChunk{Content: fullAnswer}
		default:
			fullAnswer = uc.responseGenerator.GenerateFallbackMessage()
			chunkChan <- &StreamChunk{Content: fullAnswer}
		}

		// 5. 添加助手消息到会话
		assistantMessage := entity.NewMessage(fullAnswer, "assistant")
		if err := session.AddMessage(assistantMessage); err != nil {
			uc.logger.Error(ctx, "failed to add assistant message", map[string]interface{}{"error": err})
		}

		// 6. 保存会话
		if err := uc.sessionRepo.Save(ctx, session); err != nil {
			uc.logger.Error(ctx, "failed to save session", map[string]interface{}{"error": err})
		}

		// 记录请求完成
		duration := time.Since(startTime)
		uc.logger.Info(ctx, "stream chat request completed", map[string]interface{}{
			"duration_ms": duration.Milliseconds(),
			"intent":      intent.Type,
		})

		// 7. 发送完成标记
		chunkChan <- &StreamChunk{
			Done: true,
			Metadata: map[string]any{
				"intent":      intent.Type,
				"confidence":  intent.Confidence,
				"duration_ms": duration.Milliseconds(),
				"session_id":  session.ID,
				"sources":     sources,
			},
		}
	}()

	return chunkChan, nil
}

// handleCourseIntentStream 处理课程咨询意图（流式）
func (uc *ChatUseCase) handleCourseIntentStream(ctx context.Context, query string, chunkChan chan<- *StreamChunk) (string, []*entity.Document) {
	uc.logger.Info(ctx, "handling course intent (stream)", map[string]interface{}{"query": query})

	// RAG 检索不支持流式，直接返回完整结果
	answer, sources, err := uc.ragRetriever.Retrieve(ctx, query)
	if err != nil {
		uc.logger.Error(ctx, "RAG retrieval failed", map[string]interface{}{"error": err})
		answer = uc.responseGenerator.GenerateFallbackMessage()
		sources = nil
	}

	// 发送完整答案
	chunkChan <- &StreamChunk{Content: answer}

	return answer, sources
}

// handleOrderIntentStream 处理订单查询意图（流式）
func (uc *ChatUseCase) handleOrderIntentStream(ctx context.Context, query string, chunkChan chan<- *StreamChunk) string {
	uc.logger.Info(ctx, "handling order intent (stream)", map[string]interface{}{"query": query})

	// 订单查询不支持流式，直接返回完整结果
	answer, err := uc.orderQuerier.Query(ctx, query)
	if err != nil {
		uc.logger.Error(ctx, "order query failed", map[string]interface{}{"error": err})
		answer = uc.responseGenerator.GenerateErrorMessage(err)
	}

	// 发送完整答案
	chunkChan <- &StreamChunk{Content: answer}

	return answer
}

// handleDirectIntentStream 处理直接回答意图（流式）
func (uc *ChatUseCase) handleDirectIntentStream(ctx context.Context, query string, history []*entity.Message, chunkChan chan<- *StreamChunk) string {
	uc.logger.Info(ctx, "handling direct intent (stream)", map[string]interface{}{"query": query})

	// 使用响应生成器的流式接口
	contentChan, errorChan := uc.responseGenerator.GenerateStream(ctx, query, history)

	var fullAnswer string

	// 读取流式响应
	for {
		select {
		case content, ok := <-contentChan:
			if !ok {
				// 通道关闭，流式响应结束
				return fullAnswer
			}
			fullAnswer += content
			chunkChan <- &StreamChunk{Content: content}

		case err, ok := <-errorChan:
			if ok && err != nil {
				uc.logger.Error(ctx, "stream generation failed", map[string]interface{}{"error": err})
				errorMsg := uc.responseGenerator.GenerateErrorMessage(err)
				chunkChan <- &StreamChunk{Content: errorMsg}
				return fullAnswer + errorMsg
			}

		case <-ctx.Done():
			uc.logger.Warn(ctx, "stream context cancelled", map[string]interface{}{})
			return fullAnswer
		}
	}
}
