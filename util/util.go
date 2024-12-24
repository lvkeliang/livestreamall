package util

import (
	"crypto/sha256"
	"encoding/hex"
)

// HashString 使用SHA256哈希函数将字符串映射到唯一字符串
func HashString(input string) string {
	hasher := sha256.New()
	hasher.Write([]byte(input))
	return hex.EncodeToString(hasher.Sum(nil))
}
