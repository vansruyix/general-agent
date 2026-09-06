package middleware

import (
	"bytes"
	"io"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const MAX_LOG_CONTENT_LENGTH = 1024

func Logger(log *zap.Logger, skippers ...SkipperFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		if skipHandler(c, skippers...) {
			c.Next()
			return
		}

		start := time.Now()
		logRequest(c, log)
		c.Next()
		logResponse(c, start, log)
	}
}

func logRequest(c *gin.Context, log *zap.Logger) {
	body := getPrintableBody(c)
	query := c.Request.URL.RawQuery

	args := []zap.Field{
		zap.Any("method", c.Request.Method),
		zap.Any("uri", c.Request.URL.Path),
		// logz.Any("content_type", c.Request.Header.Get("Content-Type")),
		// logz.Any("content_length", c.Request.ContentLength),
		// logz.Any("headers", c.Request.Header),
		// logz.Any("ip", c.ClientIP()),
		// logz.Any("ua", c.Request.Header.Get("User-Agent")),
	}
	if query != "" {
		args = append(args, zap.Any("query", query))
	}
	if body != "" {
		args = append(args, zap.Any("body", body))
	}
	// logz.Info(c.Request.Context(), "[HTTP Request ]", args...)
	log.Info("[HTTP Request ]",
		args...,
	)
}

func logResponse(c *gin.Context, start time.Time, log *zap.Logger) {
	// logz.Info(c.Request.Context(), "[HTTP Response]",
	// 	logz.Any("method", c.Request.Method),
	// 	logz.Any("uri", c.Request.URL.Path),
	// 	logz.Any("status", c.Writer.Status()),
	// 	logz.Any("latency", time.Since(start).String()),
	// )
	bytes := make([]byte, 0)
	log.Info("[HTTP Response]",
		zap.String("method", c.Request.Method),
		zap.String("url", c.Request.URL.Path),
		zap.Int("status", c.Writer.Status()),
		zap.String("body", string(bytes)),
		zap.Duration("latency", time.Since(start)),
	)
}

func getPrintableBody(c *gin.Context) string {
	if c.ContentType() == "application/json" || c.ContentType() == "application/x-www-form-urlencoded" {
		return string(readBody(c))
	}
	return c.ContentType()
}

func readBody(c *gin.Context) []byte {
	body, _ := c.GetRawData()
	c.Request.Body = io.NopCloser(bytes.NewBuffer(body))

	return body
}
