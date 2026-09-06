package pwd_test

import (
	"fmt"
	"general-agent/extension/crypto/pwd"
	"log"
	"testing"
)

func TestPass(t *testing.T) {
	ps := []string{"Dimrealm123@#%"}
	fmt.Printf("原密码\t|\t加密后\t|\t解密后\t|\tMD5\n")
	for _, p := range ps {
		encrypt(p)
	}
}
func encrypt(password string) {
	// 加密
	xpass, err := pwd.LoginEncrypt(password)
	if err != nil {
		log.Fatal(err)
		return
	}

	// 解密
	realPass, err := pwd.LoginDecrypt(xpass)
	if err != nil {
		log.Fatal(err)
		return
	}

	fmt.Printf("%s\t|\t%s\t|\t%s\t|\t%s\n", password, xpass, realPass, pwd.Md5(string(realPass)))
}

func TestIsValidPassword(t *testing.T) {
	fmt.Println(pwd.IsValidPassword("FeZv4cOQuFBHsGj4"))
}
