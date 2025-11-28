package utils

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"
)

// GenerateID 生成唯一 ID
func GenerateID() string {
	timestamp := time.Now().UnixNano()
	randomBytes := make([]byte, 8)
	rand.Read(randomBytes)
	return fmt.Sprintf("%d-%s", timestamp, hex.EncodeToString(randomBytes))
}

// GenerateSessionID 生成会话 ID
func GenerateSessionID() string {
	return "session-" + GenerateID()
}

// GenerateRequestID 生成请求 ID
func GenerateRequestID() string {
	return "req-" + GenerateID()
}
