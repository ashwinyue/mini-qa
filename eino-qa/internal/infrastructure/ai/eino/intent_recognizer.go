package eino

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"eino-qa/internal/domain/entity"
	"eino-qa/internal/infrastructure/config"

	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/schema"
)

// IntentRecognizer 意图识别器
type IntentRecognizer struct {
	chatModel           model.ChatModel
	confidenceThreshold float64
}

// NewIntentRecognizer 创建新的意图识别器
func NewIntentRecognizer(client *Client, cfg *config.IntentConfig) *IntentRecognizer {
	threshold := 0.7
	if cfg != nil && cfg.ConfidenceThreshold > 0 {
		threshold = cfg.ConfidenceThreshold
	}

	return &IntentRecognizer{
		chatModel:           client.GetChatModel(),
		confidenceThreshold: threshold,
	}
}

// Recognize 识别用户查询的意图
func (r *IntentRecognizer) Recognize(ctx context.Context, query string, history []*entity.Message) (*entity.Intent, error) {
	// 构建提示词
	systemPrompt := r.buildSystemPrompt()
	userPrompt := r.buildUserPrompt(query, history)

	// 构建消息列表
	messages := []*schema.Message{
		schema.SystemMessage(systemPrompt),
		schema.UserMessage(userPrompt),
	}

	// 调用 LLM
	resp, err := r.chatModel.Generate(ctx, messages)
	if err != nil {
		return nil, fmt.Errorf("failed to generate intent: %w", err)
	}

	// 解析意图
	intent, err := r.parseIntent(resp.Content)
	if err != nil {
		return nil, fmt.Errorf("failed to parse intent: %w", err)
	}

	// 验证意图
	if err := intent.Validate(); err != nil {
		return nil, fmt.Errorf("invalid intent: %w", err)
	}

	// 如果置信度低于阈值，转人工
	if !intent.IsHighConfidence(r.confidenceThreshold) {
		intent.Type = entity.IntentHandoff
	}

	return intent, nil
}

// buildSystemPrompt 构建系统提示词
func (r *IntentRecognizer) buildSystemPrompt() string {
	return `你是一个智能客服意图识别助手。你的任务是分析用户的查询，判断用户的意图类型。

意图类型定义：
1. course - 课程咨询：用户询问课程内容、课程安排、学习资料等与课程相关的问题
2. order - 订单查询：用户查询订单状态、订单详情、退款等与订单相关的问题
3. direct - 直接回答：简单的问候、闲聊或可以直接回答的一般性问题
4. handoff - 人工转接：复杂问题、投诉、或需要人工处理的情况

请以 JSON 格式返回结果，包含以下字段：
{
  "intent": "意图类型（course/order/direct/handoff）",
  "confidence": 置信度分数（0-1之间的浮点数）,
  "reason": "判断理由"
}

注意：
- 只返回 JSON，不要包含其他文字
- confidence 必须是 0 到 1 之间的数字
- 如果不确定，将 confidence 设置为较低的值`
}

// buildUserPrompt 构建用户提示词
func (r *IntentRecognizer) buildUserPrompt(query string, history []*entity.Message) string {
	var sb strings.Builder

	// 添加历史对话上下文（最近3轮）
	if len(history) > 0 {
		sb.WriteString("对话历史：\n")
		start := len(history) - 6 // 最近3轮（6条消息）
		if start < 0 {
			start = 0
		}
		for i := start; i < len(history); i++ {
			msg := history[i]
			sb.WriteString(fmt.Sprintf("%s: %s\n", msg.Role, msg.Content))
		}
		sb.WriteString("\n")
	}

	// 添加当前查询
	sb.WriteString(fmt.Sprintf("当前用户查询：%s\n\n", query))
	sb.WriteString("请分析用户意图并返回 JSON 结果。")

	return sb.String()
}

// parseIntent 解析意图结果
func (r *IntentRecognizer) parseIntent(content string) (*entity.Intent, error) {
	// 清理可能的 markdown 代码块标记
	content = strings.TrimSpace(content)
	content = strings.TrimPrefix(content, "```json")
	content = strings.TrimPrefix(content, "```")
	content = strings.TrimSuffix(content, "```")
	content = strings.TrimSpace(content)

	// 解析 JSON
	var result struct {
		Intent     string  `json:"intent"`
		Confidence float64 `json:"confidence"`
		Reason     string  `json:"reason"`
	}

	if err := json.Unmarshal([]byte(content), &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal intent json: %w", err)
	}

	// 转换为实体
	intentType := r.mapIntentType(result.Intent)
	intent := entity.NewIntent(intentType, result.Confidence)
	intent.Metadata["reason"] = result.Reason

	return intent, nil
}

// mapIntentType 映射意图类型字符串到实体类型
func (r *IntentRecognizer) mapIntentType(intentStr string) entity.IntentType {
	switch strings.ToLower(strings.TrimSpace(intentStr)) {
	case "course":
		return entity.IntentCourse
	case "order":
		return entity.IntentOrder
	case "direct":
		return entity.IntentDirect
	case "handoff":
		return entity.IntentHandoff
	default:
		// 默认转人工
		return entity.IntentHandoff
	}
}
