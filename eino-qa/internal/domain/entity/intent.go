package entity

// IntentType 定义意图类型
type IntentType string

const (
	// IntentCourse 课程咨询意图
	IntentCourse IntentType = "course"
	// IntentOrder 订单查询意图
	IntentOrder IntentType = "order"
	// IntentDirect 直接回答意图
	IntentDirect IntentType = "direct"
	// IntentHandoff 人工转接意图
	IntentHandoff IntentType = "handoff"
)

// Intent 表示用户查询的意图
type Intent struct {
	Type       IntentType
	Confidence float64
	Metadata   map[string]any
}

// NewIntent 创建新的意图实例
func NewIntent(intentType IntentType, confidence float64) *Intent {
	return &Intent{
		Type:       intentType,
		Confidence: confidence,
		Metadata:   make(map[string]any),
	}
}

// Validate 验证意图的有效性
func (i *Intent) Validate() error {
	validTypes := map[IntentType]bool{
		IntentCourse:  true,
		IntentOrder:   true,
		IntentDirect:  true,
		IntentHandoff: true,
	}

	if !validTypes[i.Type] {
		return ErrInvalidIntentType
	}

	if i.Confidence < 0 || i.Confidence > 1 {
		return ErrInvalidConfidence
	}

	return nil
}

// IsHighConfidence 判断是否为高置信度
func (i *Intent) IsHighConfidence(threshold float64) bool {
	return i.Confidence >= threshold
}

// IsCourse 判断是否为课程咨询意图
func (i *Intent) IsCourse() bool {
	return i.Type == IntentCourse
}

// IsOrder 判断是否为订单查询意图
func (i *Intent) IsOrder() bool {
	return i.Type == IntentOrder
}

// IsDirect 判断是否为直接回答意图
func (i *Intent) IsDirect() bool {
	return i.Type == IntentDirect
}

// IsHandoff 判断是否为人工转接意图
func (i *Intent) IsHandoff() bool {
	return i.Type == IntentHandoff
}
