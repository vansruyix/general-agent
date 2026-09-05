package user

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

// Module 是 user 模块的 fx 定义，注册 Repository、Service、Controller 三个依赖，
// 并通过 fx.Invoke 将路由挂载到 server 提供的 *gin.RouterGroup 上。
// 新增模块时复制此文件即可，只需修改路由注册函数。
var Module = fx.Module("user",
	fx.Provide(NewRepository, NewService, NewController),
	fx.Invoke(registerRoutes),
)

// registerRoutes 将 user 模块的 5 个 REST 端点注册到 basePath 路由组下。
func registerRoutes(rg *gin.RouterGroup, ctrl *Controller) {
	rg.GET("/user/:id", ctrl.GetByID)
	rg.GET("/user", ctrl.List)
	rg.POST("/user", ctrl.Create)
	rg.PUT("/user/:id", ctrl.Update)
	rg.DELETE("/user/:id", ctrl.Delete)
}