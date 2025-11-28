package handler

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"eino-qa/internal/adapter/http/middleware"
	"eino-qa/internal/usecase/chat"

	"github.com/gin-gonic/gin"
)

// ChatHandler 对话处理器
type ChatHandler struct {
	chatUseCase chat.ChatUseCaseInterface
}

// NewChatHandler 创建对话处理器
func NewChatHandler(chatUseCase chat.ChatUseCaseInterface) *ChatHandler {
	return &ChatHandler{
		chatUseCase: chatUseCase,
	}
}

// ChatRequestDTO HTTP 请求 DTO
type ChatRequestDTO struct {
	Query     string `json:"query" binding:"required"`
	TenantID  string `json:"tenant_id"`
	SessionID string `json:"session_id"`
	Stream    bool   `json:"stream"`
}

// ChatResponseDTO HTTP 响应 DTO
type ChatResponseDTO struct {
	Answer    string         `json:"answer"`
	Route     string         `json:"route"`
	Sources   []SourceDTO    `json:"sources,omitempty"`
	SessionID string         `json:"session_id"`
	Metadata  map[string]any `json:"metadata"`
}

// SourceDTO 来源文档 DTO
type SourceDTO struct {
	Content  string         `json:"content"`
	Score    float64        `json:"score"`
	Metadata map[string]any `json:"metadata,omitempty"`
}

// HandleChat 处理对话请求
// POST /chat
// 需求: 6.1, 6.2, 6.3, 6.4, 6.5
func (h *ChatHandler) HandleChat(c *gin.Context) {
	var req ChatRequestDTO

	// 解析请求体
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(middleware.NewBadRequestError(fmt.Sprintf("invalid request: %s", err.Error())))
		return
	}

	// 从 context 获取租户 ID（由中间件设置）
	if tenantID, exists := c.Get("tenant_id"); exists {
		if tid, ok := tenantID.(string); ok && req.TenantID == "" {
			req.TenantID = tid
		}
	}

	// 如果仍然没有租户 ID，使用默认值
	if req.TenantID == "" {
		req.TenantID = "default"
	}

	// 构建用例请求
	useCaseReq := &chat.ChatRequest{
		Query:     req.Query,
		TenantID:  req.TenantID,
		SessionID: req.SessionID,
		Stream:    req.Stream,
	}

	// 如果请求流式响应
	if req.Stream {
		h.handleStreamChat(c, useCaseReq)
		return
	}

	// 普通响应
	h.handleNormalChat(c, useCaseReq)
}

// handleNormalChat 处理普通对话请求
func (h *ChatHandler) handleNormalChat(c *gin.Context, req *chat.ChatRequest) {
	// 执行对话用例
	resp, err := h.chatUseCase.Execute(c.Request.Context(), req)
	if err != nil {
		c.Error(err)
		return
	}

	// 转换为 DTO
	dto := h.toChatResponseDTO(resp)

	// 返回响应
	c.JSON(http.StatusOK, dto)
}

// handleStreamChat 处理流式对话请求
// 需求: 6.4 - 流式响应支持（SSE）
func (h *ChatHandler) handleStreamChat(c *gin.Context, req *chat.ChatRequest) {
	// 设置 SSE 响应头
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("Transfer-Encoding", "chunked")

	// 执行流式对话用例
	chunkChan, err := h.chatUseCase.ExecuteStream(c.Request.Context(), req)
	if err != nil {
		c.Error(err)
		return
	}

	// 获取 ResponseWriter 的 flusher
	flusher, ok := c.Writer.(http.Flusher)
	if !ok {
		c.Error(middleware.NewServiceError("streaming not supported", "http"))
		return
	}

	// 流式发送响应
	c.Stream(func(w io.Writer) bool {
		select {
		case chunk, ok := <-chunkChan:
			if !ok {
				// 通道关闭，流结束
				return false
			}

			// 如果有错误，发送错误事件
			if chunk.Error != nil {
				c.SSEvent("error", map[string]any{
					"message": chunk.Error.Error(),
				})
				flusher.Flush()
				return false
			}

			// 如果完成，发送完成事件
			if chunk.Done {
				c.SSEvent("done", map[string]any{
					"metadata": chunk.Metadata,
				})
				flusher.Flush()
				return false
			}

			// 发送内容块
			c.SSEvent("message", map[string]any{
				"content": chunk.Content,
			})
			flusher.Flush()
			return true

		case <-c.Request.Context().Done():
			// 客户端断开连接
			return false
		}
	})
}

// toChatResponseDTO 转换为响应 DTO
func (h *ChatHandler) toChatResponseDTO(resp *chat.ChatResponse) *ChatResponseDTO {
	dto := &ChatResponseDTO{
		Answer:    resp.Answer,
		Route:     resp.Route,
		SessionID: resp.SessionID,
		Metadata:  resp.Metadata,
	}

	// 转换来源文档
	if len(resp.Sources) > 0 {
		dto.Sources = make([]SourceDTO, len(resp.Sources))
		for i, source := range resp.Sources {
			dto.Sources[i] = SourceDTO{
				Content:  source.Content,
				Score:    source.Score,
				Metadata: source.Metadata,
			}
		}
	}

	return dto
}

// HandleChatWithTimeout 带超时的对话处理
// 可选的辅助方法，用于设置请求超时
func (h *ChatHandler) HandleChatWithTimeout(timeout time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 创建带超时的 context
		ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)
		defer cancel()

		// 替换 request context
		c.Request = c.Request.WithContext(ctx)

		// 调用主处理器
		h.HandleChat(c)
	}
}
