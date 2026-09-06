// Package xnet
// @author: liulun
// @date: 2024/5/27
// @note: 网口工具
package xnet

import (
	"log"
	"net"
)

// GetAllInterfaceIp 获取所有网卡ip地址
func GetAllInterfaceIp() (*[]string, *[]string, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, nil, err
	}
	var ipv4s = make([]string, 0)
	var ipv6s = make([]string, 0)
	// 遍历每个网络接口
	for _, iface := range interfaces {
		// 获取接口的地址列表
		addrs, err := iface.Addrs()
		if err != nil {
			log.Printf("Error getting addresses for interface %s: %v", iface.Name, err)
			continue
		}

		// 遍历每个地址
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}

			// 打印 IPv4 和 IPv6 地址
			if ip == nil {
				continue
			}
			if ip.To4() != nil {
				ipv4s = append(ipv4s, ip.String())
			} else {
				ipv6s = append(ipv6s, ip.String())
			}
		}
	}
	return &ipv4s, &ipv6s, nil
}

func CalcNetworkAddress(ipNetStr string) (string, int, error) {
	// 解析带掩码的IP地址
	_, ipNet, err := net.ParseCIDR(ipNetStr)
	if err != nil {
		return "", 0, err
	}

	// 计算网络地址
	ip := ipNet.IP.Mask(ipNet.Mask)
	ones, _ := ipNet.Mask.Size()
	return ip.String(), ones, nil
}
