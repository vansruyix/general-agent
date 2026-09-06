package middleware

import (
	"errors"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"general-agent/extension/errorx"
)

func LogError(log *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		for _, err := range c.Errors {
			if err.Type != gin.ErrorTypePrivate {
				continue
			}

			var ex errorx.Error
			if errors.As(err.Unwrap(), &ex) {
				switch ex.Level() {
				case errorx.LevelError:
					log.Error("http error", zap.Error(ex))
					// logz.Error(c.Request.Context(), "http error", logz.Err(ex))
				case errorx.LevelWarning:
					log.Warn("http error", zap.Error(ex))
					// logz.Warn(c.Request.Context(), "http warn", logz.Err(ex))
				case errorx.LevelCritical:
					log.Error("http critical", zap.Error(ex))
					// logz.Error(c.Request.Context(), "http critical", logz.Err(ex))
				case errorx.LevelInfo:
					log.Info("http info", zap.Error(ex))
					// logz.Info(c.Request.Context(), "http info", logz.Err(ex))
				case errorx.LevelDebug:
					log.Debug("http debug", zap.Error(ex))
					// logz.Debug(c.Request.Context(), "http debug", logz.Err(ex))
				}
			}
		}
	}
}
