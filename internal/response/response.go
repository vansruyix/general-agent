package response

import (
	"general-agent/internal/errs"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response 是统一的 JSON 响应结构体。
// Code 为 200 表示成功，非 200 表示业务错误码；Data 仅在成功时填充。
type Response struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data,omitempty"`
}

// OK 返回 HTTP 200 + code:200 的成功响应。
func OK(ctx *gin.Context, data any) {
	ctx.JSON(http.StatusOK, Response{Code: 0, Msg: "success", Data: data})
}

// Fail 根据 error 类型返回不同响应：
//   - *BizError：HTTP 200 + 对应的 code/msg（前端按 code 分支处理业务错误）
//   - 其他 error：zap 记录错误日志后返回 HTTP 500 + code:50000，不泄漏内部细节
func Fail(ctx *gin.Context, err error) {
	if bizErr, ok := err.(*errs.BizError); ok {
		ctx.JSON(http.StatusOK, Response{Code: bizErr.Code, Msg: bizErr.Msg})
		return
	}
	ctx.JSON(http.StatusInternalServerError, Response{Code: errs.ErrInternal.Code, Msg: errs.ErrInternal.Msg})
}
