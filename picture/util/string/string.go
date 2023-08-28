package string

import (
	"crypto/rand"
	"encoding/hex"
	//"math/rand"
)

// RandStr 随机生成字符串
func RandStr(n int) string {
	result := make([]byte, n)
	if _, err := rand.Read(result); err != nil {
		return ""
	}
	return hex.EncodeToString(result)
}
