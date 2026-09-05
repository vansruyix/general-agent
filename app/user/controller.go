package user

import (
	"general-agent/internal/response"

	"github.com/gin-gonic/gin"
)

// Controller 是 User 模块的 HTTP 处理器，依赖 Service 完成业务逻辑。
type Controller struct {
	svc *Service
}

// NewController 创建 Controller 实例。
func NewController(svc *Service) *Controller {
	return &Controller{svc: svc}
}

// GetByID godoc
// @Summary      获取用户信息
// @Description  根据ID查询单个用户
// @Tags         用户
// @Param        id  path  string  true  "用户ID"
// @Success      200  {object}  response.Response{data=UserResp}
// @Router       /user/{id} [get]
func (ctrl *Controller) GetByID(ctx *gin.Context) {
	id := ctx.Param("id")
	user, err := ctrl.svc.GetByID(id)
	if err != nil {
		response.Fail(ctx, err)
		return
	}
	response.OK(ctx, toUserResp(user))
}

// List godoc
// @Summary      用户列表
// @Description  分页查询用户列表
// @Tags         用户
// @Param        pageNum   query  int  false  "页码"
// @Param        pageSize  query  int  false  "每页大小"
// @Success      200  {object}  response.Response{data=PageResp}
// @Router       /user [get]
func (ctrl *Controller) List(ctx *gin.Context) {
	var req ListUserReq
	if err := ctx.ShouldBindQuery(&req); err != nil {
		response.Fail(ctx, err)
		return
	}
	users, total, err := ctrl.svc.List(req.PageNum, req.PageSize)
	if err != nil {
		response.Fail(ctx, err)
		return
	}
	resps := make([]UserResp, len(users))
	for i := range users {
		resps[i] = toUserResp(&users[i])
	}
	response.OK(ctx, PageResp{List: resps, Total: total})
}

// Create godoc
// @Summary      创建用户
// @Description  创建新用户
// @Tags         用户
// @Param        req  body  CreateUserReq  true  "创建请求"
// @Success      200  {object}  response.Response{data=UserResp}
// @Router       /user [post]
func (ctrl *Controller) Create(ctx *gin.Context) {
	var req CreateUserReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Fail(ctx, err)
		return
	}
	user, err := ctrl.svc.Create(&req)
	if err != nil {
		response.Fail(ctx, err)
		return
	}
	response.OK(ctx, toUserResp(user))
}

// Update godoc
// @Summary      更新用户
// @Description  根据ID更新用户信息
// @Tags         用户
// @Param        id   path  string         true  "用户ID"
// @Param        req  body  UpdateUserReq  true  "更新请求"
// @Success      200  {object}  response.Response{data=UserResp}
// @Router       /user/{id} [put]
func (ctrl *Controller) Update(ctx *gin.Context) {
	id := ctx.Param("id")
	var req UpdateUserReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Fail(ctx, err)
		return
	}
	user, err := ctrl.svc.Update(id, &req)
	if err != nil {
		response.Fail(ctx, err)
		return
	}
	response.OK(ctx, toUserResp(user))
}

// Delete godoc
// @Summary      删除用户
// @Description  根据ID删除用户
// @Tags         用户
// @Param        id  path  string  true  "用户ID"
// @Success      200  {object}  response.Response
// @Router       /user/{id} [delete]
func (ctrl *Controller) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := ctrl.svc.Delete(id); err != nil {
		response.Fail(ctx, err)
		return
	}
	response.OK(ctx, nil)
}
