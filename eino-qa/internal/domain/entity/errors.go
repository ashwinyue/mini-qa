package entity

import "errors"

var (
	// Message 相关错误
	ErrEmptyContent = errors.New("message content cannot be empty")
	ErrInvalidRole  = errors.New("invalid message role")

	// Intent 相关错误
	ErrInvalidIntentType = errors.New("invalid intent type")
	ErrInvalidConfidence = errors.New("confidence must be between 0 and 1")

	// Document 相关错误
	ErrEmptyTenantID = errors.New("tenant ID cannot be empty")

	// Order 相关错误
	ErrEmptyUserID        = errors.New("user ID cannot be empty")
	ErrEmptyCourseName    = errors.New("course name cannot be empty")
	ErrInvalidAmount      = errors.New("amount must be non-negative")
	ErrInvalidOrderStatus = errors.New("invalid order status")
	ErrOrderNotFound      = errors.New("order not found")

	// Session 相关错误
	ErrEmptySessionID = errors.New("session ID cannot be empty")
	ErrSessionExpired = errors.New("session has expired")
)
