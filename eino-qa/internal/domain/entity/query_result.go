package entity

// QueryResult 表示查询结果
type QueryResult struct {
	Answer   string
	Route    string // 路由类型: "course", "order", "direct", "handoff"
	Sources  []*Document
	Intent   *Intent
	Metadata map[string]any
}

// NewQueryResult 创建新的查询结果实例
func NewQueryResult(answer, route string) *QueryResult {
	return &QueryResult{
		Answer:   answer,
		Route:    route,
		Sources:  make([]*Document, 0),
		Metadata: make(map[string]any),
	}
}

// AddSource 添加来源文档
func (qr *QueryResult) AddSource(doc *Document) {
	qr.Sources = append(qr.Sources, doc)
}

// SetIntent 设置意图
func (qr *QueryResult) SetIntent(intent *Intent) {
	qr.Intent = intent
}

// AddMetadata 添加元数据
func (qr *QueryResult) AddMetadata(key string, value any) {
	if qr.Metadata == nil {
		qr.Metadata = make(map[string]any)
	}
	qr.Metadata[key] = value
}

// HasSources 判断是否有来源文档
func (qr *QueryResult) HasSources() bool {
	return len(qr.Sources) > 0
}

// GetSourceCount 获取来源文档数量
func (qr *QueryResult) GetSourceCount() int {
	return len(qr.Sources)
}
