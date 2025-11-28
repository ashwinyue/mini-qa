package chat

import "eino-qa/internal/domain/entity"

// ChatRequest 对话请求
type ChatRequest struct {
	Query     string // 用户查询
	TenantID  string // 租户 ID
	SessionID string // 会话 ID
	Stream    bool   // 是否流式响应
}

// ChatResponse 对话响应
type ChatResponse struct {
	Answer    string             // 回答内容
	Route     string             // 路由类型（意图类型）
	Sources   []*entity.Document // 来源文档（RAG 检索结果）
	SessionID string             // 会话 ID
	Metadata  map[string]any     // 元数据
}

// StreamChunk 流式响应块
type StreamChunk struct {
	Content  string         // 内容片段
	Done     bool           // 是否完成
	Error    error          // 错误信息
	Metadata map[string]any // 元数据
}

// Validate 验证请求
func (r *ChatRequest) Validate() error {
	if r.Query == "" {
		return entity.ErrEmptyContent
	}

	if r.TenantID == "" {
		return entity.ErrEmptyTenantID
	}

	return nil
}
