package pwd

import (
	"encoding/base64"
	"fmt"
	"general-agent/extension/crypto/xaes"
	"general-agent/extension/crypto/xmd5"
	"general-agent/extension/logz"
	"unicode"
)

const (
	AES_KEY  = "venus.liangjie.agentcc66" // key的长度必须是16、24或者32字节
	MD5_SLAT = "liangjie.venus"
)

// Md5 密码加密
func Md5(pass string) string {
	return xmd5.Md5Str(pass + MD5_SLAT)
}

// LoginEncrypt 登录加密
func LoginEncrypt(pass string) (string, error) {
	key := []byte(AES_KEY)
	aesPass, err := xaes.AesEncrypt([]byte(pass), key)
	if err != nil {
		logz.WarnNoCtx(fmt.Sprintf("aes encrypt error: %v", err))
		return "", err
	}
	return base64.StdEncoding.EncodeToString(aesPass), nil
}

// LoginDecrypt 登录解密
func LoginDecrypt(pass string) (string, error) {
	bytesPass, err := base64.StdEncoding.DecodeString(pass)
	if err != nil {
		logz.WarnNoCtx(fmt.Sprintf("base64 decode error: %v", err))
		return "", err
	}

	tpass, err := xaes.AesDecrypt(bytesPass, []byte(AES_KEY))
	if err != nil {
		logz.WarnNoCtx(fmt.Sprintf("aes decrypt error: %v", err))
		return "", err
	}
	return string(tpass), nil
}

// IsValidPassword 检查密码是否合法（至少8位，且必须包含字母和数字）
func IsValidPassword(password string) bool {
	// 检查密码长度
	if len(password) < 8 {
		return false
	}

	// 检查是否包含字母和数字
	hasLetter := false
	hasDigit := false

	for _, char := range password {
		if unicode.IsLetter(char) {
			hasLetter = true
		}
		if unicode.IsDigit(char) {
			hasDigit = true
		}
	}

	return hasLetter && hasDigit
}
