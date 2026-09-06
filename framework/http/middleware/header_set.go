// Package middleware
// @author: luoyuansha
// @date: 2025/6/10
// @note:
package middleware

import (
	"general-agent/extension/contextz"

	"github.com/gin-gonic/gin"
)

// HeaderSet 请求头设置中间件
func HeaderSet(skippers ...SkipperFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		if skipHandler(c, skippers...) {
			c.Next()
			return
		}

		// 获取请求头中的租户ID、用户ID
		header := c.Request.Header
		bc := contextz.BizContext{}
		bc[contextz.XDR_TENANT_ID] = header.Get(contextz.XDR_TENANT_ID)
		//bc[contextz.XDR_USER_ID] = header.Get(contextz.XDR_USER_ID)
		c.Request = c.Request.WithContext(contextz.WithBizCtx(c.Request.Context(), bc))

		// 继续处理请求
		c.Next()
	}
}
