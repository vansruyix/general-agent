package middleware

import (
	"general-agent/config"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func NoRoute(conf *config.Config, sr *StaticResource) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// API 路径不存在时返回 JSON，而不是前端页面
		if strings.HasPrefix(ctx.Request.URL.Path, conf.HTTP.APIPrefix) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			ctx.Abort()
			return
		}

		ctx.FileFromFS("index.html", sr.HttpFS)
		ctx.Abort()
	}
}
