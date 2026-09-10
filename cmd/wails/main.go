// Package main 是 General Agent 的 Wails v3 桌面应用入口。
// 启动内嵌 HTTP 服务，同时 serve 前端静态文件，然后在 WebView2 中加载。
// 与 cmd/main.go 共享所有业务模块（ai、app、framework），仅入口方式不同。
package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"general-agent/ai"
	"general-agent/app"
	"general-agent/config"
	"general-agent/framework"
	httpfw "general-agent/framework/http"

	"github.com/gin-gonic/gin"
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

	// fx 组装所有业务模块 + 静态文件服务
	fxApp := fx.New(
		ai.Module,
		app.Module,
		framework.Module,
		fx.Invoke(func(lc fx.Lifecycle, engine *httpfw.Engine, log *zap.Logger) {
			// 注册前端静态文件服务（Wails 模式下需要 serve 构建产物）
			setupStaticFiles(engine.Origin)

			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					go func() {
						addr := fmt.Sprintf("%s:%d", cfg.HTTP.Host, cfg.HTTP.Port)
						srv := &http.Server{Addr: addr, Handler: engine.Origin}
						log.Info("HTTP 服务已启动", zap.String("addr", addr))
						if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
							log.Fatal("HTTP 服务启动失败", zap.Error(err))
						}
					}()
					return nil
				},
				OnStop: func(ctx context.Context) error {
					return nil
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

// setupStaticFiles 注册前端静态文件服务。
// 优先使用环境变量 WAILS_FRONTEND_DIR 指定的路径，默认为 frontend/dist/
func setupStaticFiles(engine *gin.Engine) {
	distDir := os.Getenv("WAILS_FRONTEND_DIR")
	if distDir == "" {
		// 从可执行文件所在目录查找 frontend/dist
		execPath, _ := os.Executable()
		distDir = filepath.Join(filepath.Dir(execPath), "frontend", "dist")
	}

	// 如果 dist 目录不存在，回退到项目相对路径（开发调试用）
	if _, err := os.Stat(distDir); os.IsNotExist(err) {
		distDir = "frontend/dist"
	}

	// 为 Gin 添加静态文件路由
	engine.Use(func(c *gin.Context) {
		filePath := filepath.Join(distDir, c.Request.URL.Path)
		if info, err := os.Stat(filePath); err == nil && !info.IsDir() {
			http.ServeFile(c.Writer, c.Request, filePath)
			c.Abort()
			return
		}
		c.Next()
	})

	// SPA fallback: 非 API 路由返回 index.html
	engine.NoRoute(func(c *gin.Context) {
		indexPath := filepath.Join(distDir, "index.html")
		if _, err := os.Stat(indexPath); err == nil {
			c.File(indexPath)
		}
	})
}