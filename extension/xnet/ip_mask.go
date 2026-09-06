// Package xnet
// @author: fengyi
// @date: 2024/6/6
// @note:
package xnet

import (
	"fmt"
	"net"
	"strconv"
	"strings"
)

// IPMask 带掩码的IP
type IPMask string

func (ipMask IPMask) String() string {
	return string(ipMask)
}
func (ipMask IPMask) Valid() bool {
	_, _, err := net.ParseCIDR(string(ipMask))
	return err == nil
}

func (ipMask IPMask) IsIP4() bool {
	return IsIPv4Mask(string(ipMask))
}

func (ipMask IPMask) IsIP6() bool {
	return IsIPv6Mask(string(ipMask))
}

func (ipMask IPMask) CompressIPv6() IPMask {
	if !ipMask.IsIP6() {
		return ipMask
	}
	ip, ipNet, err := net.ParseCIDR(string(ipMask))
	if err != nil {
		return ipMask
	}
	// 如果是ipv4直接返回
	if ip.To4() != nil {
		return ipMask
	}
	maskSize, _ := ipNet.Mask.Size()
	return IPMask(fmt.Sprintf("%s/%d", ip.String(), maskSize))
}
func (ipMask IPMask) IP() string {
	split := strings.Split(string(ipMask), "/")
	return split[0]
}

func (ipMask IPMask) Mask() int {
	split := strings.Split(string(ipMask), "/")
	if len(split) == 2 {
		atoi, err := strconv.Atoi(split[1])
		if err != nil {
			return 0
		}
		return atoi
	}
	return 0
}

func CompressIPv6Mask(ipMaskStr string) IPMask {
	ipMask := IPMask(ipMaskStr)
	return ipMask.CompressIPv6()
}

func IsIPv4Mask(ipAddr string) bool {
	_, _, err := net.ParseCIDR(ipAddr)
	// ipv4 判断改为不包含:   兼容::ffff:10.0.0.1
	return err == nil && !strings.Contains(ipAddr, ":")
}
func IsIPv6Mask(ipAddr string) bool {
	_, _, err := net.ParseCIDR(ipAddr)
	return err == nil && strings.Contains(ipAddr, ":")
}
