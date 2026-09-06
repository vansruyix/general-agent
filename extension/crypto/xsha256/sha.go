package xsha256

import (
	"crypto/hmac"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
)

func Encrypt(input string) string {
	hash := sha256.New()
	hash.Write([]byte(input))
	encrypted := hash.Sum(nil)
	return hex.EncodeToString(encrypted)
}

func HmacSha1(text string, keyStr string) []byte {
	key := []byte(keyStr)
	mac := hmac.New(sha1.New, key)
	mac.Write([]byte(text))
	return mac.Sum(nil)
}

func HmacSha256(text string, keyStr string) []byte {
	key := []byte(keyStr)
	mac := hmac.New(sha256.New, key)
	mac.Write([]byte(text))
	return mac.Sum(nil)
}
