package middleware

import (
	"errors"
	"fmt"
	"general-agent/extension/errorx"
	"general-agent/framework/http/response"
	"net"
	"os"
	"runtime/debug"
	"strings"

	"github.com/gin-gonic/gin"

	"general-agent/extension/logz"
)

func Recovery() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		defer func() {
			err := recover()
			if err == nil {
				return
			}

			switch v := err.(type) {
			case errorx.Error:
				// 符合预期的错误，可以直接返回给客户端
				logz.Warn(ctx, fmt.Sprintf("业务错误：%s", v.Error()))
				debug.PrintStack()
				response.Ctx(ctx).Error(v)
			case *errorx.Error:
				// 符合预期的错误，可以直接返回给客户端
				logz.Warn(ctx, fmt.Sprintf("业务错误：%s", v.Error()))
				debug.PrintStack()
				response.Ctx(ctx).Error(v)
			case error:
				// 一律返回服务器错误，避免返回堆栈错误给客户端，实际还可以针对系统错误做其他处理
				logz.Error(ctx, fmt.Sprintf("意外错误：%s", v.Error()))
				debug.PrintStack()
				response.Ctx(ctx).Error(v)
			default:
				// 同上
				logz.Error(ctx, fmt.Sprintf("未定义错误：%s", v))
				debug.PrintStack()
				response.Ctx(ctx).Error(errors.New(err.(string)))
			}
		}()
		ctx.Next()
	}
}

func isBrokenPipe(err any) bool {
	if ne, ok := err.(*net.OpError); ok {
		if se, ok := ne.Err.(*os.SyscallError); ok {
			if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
				return true
			}
		}
	}
	return false
}
