// Package main 是应用入口，通过 fx 依赖注入组装所有模块并启动 HTTP 服务。
// 模块注册顺序：config → logger → database → errs → server → user，
// fx 自动解析依赖关系，无需手动排序。
package main

import (
	"general-agent/app/user"
	"general-agent/internal/config"
	"general-agent/internal/database"
	"general-agent/internal/logger"
	"general-agent/internal/server"

	_ "general-agent/docs" // swagger docs

	"go.uber.org/fx"
)

// @title           General Agent
// @version         1.0
// @description     这是一个通用Agent智能体.
// @host            localhost:8080
// @BasePath        /api/v1/general-agent
func main() {
	fx.New(
		config.Module,
		logger.Module,
		database.Module,
		server.Module,
		user.Module,
	).Run()
}
