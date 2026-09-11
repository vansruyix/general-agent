// Package main 是应用入口，通过 fx 依赖注入组装所有模块并启动 HTTP 服务。
// 模块注册顺序：config → logger → database → errs → server → user，
// fx 自动解析依赖关系，无需手动排序。
package main

import (
	"general-agent/cmd"
	_ "general-agent/docs" // swagger docs
)

// @title           General Agent
// @version         1.0
// @description     这是一个通用Agent智能体.
// @host            localhost:8080
// @BasePath        /api/v1/general-agent
func main() {
	cmd.CreateHttpServer().Run()
}
