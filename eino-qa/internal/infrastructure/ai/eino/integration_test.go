package eino

import (
	"testing"
	"time"

	"eino-qa/internal/domain/entity"
)

// TestClientInitialization 测试客户端初始化
func TestClientInitialization(t *testing.T) {
	// 跳过集成测试（需要真实的 API Key）
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	config := ClientConfig{
		APIKey:     "test-api-key",
		ChatModel:  "qwen-turbo",
		EmbedModel: "text-embedding-v2",
		MaxRetries: 3,
		Timeout:    30 * time.Second,
	}

	client, err := NewClient(config)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	if client.GetChatModel() == nil {
		t.Error("Chat model is nil")
	}

	if client.GetEmbedModel() == nil {
		t.Error("Embed model is nil")
	}
}

// TestIntentRecognizerBasic 测试意图识别器基本功能
func TestIntentRecognizerBasic(t *testing.T) {
	// 这是一个单元测试，测试意图识别器的结构
	// 实际的 LLM 调用需要在集成测试中进行

	// 测试意图类型映射
	recognizer := &IntentRecognizer{
		confidenceThreshold: 0.7,
	}

	tests := []struct {
		input    string
		expected entity.IntentType
	}{
		{"course", entity.IntentCourse},
		{"order", entity.IntentOrder},
		{"direct", entity.IntentDirect},
		{"handoff", entity.IntentHandoff},
		{"COURSE", entity.IntentCourse},
		{"  order  ", entity.IntentOrder},
		{"unknown", entity.IntentHandoff}, // 未知类型默认转人工
	}

	for _, tt := range tests {
		result := recognizer.mapIntentType(tt.input)
		if result != tt.expected {
			t.Errorf("mapIntentType(%q) = %v, want %v", tt.input, result, tt.expected)
		}
	}
}

// TestOrderQuerierRegexExtraction 测试订单 ID 正则提取
func TestOrderQuerierRegexExtraction(t *testing.T) {
	querier := &OrderQuerier{}

	tests := []struct {
		query    string
		expected string
	}{
		{"查询订单#20251114001", "20251114001"},
		{"订单号：20251114001", "20251114001"},
		{"订单:20251114001", "20251114001"},
		{"我的订单 20251114001 怎么样了", "20251114001"},
		{"没有订单号", ""},
	}

	for _, tt := range tests {
		result := querier.extractOrderIDByRegex(tt.query)
		if result != tt.expected {
			t.Errorf("extractOrderIDByRegex(%q) = %q, want %q", tt.query, result, tt.expected)
		}
	}
}

// TestOrderQuerierSQLValidation 测试 SQL 验证
func TestOrderQuerierSQLValidation(t *testing.T) {
	querier := &OrderQuerier{}

	tests := []struct {
		sql       string
		shouldErr bool
	}{
		{"SELECT * FROM orders WHERE id = ?", false},
		{"SELECT id, amount FROM orders", false},
		{"DROP TABLE orders", true},
		{"DELETE FROM orders", true},
		{"UPDATE orders SET status = 'paid'", true},
		{"INSERT INTO orders VALUES (?)", true},
		{"SELECT * FROM orders; DROP TABLE users", true},
		{"SELECT * FROM orders -- comment", true},
	}

	for _, tt := range tests {
		err := querier.ValidateSQL(tt.sql)
		if (err != nil) != tt.shouldErr {
			t.Errorf("ValidateSQL(%q) error = %v, shouldErr = %v", tt.sql, err, tt.shouldErr)
		}
	}
}

// TestResponseGeneratorMessages 测试响应生成器消息生成
func TestResponseGeneratorMessages(t *testing.T) {
	generator := &ResponseGenerator{}

	// 测试人工转接消息
	msg := generator.GenerateHandoffMessage("置信度过低")
	if msg == "" {
		t.Error("GenerateHandoffMessage returned empty string")
	}

	// 测试降级消息
	fallback := generator.GenerateFallbackMessage()
	if fallback == "" {
		t.Error("GenerateFallbackMessage returned empty string")
	}

	// 测试错误消息
	testErr := entity.ErrOrderNotFound
	errMsg := generator.GenerateErrorMessage(testErr)
	if errMsg == "" {
		t.Error("GenerateErrorMessage returned empty string")
	}
}

// TestRAGRetrieverScoreFiltering 测试 RAG 检索器分数过滤
func TestRAGRetrieverScoreFiltering(t *testing.T) {
	retriever := &RAGRetriever{
		scoreThresh: 0.7,
	}

	docs := []*entity.Document{
		{ID: "1", Content: "Doc 1", Score: 0.9},
		{ID: "2", Content: "Doc 2", Score: 0.8},
		{ID: "3", Content: "Doc 3", Score: 0.6},
		{ID: "4", Content: "Doc 4", Score: 0.5},
	}

	filtered := retriever.filterByScore(docs)

	if len(filtered) != 2 {
		t.Errorf("Expected 2 documents after filtering, got %d", len(filtered))
	}

	for _, doc := range filtered {
		if doc.Score < 0.7 {
			t.Errorf("Document %s has score %.2f, below threshold", doc.ID, doc.Score)
		}
	}
}

// TestIntentRecognizerPromptBuilding 测试意图识别器提示词构建
func TestIntentRecognizerPromptBuilding(t *testing.T) {
	recognizer := &IntentRecognizer{}

	// 测试系统提示词
	systemPrompt := recognizer.buildSystemPrompt()
	if systemPrompt == "" {
		t.Error("System prompt is empty")
	}

	// 测试用户提示词（无历史）
	userPrompt := recognizer.buildUserPrompt("测试查询", nil)
	if userPrompt == "" {
		t.Error("User prompt is empty")
	}

	// 测试用户提示词（有历史）
	history := []*entity.Message{
		entity.NewMessage("你好", "user"),
		entity.NewMessage("您好！有什么可以帮您？", "assistant"),
	}
	userPromptWithHistory := recognizer.buildUserPrompt("测试查询", history)
	if userPromptWithHistory == "" {
		t.Error("User prompt with history is empty")
	}
}

// TestRAGRetrieverContextBuilding 测试 RAG 检索器上下文构建
func TestRAGRetrieverContextBuilding(t *testing.T) {
	retriever := &RAGRetriever{}

	docs := []*entity.Document{
		{ID: "1", Content: "Python 是一门编程语言", Score: 0.9},
		{ID: "2", Content: "Go 是一门编程语言", Score: 0.8},
	}

	context := retriever.buildContext(docs)
	if context == "" {
		t.Error("Context is empty")
	}

	// 验证上下文包含文档内容
	if !contains(context, "Python") || !contains(context, "Go") {
		t.Error("Context does not contain document content")
	}
}

// contains 检查字符串是否包含子串
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// BenchmarkIntentRecognizerMapType 基准测试意图类型映射
func BenchmarkIntentRecognizerMapType(b *testing.B) {
	recognizer := &IntentRecognizer{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		recognizer.mapIntentType("course")
	}
}

// BenchmarkOrderQuerierRegexExtraction 基准测试订单 ID 提取
func BenchmarkOrderQuerierRegexExtraction(b *testing.B) {
	querier := &OrderQuerier{}
	query := "查询订单#20251114001"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		querier.extractOrderIDByRegex(query)
	}
}
