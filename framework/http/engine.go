package http

import (
	"fmt"
	"general-agent/config"
	"general-agent/extension/logz"
	"general-agent/extension/xfile"
	"general-agent/framework/http/middleware"
	"general-agent/framework/http/response"
	"general-agent/framework/http/session"
	"net/http"
	"os"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"

	"github.com/DeanThompson/ginpprof"
	"github.com/gin-gonic/gin"
)

type RuntimeOptions struct {
	SessionLoader     *session.SessionLoader
	GlobalMiddlewares []gin.HandlerFunc
	Render            func(ctx *gin.Context, data interface{}, err error)
}

// DefaultRuntimeOptions 预设的运行时配置
func DefaultRuntimeOptions(conf *config.Config, log *zap.Logger) *RuntimeOptions {
	rOpts := &RuntimeOptions{

		GlobalMiddlewares: []gin.HandlerFunc{
			middleware.Recovery(), // gin.Recovery(),
			middleware.CORS(),
			middleware.HeaderSet(),
			middleware.Logger(log, middleware.SkipWithPathPrefix("/healthz")),
			middleware.LogError(log),
		},
		Render: func(ctx *gin.Context, data interface{}, err error) {
			if data == nil {
				data = struct{}{}
			}
			if err != nil {
				// 直接抛出错误，交给全局异常处理中间件处理
				panic(err)
			}

			response.Ctx(ctx).Data(data)
		},
	}
	return rOpts
}

// Engine 包装后的 Engine
type Engine struct {
	// 原生的 Engine
	Origin *gin.Engine
	// SessionLoader 加载器
	sessionLoader *session.SessionLoader
	// 将响应结构渲染为body内容的函数
	render func(ctx *gin.Context, data interface{}, err error)
}

func NewDefaultEngine(conf *config.Config, log *zap.Logger) *Engine {
	runtimeOptions := DefaultRuntimeOptions(conf, log)
	return NewEngine(conf, runtimeOptions)
}

func NewEngine(conf *config.Config, rOpts *RuntimeOptions) *Engine {
	// 调试模式
	if config.IsDev() {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	ginEngine := gin.New()

	engine := &Engine{
		Origin:        ginEngine,
		sessionLoader: rOpts.SessionLoader,
		render:        rOpts.Render,
	}
	ginEngine.GET(conf.HTTP.APIPrefix+"/healthz", HealthyHandler)
	ginEngine.GET("/healthz", HealthyHandler)
	ginEngine.Static(conf.HTTP.APIPrefix+"/monitor/report/v2", xfile.AppPath+"/etc/templates/v2")
	if config.IsDev() {
		ginEngine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
		ginpprof.Wrap(ginEngine)
		logz.InfoNoCtx(fmt.Sprintf("Swagger docs http://localhost:%d/swagger/index.html", conf.HTTP.Port))
	}
	engine.Origin.Use(rOpts.GlobalMiddlewares...)

	return engine
}

// func (e *Engine) Authorization() gin.HandlerFunc {
// 	return middleware.Authorization(e.render, e.sessionLoader)
// }

func (e *Engine) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	e.Origin.ServeHTTP(writer, request)
}

// HealthyHandler health checks for the server.
// guide: https://tools.ietf.org/id/draft-inadarei-api-health-check-01.html#rfc.section.3
func HealthyHandler(ctx *gin.Context) {
	ctx.Header("Cache-Control", "max-age=3600")
	ctx.Header("Content-Type", "application/health+json")

	ctx.JSON(http.StatusOK, gin.H{
		"hostname":    hostname,
		"status":      "pwd",
		"version":     "1",
		"releaseID":   "1.0.0",
		"description": "health of scan-management service",
	})
}

var hostname string

func init() {
	hostname, _ = os.Hostname()
}
