package chat

import "context"

// ChatUseCaseInterface 对话用例接口
type ChatUseCaseInterface interface {
	Execute(ctx context.Context, req *ChatRequest) (*ChatResponse, error)
	ExecuteStream(ctx context.Context, req *ChatRequest) (<-chan *StreamChunk, error)
}
