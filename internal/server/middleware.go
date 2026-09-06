// Package server 提供 gin HTTP 引擎及统一的中间件链。
// 所有中间件在此集中定义，NewEngine 创建引擎时按序注册。
package server

import (
	"general-agent/internal/errs"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// Recovery 返回一个 panic 恢复中间件。
// 发生 panic 时记录 zap 错误日志与堆栈，并返回 HTTP 500。
func Recovery(log *zap.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				log.Error("panic recovered",
					zap.Any("panic", r),
					zap.String("stack", ""),
				)
				ctx.AbortWithStatusJSON(500, gin.H{
					"code": errs.ErrInternal.Code,
					"msg":  errs.ErrInternal.Msg,
				})
			}
		}()
		ctx.Next()
	}
}

// RequestID 为每个请求注入唯一 X-Request-ID 响应头，便于链路追踪。
func RequestID() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestID := ctx.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}
		ctx.Header("X-Request-ID", requestID)
		ctx.Next()
	}
}

// CORS 返回跨域资源共享中间件，允许常见 HTTP 方法与请求头。
func CORS() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Header("Access-Control-Allow-Origin", "*")
		ctx.Header("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS,PATCH")
		ctx.Header("Access-Control-Allow-Headers", "Origin,Content-Type,Accept,Authorization,X-Request-ID,X-Requested-With")
		ctx.Header(
			"Access-Control-Expose-Headers",
			"Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type",
		)
		ctx.Header("Access-Control-Max-Age", "86400")
		ctx.Header("Access-Control-Allow-Credentials", "true")

		if ctx.Request.Method == http.MethodOptions {
			ctx.AbortWithStatus(204)
			return
		}
		ctx.Next()
	}
}

// AccessLog 返回基于 zap 的 HTTP 访问日志中间件，记录 method、path、status、latency。
func AccessLog(log *zap.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()
		ctx.Next()
		log.Info("access",
			zap.String("method", ctx.Request.Method),
			zap.String("path", ctx.Request.URL.Path),
			zap.Int("status", ctx.Writer.Status()),
			zap.Duration("latency", time.Since(start)),
		)
	}
}

// SecurityHeaders 为响应添加常见安全 HTTP 头。
func SecurityHeaders() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Header("X-Content-Type-Options", "nosniff")
		ctx.Header("X-Frame-Options", "DENY")
		ctx.Header("X-XSS-Protection", "1; mode=block")
		ctx.Next()
	}
}
