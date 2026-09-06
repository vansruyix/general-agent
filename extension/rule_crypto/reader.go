package rule_crypto

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// ReadIniFile 读取 GBK 编码的 INI 文件，解析为数组格式
func ReadIniFile(filename string, decoderMap map[string]Decoder) ([]map[string]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return nil, err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(file)

	data := make([]map[string]string, 0)
	entry := make(map[string]string)

	scanner := bufio.NewScanner(file)
	//encoder := simplifiedchinese.GBK.NewDecoder()

	for scanner.Scan() {
		//str, _, err := transform.Bytes(encoder, scanner.Bytes())
		//if err != nil {
		//	fmt.Println("Error decoding GBK:", err)
		//	continue
		//}
		if scanner.Err() != nil {
			fmt.Printf("scanner.Err().Error(): %v\n", scanner.Err().Error())
		}
		line := string(scanner.Bytes())
		line = strings.TrimSpace(line)

		// 忽略空行
		if line == "" {
			continue
		}

		// 检查是否是节（section）
		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			// 如果当前有未保存的条目，则保存到 data 中
			if len(entry) > 0 {
				data = append(data, entry)
				entry = make(map[string]string)
			}
			continue
		}

		// 解析键值对
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			dc, ok := decoderMap[key]
			if ok {
				v, err := dc.Decode(value)
				if err != nil {
					return nil, err
				}
				entry[key] = v
			} else {
				entry[key] = value
			}

		}
	}

	// 保存最后一个节（section）的条目
	if len(entry) > 0 {
		data = append(data, entry)
	}

	return data, nil
}

type Decoder interface {
	Decode(plat string) (string, error)
}

type Int16Decoder struct{}

func (*Int16Decoder) Decode(hexStr string) (string, error) {
	if hexStr == "" {
		return "", nil
	}
	decimal, err := strconv.ParseInt(strings.TrimPrefix(hexStr, "0x"), 16, 64)
	if err != nil {
		fmt.Println("Error:", err)
		return "", err
	}
	return strconv.FormatInt(decimal, 10), err
}
