package controller

import (
	"general-agent/app/services"
	"general-agent/common"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	service *services.UserService
}

func NewUserController(service *services.UserService) *UserController {
	return &UserController{
		service: service,
	}
}

// GetById 获取用户信息通过Id
// @Summary 获取用户信息
// @Description 根据ID查询单个用户
// @Tags 用户/用户信息
// @Success 200 {object} entity.User
// @Router /user/getById/{id} [get]
func (ctrl *UserController) GetById(ctx *gin.Context) {
	ctx.JSON(200, common.ResultResp{Data: "这是一个用户"})
}
