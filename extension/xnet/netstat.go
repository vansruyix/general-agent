// Package xnet
// @author: fengyi
// @date: 2024/7/23
// @note:
package xnet

import (
	"bufio"
	"fmt"
	"general-agent/extension/logz"
	"net"
	"os"
	"strconv"
	"strings"
)

// CheckPortUse 校验端口是否被使用
// use return true
func CheckPortUse(port int) bool {
	listen, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return true
	}
	_ = listen.Close()
	return false
}

func ListUsePort() ([]int64, error) {
	tcpPorts, err := listPort("/proc/net/tcp")
	if err != nil {
		return nil, err
	}
	udpPorts, err := listPort("/proc/net/udp")
	if err != nil {
		return nil, err
	}
	return append(tcpPorts, udpPorts...), nil
}

func listPort(filePath string) ([]int64, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)

	var ports []int64
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		if len(fields) == 0 {
			continue
		}
		localAddrHex := fields[1]
		localPortHex := localAddrHex[strings.LastIndex(localAddrHex, ":")+1:]
		localPort, err := strconv.ParseInt(localPortHex, 16, 16)
		if err != nil {
			logz.WarnNoCtx("error parse port", logz.Err(err))
			continue
		}
		ports = append(ports, localPort)
	}
	return ports, nil
}
