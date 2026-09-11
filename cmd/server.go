package cmd

import (
	"context"
	"general-agent/ai"
	"general-agent/app"
	"general-agent/framework"
	"general-agent/framework/http"
	"general-agent/framework/scheduler"

	"go.uber.org/fx"
)

func CreateHttpServer() *fx.App {
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
	return App
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
