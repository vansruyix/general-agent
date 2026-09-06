// Package xnet
// @author: fengyi
// @date: 2024/5/6
// @note:
package xnet

import (
	"fmt"
	"testing"
)

func TestIPMake(t *testing.T) {
	printIIPMaskInfo("192.168.1.1/24")
	printIIPMaskInfo("::ffff:10.0.0.1/64")
	printIIPMaskInfo("2001:0db8:0000:0042:0000:8a2e:0370:7334/64")
	printIIPMaskInfo("2005::3/64")
}

func printIIPMaskInfo(ip IPMask) {
	fmt.Println("=============", ip, "=============")
	fmt.Println("IP Valid:", ip.Valid())
	fmt.Println("IP IsIP4:", ip.IsIP4())
	fmt.Println("IP IsIP6:", ip.IsIP6())
	fmt.Println("IP CompressIPv6:", ip.CompressIPv6())
}

func TestIP_CompressIPv6Mask(t *testing.T) {
	ipv6Addresses := []IPMask{
		"2001:0db8:0000:0042:0000:8a2e:0370:7334/32",
		"2001:0db8:0000:0000:0000:0000:0000:0001/64",
		"2001:0db8:0000:0000:0000:0000:0000:0000/128",
		"0000:0000:0000:0000:0000:0000:0000:0001/64",
		"0000:0000:0000:0000:0000:0000:0000:0000/128",
		"2005::3/64",
	}

	for _, ipv6 := range ipv6Addresses {
		fmt.Printf("%s -> %s\n", ipv6, ipv6.CompressIPv6())
	}
}

func TestIPMask(t *testing.T) {
	ipNetStrs := []string{
		"192.168.1.10/24",
		"10.0.0.1/8",
		"10.15.22.22/24",
		"172.16.5.4/16",
		"2001:db8::1/32",
		"2001:db8::1/128",
		"2005::3/64",
	}

	for _, ipNetStr := range ipNetStrs {
		networkAddr, mask, err := CalcNetworkAddress(ipNetStr)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		} else {
			fmt.Printf("%-20s -> %s/%d\n", ipNetStr, networkAddr, mask)
		}
	}
}
