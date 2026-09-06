// Package xnet
// @author: fengyi
// @date: 2024/5/6
// @note:
package xnet

import (
	"fmt"
	"testing"
)

func TestIP(t *testing.T) {
	printIpInfo("192.168.1.1")
	printIpInfo("::ffff:10.0.0.1")
	printIpInfo("2001:0db8:0000:0042:0000:8a2e:0370:7334")
}

func printIpInfo(ip IP) {
	fmt.Println("=============", ip, "=============")
	fmt.Println("IP Valid:", ip.Valid())
	fmt.Println("IP IsIP4:", ip.IsIP4())
	fmt.Println("IP IsIP6:", ip.IsIP6())
	fmt.Println("IP CompressIPv6:", ip.CompressIPv6())
}

func TestIP_CompressIPv6(t *testing.T) {
	ipv6Addresses := []IP{
		"2001:0db8:0000:0042:0000:8a2e:0370:7334",
		"2001:0db8:0000:0000:0000:0000:0000:0001",
		"2001:0db8:0000:0000:0000:0000:0000:0000",
		"0000:0000:0000:0000:0000:0000:0000:0001",
		"0000:0000:0000:0000:0000:0000:0000:0000",
		"2003:0000:0000:0000:0000:0000:0000:91",
		"2003:0000:1111:abcd:2222:6666:5555:9191",
	}

	for _, ipv6 := range ipv6Addresses {
		fmt.Printf("%s -> %s\n", ipv6, ipv6.CompressIPv6())
	}
}

func TestCompressIPv6(t *testing.T) {
	ipv6Addresses := []IP{
		"2001:0db8:0000:0042:0000:8a2e:0370:7334",
		"2001:0db8:0000:0000:0000:0000:0000:0001",
		"2001:0db8:0000:0000:0000:0000:0000:0000",
		"0000:0000:0000:0000:0000:0000:0000:0001",
		"0000:0000:0000:0000:0000:0000:0000:0000",
		"2003:0000:0000:0000:0000:0000:0000:91",
	}

	for _, ipv6 := range ipv6Addresses {
		fmt.Printf("%s -> %s\n", ipv6, ipv6.CompressIPv6())
	}
}
