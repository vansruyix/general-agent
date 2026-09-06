package xmd5

import (
	"crypto/md5"
	"encoding/hex"
)

func Md5Str(src string) string {
	h := md5.New()
	h.Write([]byte(src))
	return hex.EncodeToString(h.Sum(nil))
}
