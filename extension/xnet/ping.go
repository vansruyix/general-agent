// Package xnet
// @author: liulun
// @date: 2024/6/3
// @note: ping命令
package xnet

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

func CanPing(target string) (bool, error) {
	// 构建 ping 命令
	cmd := exec.Command("timeout", "7", "ping", "-c", "3", "-W", "2", target)

	// 捕获命令输出
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	// 运行命令
	err := cmd.Run()
	if err != nil {
		// 如果命令执行失败，返回 false
		return false, fmt.Errorf("error executing ping command: %v", err)
	}

	// 解析输出以确定是否接收到响应
	output := out.String()
	if strings.Contains(output, "0% packet loss") {
		return true, nil
	}
	return false, nil
}
