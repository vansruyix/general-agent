package controller

import (
	"general-agent/app/dao"
	"general-agent/app/services"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

var Module = fx.Module("user",
	fx.Provide(dao.NewUserModel, services.NewUserService, NewUserController),
	fx.Invoke(),
)

func registryRouter(router *gin.RouterGroup, ctrl *UserController) {
	router.GET("/getById", ctrl.GetById)
}
