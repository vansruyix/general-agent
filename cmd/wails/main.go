// Package main 是 General Agent 的 Wails v3 桌面应用入口。
// 启动内嵌 HTTP 服务并 serve 前端静态文件，然后在 WebView2 中加载。
package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"general-agent/cmd"
	"general-agent/config"

	"github.com/gin-gonic/gin"
	"github.com/wailsapp/wails/v3/pkg/application"
	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewDevelopment()

	cfg, err := config.NewGlobalConfig()

	if err != nil {
		logger.Fatal("加载配置失败", zap.Error(err))
	}

	addr := fmt.Sprintf("%s:%d", cfg.HTTP.Host, cfg.HTTP.Port)
	webURL := fmt.Sprintf("http://localhost:%d", cfg.HTTP.Port)

	// fx 组装所有业务模块 + 静态文件服务
	fxApp := cmd.CreateHttpServer()
	go fxApp.Run()

	// 等待 HTTP 服务就绪（用 localhost 检测，因为 0.0.0.0 不能被 connect）
	if !waitForServer(fmt.Sprintf("127.0.0.1:%d", cfg.HTTP.Port), 5*time.Second) {
		logger.Fatal("HTTP 服务启动超时", zap.String("addr", addr))
	}

	// 创建 Wails 桌面窗口
	wailsApp := application.New(application.Options{
		Name:        "General Agent",
		Description: "通用智能体桌面应用",
	})

	window := wailsApp.Window.NewWithOptions(application.WebviewWindowOptions{
		Name:   "main",
		Title:  "General Agent",
		Width:  1200,
		Height: 800,
		URL:    webURL,
	})
	window.Show()

	wailsApp.OnShutdown(func() {
		go fxApp.Stop(context.Background())
	})

	if err := wailsApp.Run(); err != nil {
		logger.Fatal("应用启动失败", zap.Error(err))
	}
	fxApp.Stop(context.Background())
}

// registerStaticFiles 为 Gin 引擎注册前端静态文件服务和 SPA fallback。
// 查找顺序：环境变量 WAILS_FRONTEND_DIR → 当前目录/frontend/dist → 上级目录/frontend/dist
func registerStaticFiles(engine *gin.Engine) {
	distDir := os.Getenv("WAILS_FRONTEND_DIR")
	if distDir == "" {
		// 尝试多个常见位置
		candidates := []string{
			"frontend/dist",
			"../frontend/dist",
		}
		for _, d := range candidates {
			if info, err := os.Stat(d); err == nil && info.IsDir() {
				distDir = d
				break
			}
		}
	}

	if distDir == "" {
		fmt.Fprintf(os.Stderr, "警告: 找不到前端静态文件目录\n")
		return
	}
	fmt.Printf("前端静态文件目录: %s\n", distDir)

	// 静态文件中间件
	fileServer := http.FileServer(http.Dir(distDir))
	engine.Use(func(c *gin.Context) {
		p := c.Request.URL.Path
		// API 路径放行
		if len(p) >= 4 && p[:4] == "/api" {
			c.Next()
			return
		}
		// 尝试提供静态文件
		fullPath := filepath.Join(distDir, filepath.Clean(p))
		if info, err := os.Stat(fullPath); err == nil && !info.IsDir() {
			fileServer.ServeHTTP(c.Writer, c.Request)
			c.Abort()
			return
		}
		// SPA fallback
		indexPath := filepath.Join(distDir, "index.html")
		if _, err := os.Stat(indexPath); err == nil {
			c.File(indexPath)
			c.Abort()
			return
		}
		c.Next()
	})
}

// waitForServer 轮询等待 HTTP 服务就绪
func waitForServer(addr string, timeout time.Duration) bool {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		conn, err := net.DialTimeout("tcp", addr, 100*time.Millisecond)
		if err == nil {
			conn.Close()
			return true
		}
		time.Sleep(100 * time.Millisecond)
	}
	return false
}
