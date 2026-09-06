package rule_crypto

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestReadIniFile(t *testing.T) {
	decoderMap := make(map[string]Decoder)
	decoderMap["firstid"] = &Int16Decoder{}
	decoderMap["wagsafe"] = &Int16Decoder{}
	decoderMap["id"] = &Int16Decoder{}

	data, err := ReadIniFile("example/80/eventsdisc.ini", decoderMap)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("总共", len(data), "条")
	for _, entry := range data {
		fmt.Println(entry)
	}
}

func TestReadAll(t *testing.T) {
	decoderMap := make(map[string]Decoder)
	decoderMap["securityid"] = &Int16Decoder{}

	dir := "example/db/"
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println("Error:", err)
			return err
		}
		// 检查文件是否是 INI 文件
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".ini") {
			data, err := ReadIniFile(dir+info.Name(), decoderMap)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println(info.Name(), "总共", len(data), "条")
			for _, entry := range data {
				_, err := fmt.Println(entry)
				if err != nil {
					fmt.Println("Print error:", err)
				}
			}
		}
		return nil
	})
	if err != nil {
		fmt.Println(err)
	}
}
