package entity

import "time"

// Document 表示知识库中的文档
type Document struct {
	ID        string
	Content   string
	Vector    []float32
	Metadata  map[string]any
	Score     float64 // 相似度分数
	TenantID  string
	CreatedAt time.Time
}

// NewDocument 创建新的文档实例
func NewDocument(content string, tenantID string) *Document {
	return &Document{
		ID:        generateDocumentID(),
		Content:   content,
		Metadata:  make(map[string]any),
		TenantID:  tenantID,
		CreatedAt: time.Now(),
	}
}

// Validate 验证文档的有效性
func (d *Document) Validate() error {
	if d.Content == "" {
		return ErrEmptyContent
	}

	if d.TenantID == "" {
		return ErrEmptyTenantID
	}

	return nil
}

// HasVector 判断文档是否有向量
func (d *Document) HasVector() bool {
	return len(d.Vector) > 0
}

// SetVector 设置文档向量
func (d *Document) SetVector(vector []float32) {
	d.Vector = vector
}

// SetScore 设置相似度分数
func (d *Document) SetScore(score float64) {
	d.Score = score
}

// AddMetadata 添加元数据
func (d *Document) AddMetadata(key string, value any) {
	if d.Metadata == nil {
		d.Metadata = make(map[string]any)
	}
	d.Metadata[key] = value
}

// GetMetadata 获取元数据
func (d *Document) GetMetadata(key string) (any, bool) {
	if d.Metadata == nil {
		return nil, false
	}
	value, exists := d.Metadata[key]
	return value, exists
}

// generateDocumentID 生成文档 ID
func generateDocumentID() string {
	return "doc_" + time.Now().Format("20060102150405") + randomString(12)
}
