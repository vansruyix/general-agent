// Package logger 提供基于 zap 的日志器，通过 fx 注入为全局单例。
// 日志级别与输出格式（console / json）由 config.Log 控制，
// 应用退出时通过 fx lifecycle OnStop 自动 Sync 缓冲区。
package logger

import (
	"context"
	"os"

	"general-agent/internal/config"

	"go.uber.org/fx"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// New 根据配置创建 *zap.Logger。
// 日志级别不可解析时默认 InfoLevel；输出到 stdout。
func New(cfg *config.Config) (*zap.Logger, error) {
	level, err := zapcore.ParseLevel(cfg.Log.Level)
	if err != nil {
		level = zapcore.InfoLevel
	}

	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.TimeKey = "time"
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	var encoder zapcore.Encoder
	if cfg.Log.Encoding == "json" {
		encoder = zapcore.NewJSONEncoder(encoderCfg)
	} else {
		encoder = zapcore.NewConsoleEncoder(encoderCfg)
	}

	core := zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), zap.NewAtomicLevelAt(level))
	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
	return logger, nil
}

// Module 是 fx 模块，提供 *zap.Logger 单例，并在 OnStop 时 Sync。
var Module = fx.Module("logger",
	fx.Provide(New),
	fx.Invoke(func(lc fx.Lifecycle, log *zap.Logger) {
		lc.Append(fx.Hook{OnStop: func(ctx context.Context) error { return log.Sync() }})
	}),
)
