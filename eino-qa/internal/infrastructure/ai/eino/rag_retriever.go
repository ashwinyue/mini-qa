package eino

import (
	"context"
	"fmt"
	"strings"

	"eino-qa/internal/domain/entity"
	"eino-qa/internal/domain/repository"
	"eino-qa/internal/infrastructure/config"

	"github.com/cloudwego/eino/components/embedding"
	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/schema"
)

// RAGRetriever RAG 检索器
type RAGRetriever struct {
	embedder    embedding.Embedder
	chatModel   model.ChatModel
	vectorRepo  repository.VectorRepository
	topK        int
	scoreThresh float64
}

// NewRAGRetriever 创建新的 RAG 检索器
func NewRAGRetriever(
	client *Client,
	vectorRepo repository.VectorRepository,
	cfg *config.RAGConfig,
) *RAGRetriever {
	topK := 5
	scoreThresh := 0.7

	if cfg != nil {
		if cfg.TopK > 0 {
			topK = cfg.TopK
		}
		if cfg.ScoreThreshold > 0 {
			scoreThresh = cfg.ScoreThreshold
		}
	}

	return &RAGRetriever{
		embedder:    client.GetEmbedModel(),
		chatModel:   client.GetChatModel(),
		vectorRepo:  vectorRepo,
		topK:        topK,
		scoreThresh: scoreThresh,
	}
}

// Retrieve 执行 RAG 检索并生成答案
func (r *RAGRetriever) Retrieve(ctx context.Context, query string) (string, []*entity.Document, error) {
	// 1. 生成查询向量
	vector, err := r.generateQueryVector(ctx, query)
	if err != nil {
		return "", nil, fmt.Errorf("failed to generate query vector: %w", err)
	}

	// 2. 执行向量搜索
	docs, err := r.vectorRepo.Search(ctx, vector, r.topK)
	if err != nil {
		return "", nil, fmt.Errorf("failed to search vectors: %w", err)
	}

	// 3. 过滤低分文档
	filteredDocs := r.filterByScore(docs)

	// 4. 如果没有相关文档，返回未命中
	if len(filteredDocs) == 0 {
		return "", nil, fmt.Errorf("no relevant documents found")
	}

	// 5. 使用检索到的文档生成答案
	answer, err := r.generateAnswer(ctx, query, filteredDocs)
	if err != nil {
		return "", filteredDocs, fmt.Errorf("failed to generate answer: %w", err)
	}

	return answer, filteredDocs, nil
}

// generateQueryVector 生成查询向量
func (r *RAGRetriever) generateQueryVector(ctx context.Context, query string) ([]float32, error) {
	// 使用嵌入模型生成向量
	resp, err := r.embedder.EmbedStrings(ctx, []string{query})
	if err != nil {
		return nil, fmt.Errorf("failed to embed query: %w", err)
	}

	if len(resp) == 0 || len(resp[0]) == 0 {
		return nil, fmt.Errorf("empty embedding result")
	}

	// 转换 float64 到 float32
	vector := make([]float32, len(resp[0]))
	for i, v := range resp[0] {
		vector[i] = float32(v)
	}

	return vector, nil
}

// filterByScore 根据分数过滤文档
func (r *RAGRetriever) filterByScore(docs []*entity.Document) []*entity.Document {
	filtered := make([]*entity.Document, 0, len(docs))
	for _, doc := range docs {
		if doc.Score >= r.scoreThresh {
			filtered = append(filtered, doc)
		}
	}
	return filtered
}

// generateAnswer 使用检索到的文档生成答案
func (r *RAGRetriever) generateAnswer(ctx context.Context, query string, docs []*entity.Document) (string, error) {
	// 构建上下文
	context := r.buildContext(docs)

	// 构建提示词
	systemPrompt := r.buildSystemPrompt()
	userPrompt := r.buildUserPrompt(query, context)

	// 构建消息列表
	messages := []*schema.Message{
		schema.SystemMessage(systemPrompt),
		schema.UserMessage(userPrompt),
	}

	// 调用 LLM 生成答案
	resp, err := r.chatModel.Generate(ctx, messages)
	if err != nil {
		return "", fmt.Errorf("failed to generate answer: %w", err)
	}

	return resp.Content, nil
}

// buildContext 构建检索文档的上下文
func (r *RAGRetriever) buildContext(docs []*entity.Document) string {
	var sb strings.Builder

	for i, doc := range docs {
		sb.WriteString(fmt.Sprintf("文档 %d (相似度: %.2f):\n", i+1, doc.Score))
		sb.WriteString(doc.Content)
		sb.WriteString("\n\n")
	}

	return sb.String()
}

// buildSystemPrompt 构建系统提示词
func (r *RAGRetriever) buildSystemPrompt() string {
	return `你是一个专业的课程咨询助手。你的任务是根据提供的知识库文档，准确回答用户关于课程的问题。

回答要求：
1. 基于提供的文档内容回答，不要编造信息
2. 如果文档中没有相关信息，明确告知用户
3. 回答要清晰、准确、有条理
4. 使用友好、专业的语气
5. 如果需要，可以引用文档中的具体内容

注意：
- 只回答与课程相关的问题
- 不要回答与课程无关的问题
- 如果问题超出知识库范围，建议用户联系人工客服`
}

// buildUserPrompt 构建用户提示词
func (r *RAGRetriever) buildUserPrompt(query, context string) string {
	var sb strings.Builder

	sb.WriteString("知识库文档：\n")
	sb.WriteString(context)
	sb.WriteString("\n")
	sb.WriteString(fmt.Sprintf("用户问题：%s\n\n", query))
	sb.WriteString("请基于上述知识库文档回答用户问题。")

	return sb.String()
}
