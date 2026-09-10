// Package main 是 General Agent 的 Wails v3 桌面应用入口。
// 启动内嵌 HTTP 服务后在 WebView2（Windows）/ WebKit（Linux）中加载前端 UI。
// 与 cmd/main.go 共享所有业务模块（ai、app、framework），仅入口方式不同。
package main

import (
	"context"
	"fmt"

	"general-agent/ai"
	"general-agent/app"
	"general-agent/config"
	"general-agent/framework"
	httpfw "general-agent/framework/http"

	"github.com/wailsapp/wails/v3/pkg/application"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewDevelopment()

	// 加载配置
	cfg, err := config.NewGlobalConfig()
	if err != nil {
		logger.Fatal("加载配置失败", zap.Error(err))
	}

	// fx 组装所有业务模块
	fxApp := fx.New(
		ai.Module,
		app.Module,
		framework.Module,
		fx.Invoke(func(lc fx.Lifecycle, httpServer *httpfw.HTTPServer, log *zap.Logger) {
			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					go httpServer.Serve(ctx)
					log.Info("HTTP 服务已启动",
						zap.String("addr", fmt.Sprintf("%s:%d", cfg.HTTP.Host, cfg.HTTP.Port)),
					)
					return nil
				},
				OnStop: func(ctx context.Context) error {
					return httpServer.Shutdown(ctx)
				},
			})
		}),
	)

	// 启动 fx 应用
	go fxApp.Run()

	// 创建 Wails 桌面窗口，指向本地 HTTP 服务
	wailsApp := application.New(application.Options{
		Name:        "General Agent",
		Description: "通用智能体桌面应用",
	})

	window := wailsApp.Window.NewWithOptions(application.WebviewWindowOptions{
		Name:   "main",
		Title:  "General Agent",
		Width:  1200,
		Height: 800,
		URL:    fmt.Sprintf("http://localhost:%d", cfg.HTTP.Port),
	})
	window.Show()

	wailsApp.OnShutdown(func() {
		_ = fxApp.Stop(context.Background())
	})

	if err := wailsApp.Run(); err != nil {
		logger.Fatal("应用启动失败", zap.Error(err))
	}
}