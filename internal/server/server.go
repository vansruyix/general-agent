// Package server 提供 gin HTTP 引擎及统一的中间件链。
package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"general-agent/config"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// NewEngine 创建并配置 gin 引擎。
// 根据 cfg.App.Env 动态设置 gin 运行模式（dev 时为 DebugMode，其余为 ReleaseMode）。
// 中间件按序注册：Recovery → RequestID → CORS → AccessLog → SecurityHeaders。
// Swagger UI 仅在开发环境（Env == "dev"）下挂载。
func NewEngine(cfg *config.Config, log *zap.Logger) (*gin.Engine, *gin.RouterGroup) {
	// 开发模式使用 DebugMode，便于调试；生产模式使用 ReleaseMode，关闭调试输出
	if cfg.App.Env == "dev" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	engine := gin.New()

	// 统一注册中间件链
	engine.Use(
		Recovery(log),
		RequestID(),
		CORS(),
		AccessLog(log),
		SecurityHeaders(),
	)

	// Swagger 仅开发环境可用，避免生产泄漏 API 文档
	if cfg.App.Env == "dev" {
		engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	rg := engine.Group(cfg.HTTP.APIPrefix)
	return engine, rg
}

// Module 是 fx 模块，提供 *gin.Engine 与 *gin.RouterGroup，并注册服务生命周期钩子。
// OnStart 异步启动 HTTP 监听，启动失败写入 Fatal；OnStop 等待 5 秒优雅关停。
var Module = fx.Module("server",
	fx.Provide(NewEngine),
	fx.Invoke(func(lc fx.Lifecycle, cfg *config.Config, engine *gin.Engine, log *zap.Logger) {
		srv := &http.Server{Addr: fmt.Sprintf("%s:%d", cfg.HTTP.Host, cfg.HTTP.Port), Handler: engine}
		lc.Append(fx.Hook{
			OnStart: func(ctx context.Context) error {
				go func() {
					log.Sugar().Infof("server start success，address %s", fmt.Sprintf("%s:%d", cfg.HTTP.Host, cfg.HTTP.Port))
					if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
						log.Fatal("server start failed", zap.Error(err))
					}
				}()
				return nil
			},
			OnStop: func(ctx context.Context) error {
				shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()
				if err := srv.Shutdown(shutdownCtx); err != nil {
					return fmt.Errorf("server shutdown: %w", err)
				}
				return nil
			},
		})
	}),
)
