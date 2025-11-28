package entity

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"sync/atomic"
	"time"
)

var counter uint64

// generateUniqueID 生成唯一 ID
func generateUniqueID(prefix string, length int) string {
	// 使用原子计数器确保唯一性
	count := atomic.AddUint64(&counter, 1)

	// 生成随机字节
	randomBytes := make([]byte, length/2)
	if _, err := rand.Read(randomBytes); err != nil {
		// 如果随机数生成失败，使用时间戳和计数器
		return fmt.Sprintf("%s%d%d", prefix, time.Now().UnixNano(), count)
	}

	// 转换为十六进制字符串
	randomStr := hex.EncodeToString(randomBytes)

	return prefix + randomStr[:length]
}

// randomString 生成随机字符串（使用加密安全的随机数）
func randomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, n)

	// 使用加密安全的随机数
	randomBytes := make([]byte, n)
	if _, err := rand.Read(randomBytes); err != nil {
		// 降级到时间戳方案
		for i := range b {
			b[i] = letters[(time.Now().UnixNano()+int64(i))%int64(len(letters))]
		}
		return string(b)
	}

	for i := range b {
		b[i] = letters[int(randomBytes[i])%len(letters)]
	}

	return string(b)
}
