// Package xnet
// @author: fengyi
// @date: 2024/8/9
// @note:
package xnet

import (
	"fmt"
	"net"
	"strings"
)

type IPLabel struct {
	Label string
	IP    string // 支持形式 172.18.100.2、192.168.1.1/32、0.0.0.0/0
}

func CheckIPLabelOverlap(ips []IPLabel) (IPLabel, IPLabel, error) {
	for i := 0; i < len(ips); i++ {
		_, ipNetA, err := net.ParseCIDR(ips[i].IP)
		if err != nil {
			ipNetA = parseSingleIP(ips[i].IP)
			if ipNetA == nil {
				return IPLabel{}, IPLabel{}, err
			}
		}

		for j := i + 1; j < len(ips); j++ {
			_, ipNetB, err := net.ParseCIDR(ips[j].IP)
			if err != nil {
				ipNetB = parseSingleIP(ips[j].IP)
				if ipNetB == nil {
					return IPLabel{}, IPLabel{}, err
				}
			}

			if ipNetA.Contains(ipNetB.IP) || ipNetB.Contains(ipNetA.IP) {
				return ips[i], ips[j], fmt.Errorf("%s (%s) 与 %s (%s) 存在包含关系", ips[i].Label, ips[i].IP, ips[j].Label, ips[j].IP)
			}
		}
	}
	return IPLabel{}, IPLabel{}, nil
}

func CheckIPOverlap(ips []string) (string, string, error) {
	for i := 0; i < len(ips); i++ {
		_, ipNetA, err := net.ParseCIDR(ips[i])
		if err != nil {
			ipNetA = parseSingleIP(ips[i])
			if ipNetA == nil {
				return "", "", err
			}
		}

		for j := i + 1; j < len(ips); j++ {
			_, ipNetB, err := net.ParseCIDR(ips[j])
			if err != nil {
				ipNetB = parseSingleIP(ips[j])
				if ipNetB == nil {
					return "", "", err
				}
			}

			if ipNetA.Contains(ipNetB.IP) || ipNetB.Contains(ipNetA.IP) {
				return ips[i], ips[j], fmt.Errorf("%s 与 %s 存在包含关系", ips[i], ips[j])
			}
		}
	}
	return "", "", nil
}

// parseSingleIP 用于处理单个 IP（不带网段）转化为 /32 的形式
func parseSingleIP(ipStr string) *net.IPNet {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return nil
	}
	if strings.Contains(ipStr, ":") {
		return &net.IPNet{IP: ip, Mask: net.CIDRMask(128, 128)}
	} else {
		return &net.IPNet{IP: ip, Mask: net.CIDRMask(32, 32)}
	}
}
