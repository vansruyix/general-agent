package response

import (
	"errors"
	"general-agent/config"
	"general-agent/extension/errorx"
	"strings"

	"github.com/gin-gonic/gin"
)

type ginResponse struct {
	ctx *gin.Context
}

func Ctx(c *gin.Context) *ginResponse {
	return &ginResponse{ctx: c}
}

// CodeMsg 自定义错误码
func (res *ginResponse) CodeMsg(code int, msg string) {
	res.ctx.JSON(200, Response{Code: code, Msg: msg})
}

// Ok 成功响应
func (res *ginResponse) Ok() {
	res.ctx.JSON(200, Response{Code: errorx.DefaultSuccessCode, Msg: ResponseDefaultMsg})
}

// Data 成功数据响应
func (res *ginResponse) Data(data any) {
	res.ctx.JSON(200, Response{Code: errorx.DefaultSuccessCode, Msg: ResponseDefaultMsg, Data: data})
}

// Custom 自定义响应数据
func (res *ginResponse) Custom(data any) {
	res.ctx.JSON(200, data)
}

// Fail 未定义错误码的错误，默认错误码-1
func (res *ginResponse) Fail(msg string) {
	res.ctx.JSON(200, Response{Code: errorx.DefaultErrorCode, Msg: msg})
}

// Error 根据业务错误响应（推荐）
func (res *ginResponse) Error(err error) {
	var bizErr errorx.Error
	// 判断是否是业务错误，不是则转换为业务错误
	if errorx.IsBizFault(err) {
		bizErr = errorx.New()
		errors.As(err, &bizErr)
	} else {
		bizErr = errorx.ErrUnknown.WithReason(err.Error())
	}
	// 国际化处理
	if strings.HasPrefix(bizErr.Message, "_in18.") {
		// 如果是Message是约定的国际化字符，则在这里转换为对应语言
		lang := res.ctx.GetString("lang")
		bizErr.Message = langConvert(lang, bizErr.Message)
	}
	resBody := Response{Code: bizErr.Code, Msg: bizErr.Message}
	if !config.IsProd() {
		resBody.Reason = bizErr.GetReason()
	}
	// 返回数据
	res.ctx.JSON(200, resBody)
}

func langConvert(lang, key string) string {
	//TODO 语言切换，控制面与引擎的分开。否则点错将产生不同语言的告警日志数据。
	return key
}
