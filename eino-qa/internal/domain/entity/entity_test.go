package entity

import (
	"testing"
	"time"
)

// TestMessageCreation 测试消息创建
func TestMessageCreation(t *testing.T) {
	msg := NewMessage("Hello", "user")

	if msg.Content != "Hello" {
		t.Errorf("Expected content 'Hello', got '%s'", msg.Content)
	}

	if msg.Role != "user" {
		t.Errorf("Expected role 'user', got '%s'", msg.Role)
	}

	if err := msg.Validate(); err != nil {
		t.Errorf("Valid message failed validation: %v", err)
	}
}

// TestMessageValidation 测试消息验证
func TestMessageValidation(t *testing.T) {
	tests := []struct {
		name    string
		content string
		role    string
		wantErr bool
	}{
		{"valid user message", "Hello", "user", false},
		{"valid assistant message", "Hi", "assistant", false},
		{"valid system message", "System", "system", false},
		{"empty content", "", "user", true},
		{"invalid role", "Hello", "invalid", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := &Message{Content: tt.content, Role: tt.role}
			err := msg.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestIntentCreation 测试意图创建
func TestIntentCreation(t *testing.T) {
	intent := NewIntent(IntentCourse, 0.95)

	if intent.Type != IntentCourse {
		t.Errorf("Expected type 'course', got '%s'", intent.Type)
	}

	if intent.Confidence != 0.95 {
		t.Errorf("Expected confidence 0.95, got %f", intent.Confidence)
	}

	if err := intent.Validate(); err != nil {
		t.Errorf("Valid intent failed validation: %v", err)
	}
}

// TestIntentValidation 测试意图验证
func TestIntentValidation(t *testing.T) {
	tests := []struct {
		name       string
		intentType IntentType
		confidence float64
		wantErr    bool
	}{
		{"valid course intent", IntentCourse, 0.9, false},
		{"valid order intent", IntentOrder, 0.8, false},
		{"valid direct intent", IntentDirect, 0.7, false},
		{"valid handoff intent", IntentHandoff, 0.6, false},
		{"invalid type", IntentType("invalid"), 0.9, true},
		{"confidence too low", IntentCourse, -0.1, true},
		{"confidence too high", IntentCourse, 1.1, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			intent := &Intent{Type: tt.intentType, Confidence: tt.confidence}
			err := intent.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestDocumentCreation 测试文档创建
func TestDocumentCreation(t *testing.T) {
	doc := NewDocument("Test content", "tenant1")

	if doc.Content != "Test content" {
		t.Errorf("Expected content 'Test content', got '%s'", doc.Content)
	}

	if doc.TenantID != "tenant1" {
		t.Errorf("Expected tenant ID 'tenant1', got '%s'", doc.TenantID)
	}

	if err := doc.Validate(); err != nil {
		t.Errorf("Valid document failed validation: %v", err)
	}
}

// TestOrderCreation 测试订单创建
func TestOrderCreation(t *testing.T) {
	order := NewOrder("user123", "Python Course", 99.99, "tenant1")

	if order.UserID != "user123" {
		t.Errorf("Expected user ID 'user123', got '%s'", order.UserID)
	}

	if order.CourseName != "Python Course" {
		t.Errorf("Expected course name 'Python Course', got '%s'", order.CourseName)
	}

	if order.Amount != 99.99 {
		t.Errorf("Expected amount 99.99, got %f", order.Amount)
	}

	if order.Status != OrderStatusPending {
		t.Errorf("Expected status 'pending', got '%s'", order.Status)
	}

	if err := order.Validate(); err != nil {
		t.Errorf("Valid order failed validation: %v", err)
	}
}

// TestOrderStatusUpdate 测试订单状态更新
func TestOrderStatusUpdate(t *testing.T) {
	order := NewOrder("user123", "Python Course", 99.99, "tenant1")

	err := order.UpdateStatus(OrderStatusPaid)
	if err != nil {
		t.Errorf("Failed to update status: %v", err)
	}

	if order.Status != OrderStatusPaid {
		t.Errorf("Expected status 'paid', got '%s'", order.Status)
	}

	// Test invalid status
	err = order.UpdateStatus(OrderStatus("invalid"))
	if err == nil {
		t.Error("Expected error for invalid status, got nil")
	}
}

// TestSessionCreation 测试会话创建
func TestSessionCreation(t *testing.T) {
	session := NewSession("tenant1", 1*time.Hour)

	if session.TenantID != "tenant1" {
		t.Errorf("Expected tenant ID 'tenant1', got '%s'", session.TenantID)
	}

	if session.GetMessageCount() != 0 {
		t.Errorf("Expected 0 messages, got %d", session.GetMessageCount())
	}

	if err := session.Validate(); err != nil {
		t.Errorf("Valid session failed validation: %v", err)
	}
}

// TestSessionAddMessage 测试会话添加消息
func TestSessionAddMessage(t *testing.T) {
	session := NewSession("tenant1", 1*time.Hour)

	msg := NewMessage("Hello", "user")
	err := session.AddMessage(msg)
	if err != nil {
		t.Errorf("Failed to add message: %v", err)
	}

	if session.GetMessageCount() != 1 {
		t.Errorf("Expected 1 message, got %d", session.GetMessageCount())
	}

	lastMsg := session.GetLastMessage()
	if lastMsg.Content != "Hello" {
		t.Errorf("Expected last message 'Hello', got '%s'", lastMsg.Content)
	}
}

// TestTenantCreation 测试租户创建
func TestTenantCreation(t *testing.T) {
	tenant := NewTenant("tenant1", "Tenant One")

	if tenant.ID != "tenant1" {
		t.Errorf("Expected ID 'tenant1', got '%s'", tenant.ID)
	}

	if tenant.Name != "Tenant One" {
		t.Errorf("Expected name 'Tenant One', got '%s'", tenant.Name)
	}

	if tenant.CollectionName != "kb_tenant1" {
		t.Errorf("Expected collection name 'kb_tenant1', got '%s'", tenant.CollectionName)
	}

	if err := tenant.Validate(); err != nil {
		t.Errorf("Valid tenant failed validation: %v", err)
	}
}

// TestValidateOrderID 测试订单 ID 验证
func TestValidateOrderID(t *testing.T) {
	tests := []struct {
		name    string
		orderID string
		wantErr bool
	}{
		{"valid order ID", "#20251128001", false},
		{"empty order ID", "", true},
		{"invalid format - no hash", "20251128001", true},
		{"invalid format - wrong length", "#2025112800", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateOrderID(tt.orderID)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateOrderID() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestSanitizeSQL 测试 SQL 清理
func TestSanitizeSQL(t *testing.T) {
	tests := []struct {
		name    string
		sql     string
		wantErr bool
	}{
		{"safe SELECT", "SELECT * FROM orders WHERE id = ?", false},
		{"dangerous DROP", "DROP TABLE orders", true},
		{"dangerous DELETE", "DELETE FROM orders", true},
		{"dangerous UPDATE", "UPDATE orders SET status = 'paid'", true},
		{"dangerous UNION", "SELECT * FROM orders UNION SELECT * FROM users", true},
		{"dangerous comment", "SELECT * FROM orders --", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := SanitizeSQL(tt.sql)
			if (err != nil) != tt.wantErr {
				t.Errorf("SanitizeSQL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
