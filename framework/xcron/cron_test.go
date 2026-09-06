// Package cron
// @author: fengyi
// @date: 2024/4/8
// @note:
package xcron_test

import (
	"fmt"
	"general-agent/framework/xcron"
	"testing"

	"github.com/robfig/cron/v3"
)

func TestNewCron(t *testing.T) {
	c := cron.New(cron.WithSeconds())

	c.AddFunc("@every 5s", func() {
		fmt.Println("tick every 5 second")
	})

	c.AddFunc("0 0/1 * * * ?", func() {
		fmt.Println("tick every 0 0/1 * * * ? second")
	})

	c.Start()
	for _, entry := range c.Entries() {
		fmt.Println(entry.ID)
	}

	// 启动后添加
	c.AddFunc("0 0/1 * * * ?", func() {
		fmt.Println("tick every 10 second======")
	})

	select {}
}

// 测试cron表达式生成
func TestCron(t *testing.T) {
	timer := xcron.Timer{
		TriggerTime: "08:10",
		Type:        "day",
		Week:        []int{1, 2},
		Month:       []int{1, 15},
	}

	fmt.Println(timer.ToCron())
}
