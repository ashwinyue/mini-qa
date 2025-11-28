package eino

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"eino-qa/internal/domain/entity"
	"eino-qa/internal/domain/repository"

	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/schema"
)

// OrderQuerier 订单查询器
type OrderQuerier struct {
	chatModel model.ChatModel
	orderRepo repository.OrderRepository
}

// NewOrderQuerier 创建新的订单查询器
func NewOrderQuerier(client *Client, orderRepo repository.OrderRepository) *OrderQuerier {
	return &OrderQuerier{
		chatModel: client.GetChatModel(),
		orderRepo: orderRepo,
	}
}

// Query 查询订单信息
func (q *OrderQuerier) Query(ctx context.Context, query string) (string, error) {
	// 1. 提取订单 ID
	orderID, err := q.extractOrderID(ctx, query)
	if err != nil {
		return "", fmt.Errorf("failed to extract order id: %w", err)
	}

	if orderID == "" {
		return "抱歉，我没有在您的问题中找到订单号。请提供订单号，格式如：#20251114001", nil
	}

	// 2. 查询订单
	order, err := q.orderRepo.FindByID(ctx, orderID)
	if err != nil {
		if err == entity.ErrOrderNotFound {
			return fmt.Sprintf("抱歉，未找到订单号为 %s 的订单。请确认订单号是否正确。", orderID), nil
		}
		return "", fmt.Errorf("failed to query order: %w", err)
	}

	// 3. 格式化订单信息为自然语言
	answer, err := q.formatOrderInfo(ctx, query, order)
	if err != nil {
		return "", fmt.Errorf("failed to format order info: %w", err)
	}

	return answer, nil
}

// extractOrderID 从用户查询中提取订单 ID
func (q *OrderQuerier) extractOrderID(ctx context.Context, query string) (string, error) {
	// 首先尝试使用正则表达式提取
	orderID := q.extractOrderIDByRegex(query)
	if orderID != "" {
		return orderID, nil
	}

	// 如果正则提取失败，使用 LLM 提取
	return q.extractOrderIDByLLM(ctx, query)
}

// extractOrderIDByRegex 使用正则表达式提取订单 ID
func (q *OrderQuerier) extractOrderIDByRegex(query string) string {
	// 匹配订单号格式：#20251114001 或 20251114001
	patterns := []string{
		`#(\d{11,})`,          // #开头的订单号
		`订单号[：:]\s*(\d{11,})`, // 订单号：xxxxx
		`订单[：:]\s*(\d{11,})`,  // 订单：xxxxx
		`\b(\d{11,})\b`,       // 独立的11位以上数字
	}

	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindStringSubmatch(query)
		if len(matches) > 1 {
			return matches[1]
		}
	}

	return ""
}

// extractOrderIDByLLM 使用 LLM 提取订单 ID
func (q *OrderQuerier) extractOrderIDByLLM(ctx context.Context, query string) (string, error) {
	systemPrompt := `你是一个订单号提取助手。从用户的查询中提取订单号。

订单号特征：
- 通常是11位或更长的数字
- 可能以 # 开头
- 可能在"订单号"、"订单"等词后面

请以 JSON 格式返回结果：
{
  "order_id": "提取到的订单号（只包含数字，不包含#）",
  "found": true/false
}

如果没有找到订单号，返回：
{
  "order_id": "",
  "found": false
}`

	userPrompt := fmt.Sprintf("用户查询：%s\n\n请提取订单号。", query)

	messages := []*schema.Message{
		schema.SystemMessage(systemPrompt),
		schema.UserMessage(userPrompt),
	}

	resp, err := q.chatModel.Generate(ctx, messages)
	if err != nil {
		return "", fmt.Errorf("failed to generate extraction: %w", err)
	}

	// 解析结果
	content := strings.TrimSpace(resp.Content)
	content = strings.TrimPrefix(content, "```json")
	content = strings.TrimPrefix(content, "```")
	content = strings.TrimSuffix(content, "```")
	content = strings.TrimSpace(content)

	var result struct {
		OrderID string `json:"order_id"`
		Found   bool   `json:"found"`
	}

	if err := json.Unmarshal([]byte(content), &result); err != nil {
		return "", fmt.Errorf("failed to parse extraction result: %w", err)
	}

	if !result.Found {
		return "", nil
	}

	return result.OrderID, nil
}

// formatOrderInfo 将订单信息格式化为自然语言
func (q *OrderQuerier) formatOrderInfo(ctx context.Context, query string, order *entity.Order) (string, error) {
	// 构建订单信息的结构化描述
	orderInfo := fmt.Sprintf(`订单信息：
- 订单号：%s
- 课程名称：%s
- 订单金额：%.2f 元
- 订单状态：%s
- 创建时间：%s`,
		order.ID,
		order.CourseName,
		order.Amount,
		q.formatOrderStatus(order.Status),
		order.CreatedAt.Format("2006-01-02 15:04:05"),
	)

	// 使用 LLM 生成自然语言回复
	systemPrompt := `你是一个专业的客服助手。根据订单信息，用自然、友好的语言回答用户的问题。

要求：
1. 语气友好、专业
2. 信息准确、完整
3. 根据用户的具体问题重点回答
4. 如果订单状态异常，提供相应的建议`

	userPrompt := fmt.Sprintf("%s\n\n用户问题：%s\n\n请根据订单信息回答用户问题。", orderInfo, query)

	messages := []*schema.Message{
		schema.SystemMessage(systemPrompt),
		schema.UserMessage(userPrompt),
	}

	resp, err := q.chatModel.Generate(ctx, messages)
	if err != nil {
		return "", fmt.Errorf("failed to generate answer: %w", err)
	}

	return resp.Content, nil
}

// formatOrderStatus 格式化订单状态
func (q *OrderQuerier) formatOrderStatus(status entity.OrderStatus) string {
	statusMap := map[entity.OrderStatus]string{
		entity.OrderStatusPending:   "待支付",
		entity.OrderStatusPaid:      "已支付",
		entity.OrderStatusRefunded:  "已退款",
		entity.OrderStatusCancelled: "已取消",
	}

	if formatted, ok := statusMap[status]; ok {
		return formatted
	}

	return string(status)
}

// ValidateSQL 验证 SQL 查询的安全性
func (q *OrderQuerier) ValidateSQL(sql string) error {
	// 转换为大写进行检查
	upperSQL := strings.ToUpper(sql)

	// 检查危险操作
	dangerousKeywords := []string{
		"DROP",
		"DELETE",
		"UPDATE",
		"INSERT",
		"ALTER",
		"CREATE",
		"TRUNCATE",
		"EXEC",
		"EXECUTE",
		"--",
		";",
	}

	for _, keyword := range dangerousKeywords {
		if strings.Contains(upperSQL, keyword) {
			return fmt.Errorf("dangerous SQL keyword detected: %s", keyword)
		}
	}

	// 确保只包含 SELECT 语句
	if !strings.HasPrefix(strings.TrimSpace(upperSQL), "SELECT") {
		return fmt.Errorf("only SELECT statements are allowed")
	}

	return nil
}
