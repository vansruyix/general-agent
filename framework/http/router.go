package http

import (
	"general-agent/config"
	"general-agent/framework/http/middleware"

	"github.com/gin-gonic/gin"
)

func NewDefaultRouterGroup(conf *config.Config, engine *Engine, sr *middleware.StaticResource) *gin.RouterGroup {
	engine.Origin.NoRoute(middleware.NoRoute(conf, sr))
	return engine.Origin.Group(conf.HTTP.APIPrefix)
}
