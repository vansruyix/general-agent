package agent

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

var Module = fx.Module("agent",
	fx.Provide(NewChatAgentController),
	fx.Invoke(registerRoutes),
)

// registerRoutes 将 user 模块的 5 个 REST 端点注册到 basePath 路由组下。
func registerRoutes(rg *gin.RouterGroup, ctrl *ChatAgentController) {
	rg.POST("/agent/chat", ctrl.Chat)
}
