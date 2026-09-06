package xnet

import (
	"fmt"
	"net"
	"strings"
)

const DefaultIPV4 = "0.0.0.0"
const DefaultIPV4Mask = "0.0.0.0/0"
const DefaultIPV6 = "::"
const DefaultIPV6Mask = "::/0"

type IP string

func (ip IP) IsIP4() bool {
	return IsIPv4(string(ip))
}

func (ip IP) IsIP6() bool {
	return IsIPv6(string(ip))
}

func (ip IP) Valid() bool {
	return net.ParseIP(string(ip)) != nil
}

func (ip IP) ParseIP() net.IP {
	return net.ParseIP(string(ip))
}
func (ip IP) String() string {
	return string(ip)
}
func (ip IP) Host(port int) string {
	if ip.IsIP6() {
		return fmt.Sprintf("[%s]:%d", ip.String(), port)
	}
	return fmt.Sprintf("%s:%d", ip.String(), port)
}

func (ip IP) CompressIPv6() IP {
	if !ip.IsIP6() {
		return ip
	}
	address := net.ParseIP(ip.String())
	address.IsLoopback()
	if address == nil {
		return ip
	}
	if address.To4() != nil {
		return ip
	}
	return IP(address.String())
}
func (ip IP) CheckIPOverlap(other IP) (bool, error) {
	_, ipNetA, err := net.ParseCIDR(ip.String())
	if err != nil {
		return false, err
	}
	_, ipNetB, err := net.ParseCIDR(other.String())
	if err != nil {
		return false, err
	}
	if ipNetA.Contains(ipNetB.IP) || ipNetB.Contains(ipNetA.IP) {
		return true, nil
	}
	return false, nil
}

func CompressIPv6(ipv6 string) IP {
	ip := IP(ipv6)
	if !ip.IsIP6() {
		return ip
	}
	address := net.ParseIP(ip.String())
	if address == nil {
		return ip
	}
	if address.To4() != nil {
		return ip
	}
	return IP(address.String())
}

func IsIPv4(ipAddr string) bool {
	ip := net.ParseIP(ipAddr)
	// ipv4 判断改为不包含:   兼容::ffff:10.0.0.1
	return ip != nil && !strings.Contains(ipAddr, ":")
}
func IsIPv6(ipAddr string) bool {
	ip := net.ParseIP(ipAddr)
	return ip != nil && strings.Contains(ipAddr, ":")
}

func IsIP(ipAddr string) bool {
	ip := net.ParseIP(ipAddr)
	return ip != nil
}
