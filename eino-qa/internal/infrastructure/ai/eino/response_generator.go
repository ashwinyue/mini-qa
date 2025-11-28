package eino

import (
	"context"
	"fmt"
	"strings"

	"eino-qa/internal/domain/entity"

	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/schema"
)

// ResponseGenerator 响应生成器
type ResponseGenerator struct {
	chatModel model.ChatModel
}

// NewResponseGenerator 创建新的响应生成器
func NewResponseGenerator(client *Client) *ResponseGenerator {
	return &ResponseGenerator{
		chatModel: client.GetChatModel(),
	}
}

// Generate 生成直接回答
func (g *ResponseGenerator) Generate(ctx context.Context, query string, history []*entity.Message) (string, error) {
	// 构建提示词
	systemPrompt := g.buildSystemPrompt()
	userPrompt := g.buildUserPrompt(query, history)

	// 构建消息列表
	messages := []*schema.Message{
		schema.SystemMessage(systemPrompt),
	}

	// 添加历史消息
	for _, msg := range history {
		if msg.IsUser() {
			messages = append(messages, schema.UserMessage(msg.Content))
		} else if msg.IsAssistant() {
			messages = append(messages, schema.AssistantMessage(msg.Content, nil))
		}
	}

	// 添加当前查询
	messages = append(messages, schema.UserMessage(userPrompt))

	// 调用 LLM 生成回答
	resp, err := g.chatModel.Generate(ctx, messages)
	if err != nil {
		return "", fmt.Errorf("failed to generate response: %w", err)
	}

	return resp.Content, nil
}

// GenerateStream 生成流式回答
func (g *ResponseGenerator) GenerateStream(ctx context.Context, query string, history []*entity.Message) (<-chan string, <-chan error) {
	resultChan := make(chan string, 10)
	errorChan := make(chan error, 1)

	go func() {
		defer close(resultChan)
		defer close(errorChan)

		// 构建提示词
		systemPrompt := g.buildSystemPrompt()
		userPrompt := g.buildUserPrompt(query, history)

		// 构建消息列表
		messages := []*schema.Message{
			schema.SystemMessage(systemPrompt),
		}

		// 添加历史消息
		for _, msg := range history {
			if msg.IsUser() {
				messages = append(messages, schema.UserMessage(msg.Content))
			} else if msg.IsAssistant() {
				messages = append(messages, schema.AssistantMessage(msg.Content, nil))
			}
		}

		// 添加当前查询
		messages = append(messages, schema.UserMessage(userPrompt))

		// 调用 LLM 生成流式回答
		streamReader, err := g.chatModel.Stream(ctx, messages)
		if err != nil {
			errorChan <- fmt.Errorf("failed to start stream: %w", err)
			return
		}

		// 读取流式响应
		for {
			chunk, err := streamReader.Recv()
			if err != nil {
				if err.Error() != "EOF" {
					errorChan <- fmt.Errorf("stream error: %w", err)
				}
				break
			}

			if chunk.Content != "" {
				resultChan <- chunk.Content
			}
		}
	}()

	return resultChan, errorChan
}

// buildSystemPrompt 构建系统提示词
func (g *ResponseGenerator) buildSystemPrompt() string {
	return `你是一个友好、专业的智能客服助手。你的任务是回答用户的问题，提供帮助和支持。

回答要求：
1. 语气友好、热情、专业
2. 回答简洁明了，重点突出
3. 对于简单的问候和闲聊，给予适当的回应
4. 对于不确定的问题，诚实告知并建议联系人工客服
5. 保持礼貌和耐心

注意事项：
- 不要编造信息
- 不要回答与业务无关的问题
- 如果问题超出能力范围，建议用户联系人工客服
- 保护用户隐私，不要询问敏感信息`
}

// buildUserPrompt 构建用户提示词
func (g *ResponseGenerator) buildUserPrompt(query string, history []*entity.Message) string {
	// 对于直接回答，通常不需要额外的上下文构建
	// 历史消息已经在消息列表中了
	return query
}

// GenerateHandoffMessage 生成人工转接消息
func (g *ResponseGenerator) GenerateHandoffMessage(reason string) string {
	messages := []string{
		"正在为您转接人工客服，请稍候...",
		"您的问题比较复杂，让我为您转接专业客服人员。",
		"为了更好地帮助您，我将为您转接人工客服。",
	}

	// 根据原因选择合适的消息
	if reason != "" {
		return fmt.Sprintf("由于%s，正在为您转接人工客服，请稍候...", reason)
	}

	// 默认返回第一条消息
	return messages[0]
}

// GenerateErrorMessage 生成错误消息
func (g *ResponseGenerator) GenerateErrorMessage(err error) string {
	// 根据错误类型生成友好的错误消息
	errMsg := err.Error()

	if strings.Contains(errMsg, "timeout") {
		return "抱歉，系统响应超时，请稍后重试。"
	}

	if strings.Contains(errMsg, "not found") {
		return "抱歉，未找到相关信息。"
	}

	if strings.Contains(errMsg, "network") {
		return "抱歉，网络连接出现问题，请稍后重试。"
	}

	// 默认错误消息
	return "抱歉，系统遇到了一些问题，请稍后重试或联系人工客服。"
}

// GenerateFallbackMessage 生成降级消息
func (g *ResponseGenerator) GenerateFallbackMessage() string {
	return "抱歉，我暂时无法回答您的问题。请联系人工客服获取帮助，或稍后重试。"
}
