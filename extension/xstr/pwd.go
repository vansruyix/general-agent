package xstr

import (
	"crypto/rand"
	"encoding/base64"
	"strings"
)

// GenerateSecurePassword 生成高强度随机密码
// length: 密码长度（建议 12 位以上）
func GenerateSecurePassword(n int) (string, error) {
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	pwd := base64.URLEncoding.EncodeToString(bytes)
	pwd = strings.TrimRight(pwd, "=")
	return pwd, nil
}

// 辅助函数：判断字符是否在字符串中
func isCharIn(c byte, s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] == c {
			return true
		}
	}
	return false
}

// 使用 crypto/rand 生成 0~n 的随机数
func randInt() int {
	var b [1]byte
	rand.Read(b[:])
	return int(b[0])
}
