package http

import (
	"general-agent/config"

	"github.com/gin-gonic/gin"
)

func NewDefaultRouterGroup(conf *config.Config, engine *Engine) *gin.RouterGroup {
	return engine.Origin.Group(conf.HTTP.APIPrefix)
}
