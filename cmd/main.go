// Package main 是应用入口，通过 fx 依赖注入组装所有模块并启动 HTTP 服务。
// 模块注册顺序：config → logger → database → errs → server → user，
// fx 自动解析依赖关系，无需手动排序。
package main

import (
	"context"
	"general-agent/ai"
	"general-agent/app"
	"general-agent/framework"
	"general-agent/framework/http"
	"general-agent/framework/scheduler"

	_ "general-agent/docs" // swagger docs

	"go.uber.org/fx"
)

// @title           General Agent
// @version         1.0
// @description     这是一个通用Agent智能体.
// @host            localhost:8080
// @BasePath        /api/v1/general-agent
func main() {
	App := fx.New(
		// config.Module,
		// logger.Module,
		// database.Module,
		// server.Module,
		ai.Module,
		app.Module,
		framework.Module,
		fx.Invoke(start),
		fx.Invoke(task),
	)
	App.Run()
}

func start(lifecycle fx.Lifecycle, httpServer *http.HTTPServer) {
	lifecycle.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				go httpServer.Serve(ctx)
				return nil
			},
			OnStop: func(ctx context.Context) error {
				go httpServer.Shutdown(ctx)
				return nil
			},
		},
	)
}

func task(lifecycle fx.Lifecycle, scheduler *scheduler.JobManager) {
	lifecycle.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				go scheduler.Start()
				return nil
			},
			OnStop: func(ctx context.Context) error {
				scheduler.Stop()
				return nil
			},
		},
	)
}
