package chat

import (
	"context"
	"fmt"
	"time"

	"eino-qa/internal/domain/entity"
	"eino-qa/internal/domain/repository"
	"eino-qa/internal/infrastructure/ai/eino"
	"eino-qa/internal/infrastructure/logger"
)

// ChatUseCase 对话用例
type ChatUseCase struct {
	intentRecognizer  *eino.IntentRecognizer
	ragRetriever      *eino.RAGRetriever
	orderQuerier      *eino.OrderQuerier
	responseGenerator *eino.ResponseGenerator
	sessionRepo       repository.SessionRepository
	sessionTTL        time.Duration
	logger            logger.Logger
}

// NewChatUseCase 创建新的对话用例
func NewChatUseCase(
	intentRecognizer *eino.IntentRecognizer,
	ragRetriever *eino.RAGRetriever,
	orderQuerier *eino.OrderQuerier,
	responseGenerator *eino.ResponseGenerator,
	sessionRepo repository.SessionRepository,
	sessionTTL time.Duration,
	log logger.Logger,
) *ChatUseCase {
	if sessionTTL == 0 {
		sessionTTL = 30 * time.Minute // 默认 30 分钟
	}

	return &ChatUseCase{
		intentRecognizer:  intentRecognizer,
		ragRetriever:      ragRetriever,
		orderQuerier:      orderQuerier,
		responseGenerator: responseGenerator,
		sessionRepo:       sessionRepo,
		sessionTTL:        sessionTTL,
		logger:            log,
	}
}

// Execute 执行对话用例
func (uc *ChatUseCase) Execute(ctx context.Context, req *ChatRequest) (*ChatResponse, error) {
	// 验证请求
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("invalid request: %w", err)
	}

	// 记录请求开始
	startTime := time.Now()
	uc.logger.Info(ctx, "chat request started", map[string]interface{}{
		"tenant_id":  req.TenantID,
		"session_id": req.SessionID,
		"query":      req.Query,
	})

	// 1. 加载或创建会话
	session, err := uc.loadOrCreateSession(ctx, req.TenantID, req.SessionID)
	if err != nil {
		uc.logger.Error(ctx, "failed to load session", map[string]interface{}{"error": err})
		return nil, fmt.Errorf("failed to load session: %w", err)
	}

	// 2. 添加用户消息到会话
	userMessage := entity.NewMessage(req.Query, "user")
	if err := session.AddMessage(userMessage); err != nil {
		uc.logger.Error(ctx, "failed to add user message", map[string]interface{}{"error": err})
		return nil, fmt.Errorf("failed to add user message: %w", err)
	}

	// 3. 识别意图
	intent, err := uc.intentRecognizer.Recognize(ctx, req.Query, session.GetMessages())
	if err != nil {
		uc.logger.Error(ctx, "failed to recognize intent", map[string]interface{}{"error": err})
		return nil, fmt.Errorf("failed to recognize intent: %w", err)
	}

	uc.logger.Info(ctx, "intent recognized", map[string]interface{}{
		"intent":     intent.Type,
		"confidence": intent.Confidence,
	})

	// 4. 根据意图路由到不同的处理流程
	var answer string
	var sources []*entity.Document
	var routeErr error

	switch intent.Type {
	case entity.IntentCourse:
		answer, sources, routeErr = uc.handleCourseIntent(ctx, req.Query)
	case entity.IntentOrder:
		answer, routeErr = uc.handleOrderIntent(ctx, req.Query)
	case entity.IntentDirect:
		answer, routeErr = uc.handleDirectIntent(ctx, req.Query, session.GetMessages())
	case entity.IntentHandoff:
		answer = uc.handleHandoffIntent(ctx, intent)
	default:
		answer = uc.responseGenerator.GenerateFallbackMessage()
	}

	// 处理路由错误
	if routeErr != nil {
		uc.logger.Error(ctx, "route handling failed", map[string]interface{}{
			"intent": intent.Type,
			"error":  routeErr,
		})
		answer = uc.responseGenerator.GenerateErrorMessage(routeErr)
	}

	// 5. 添加助手消息到会话
	assistantMessage := entity.NewMessage(answer, "assistant")
	if err := session.AddMessage(assistantMessage); err != nil {
		uc.logger.Error(ctx, "failed to add assistant message", map[string]interface{}{"error": err})
		// 不返回错误，因为回答已经生成
	}

	// 6. 保存会话
	if err := uc.sessionRepo.Save(ctx, session); err != nil {
		uc.logger.Error(ctx, "failed to save session", map[string]interface{}{"error": err})
		// 不返回错误，因为回答已经生成
	}

	// 记录请求完成
	duration := time.Since(startTime)
	uc.logger.Info(ctx, "chat request completed", map[string]interface{}{
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

// loadOrCreateSession 加载或创建会话
func (uc *ChatUseCase) loadOrCreateSession(ctx context.Context, tenantID, sessionID string) (*entity.Session, error) {
	// 如果提供了会话 ID，尝试加载
	if sessionID != "" {
		session, err := uc.sessionRepo.Load(ctx, sessionID)
		if err == nil {
			// 检查会话是否过期
			if !session.IsExpired() {
				// 延长会话过期时间
				session.ExtendExpiration(uc.sessionTTL)
				return session, nil
			}
			// 会话已过期，创建新会话
			uc.logger.Info(ctx, "session expired, creating new session", map[string]interface{}{"old_session_id": sessionID})
		}
	}

	// 创建新会话
	session := entity.NewSession(tenantID, uc.sessionTTL)
	uc.logger.Info(ctx, "new session created", map[string]interface{}{"session_id": session.ID})

	return session, nil
}

// handleCourseIntent 处理课程咨询意图
func (uc *ChatUseCase) handleCourseIntent(ctx context.Context, query string) (string, []*entity.Document, error) {
	uc.logger.Info(ctx, "handling course intent", map[string]interface{}{"query": query})

	// 使用 RAG 检索器
	answer, sources, err := uc.ragRetriever.Retrieve(ctx, query)
	if err != nil {
		uc.logger.Error(ctx, "RAG retrieval failed", map[string]interface{}{"error": err})
		// 如果 RAG 失败，返回降级消息
		return uc.responseGenerator.GenerateFallbackMessage(), nil, nil
	}

	return answer, sources, nil
}

// handleOrderIntent 处理订单查询意图
func (uc *ChatUseCase) handleOrderIntent(ctx context.Context, query string) (string, error) {
	uc.logger.Info(ctx, "handling order intent", map[string]interface{}{"query": query})

	// 使用订单查询器
	answer, err := uc.orderQuerier.Query(ctx, query)
	if err != nil {
		uc.logger.Error(ctx, "order query failed", map[string]interface{}{"error": err})
		return "", err
	}

	return answer, nil
}

// handleDirectIntent 处理直接回答意图
func (uc *ChatUseCase) handleDirectIntent(ctx context.Context, query string, history []*entity.Message) (string, error) {
	uc.logger.Info(ctx, "handling direct intent", map[string]interface{}{"query": query})

	// 使用响应生成器
	answer, err := uc.responseGenerator.Generate(ctx, query, history)
	if err != nil {
		uc.logger.Error(ctx, "response generation failed", map[string]interface{}{"error": err})
		return "", err
	}

	return answer, nil
}

// handleHandoffIntent 处理人工转接意图
func (uc *ChatUseCase) handleHandoffIntent(ctx context.Context, intent *entity.Intent) string {
	uc.logger.Info(ctx, "handling handoff intent", map[string]interface{}{})

	reason := ""
	if r, ok := intent.Metadata["reason"].(string); ok {
		reason = r
	}

	return uc.responseGenerator.GenerateHandoffMessage(reason)
}
