// Package xnet
// @author: fengyi
// @date: 2024/8/9
// @note:
package xnet

import (
	"fmt"
	"testing"
)

func TestCheckIPLabelOverlap(t *testing.T) {
	ips := []IPLabel{
		{Label: "A", IP: "2001:db8::1"},
		{Label: "B", IP: "2001:db8::1/128"},
		{Label: "C", IP: "2001:db8::/32"},
		{Label: "D", IP: "::/0"},
	}

	ips2 := []IPLabel{
		{Label: "A", IP: "172.18.100.2"},
		{Label: "B", IP: "192.168.1.1/32"},
		{Label: "C", IP: "172.18.0.0/16"},
		{Label: "D", IP: "0.0.0.0/0"},
	}

	_, _, err := CheckIPLabelOverlap(ips)
	fmt.Println(err)
	_, _, err = CheckIPLabelOverlap(ips2)
	fmt.Println(err)
}
