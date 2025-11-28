package entity

import "errors"

var (
	// Tenant 相关错误
	ErrInvalidTenantID = errors.New("invalid tenant ID")
)

// Tenant 表示租户值对象
type Tenant struct {
	ID             string
	Name           string
	CollectionName string
	DatabasePath   string
	Metadata       map[string]any
}

// NewTenant 创建新的租户实例
func NewTenant(id, name string) *Tenant {
	return &Tenant{
		ID:             id,
		Name:           name,
		CollectionName: generateCollectionName(id),
		DatabasePath:   generateDatabasePath(id),
		Metadata:       make(map[string]any),
	}
}

// Validate 验证租户的有效性
func (t *Tenant) Validate() error {
	if t.ID == "" {
		return ErrInvalidTenantID
	}

	if t.Name == "" {
		return errors.New("tenant name cannot be empty")
	}

	return nil
}

// IsDefault 判断是否为默认租户
func (t *Tenant) IsDefault() bool {
	return t.ID == "default"
}

// generateCollectionName 生成 Milvus Collection 名称
func generateCollectionName(tenantID string) string {
	return "kb_" + tenantID
}

// generateDatabasePath 生成 SQLite 数据库路径
func generateDatabasePath(tenantID string) string {
	return "data/" + tenantID + ".db"
}
