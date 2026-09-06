// Package xnet
// @author: liulun
// @date: 2024/4/26
// @note: url地址解析
package xnet

import (
	"fmt"
	"net/url"
	"strconv"
)

func GetUrlHostAndPort(rawURL string) (string, int, error) {
	// 解析 URL
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		fmt.Println("解析 URL 错误:", err)
		return "", 0, err
	}

	// 获取主机名和端口
	host := parsedURL.Hostname()
	port := parsedURL.Port()

	// 如果端口为空，则根据协议类型设置默认端口号
	if port == "" {
		switch parsedURL.Scheme {
		case "https":
			port = "443"
		default:
			port = "80"
		}
	}
	p, _ := strconv.Atoi(port)
	return host, p, nil
}
