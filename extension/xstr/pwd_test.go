package xstr

import (
	"fmt"
	"testing"
)

func TestGeneratePassword(t *testing.T) {
	password, err := GenerateSecurePassword(12)
	if err != nil {
		t.Error(err)
	}

	fmt.Println("生成的高强度密码：", password)
}
