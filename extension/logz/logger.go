package logz

import (
	"context"
	"fmt"
	"general-agent/extension/contextz"
	"log/slog"
	"os"
	"runtime"
	"strings"
)

var (
	l *slog.Logger

	Group    = slog.Group
	String   = slog.String
	Int64    = slog.Int64
	Int      = slog.Int
	Uint64   = slog.Uint64
	Float64  = slog.Float64
	Bool     = slog.Bool
	Time     = slog.Time
	Duration = slog.Duration
	Any      = slog.Any
)

// 获取调用者信息
func getSourceInfo() slog.Attr {
	// Skip 3 frames:
	// 0: runtime.Caller
	// 1: getSourceInfo()
	// 2: log function (Info/Debug/Error etc.)
	// 3: actual caller
	_, file, line, ok := runtime.Caller(3)
	if !ok {
		return slog.String("source", "unknown")
	}

	// 简化路径，只保留最后4级目录+文件名
	var simplifiedPath strings.Builder

	// 分割路径
	parts := strings.Split(file, string(os.PathSeparator))
	if len(parts) == 1 {
		// 处理Windows风格路径
		parts = strings.Split(file, "/")
	}

	// 计算起始索引，最多保留最后4级目录+文件名
	startIndex := len(parts) - 5 // 4级目录 + 1个文件名
	if startIndex < 0 {
		startIndex = 0
	}

	// 处理路径的每一部分
	for i := startIndex; i < len(parts); i++ {
		part := parts[i]
		if i == len(parts)-1 {
			// 最后一部分是文件名，完整保留
			simplifiedPath.WriteString(part)
		} else if len(part) > 0 {
			// 目录名只保留首字母
			simplifiedPath.WriteString(string(part[0]))
			simplifiedPath.WriteString("/")
		} else {
			// 空部分（比如路径开始的空字符串）
			// 只有在不是起始位置时才添加
			if i > startIndex {
				simplifiedPath.WriteString("/")
			}
		}
	}

	source := fmt.Sprintf("%s:%d", simplifiedPath.String(), line)
	return slog.String("source", source)
}

func init() {
	//if model.IsDev() {
	//	l = slog.New(slog.NewTextHandler(os.Stdout, nil))
	//} else {
	//	exePath, err := os.Executable()
	//	if err != nil {
	//		panic(fmt.Errorf("[logz] %v", err))
	//	}
	//	flog, err := os.OpenFile(filepath.Join(filepath.Dir(exePath), "wserver.log"), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	//	if err != nil {
	//		panic(fmt.Errorf("[logz] error opening file: %v", err))
	//	}
	//	l = slog.New(slog.NewTextHandler(io.MultiWriter(os.Stdout, flog), nil))
	//}
	opts := &slog.HandlerOptions{
		AddSource: false, // 我们自己处理源信息
	}
	l = slog.New(slog.NewTextHandler(os.Stdout, opts))
}

func Debug(ctx context.Context, msg string, args ...any) {
	args = append(args, getSourceInfo())
	l.DebugContext(ctx, prefixWithModuleName(ctx, msg), args...)
}

func DebugNoCtx(msg string, args ...any) {
	args = append(args, getSourceInfo())
	l.Debug(msg, args...)
}

func Info(ctx context.Context, msg string, args ...any) {
	args = append(args, getSourceInfo())
	l.InfoContext(ctx, prefixWithModuleName(ctx, msg), args...)
}

func InfoNoCtx(msg string, args ...any) {
	args = append(args, getSourceInfo())
	l.Info(msg, args...)
}

func Warn(ctx context.Context, msg string, args ...any) {
	args = append(args, getSourceInfo())
	l.WarnContext(ctx, prefixWithModuleName(ctx, msg), args...)
}

func WarnNoCtx(msg string, args ...any) {
	args = append(args, getSourceInfo())
	l.Warn(msg, args...)
}

func Error(ctx context.Context, msg string, args ...any) {
	args = append(args, getSourceInfo())
	if len(args) != 0 {
		if _, ok := args[0].(error); ok {
			args[0] = slog.Any("error", args[0])
		}
	}
	l.ErrorContext(ctx, prefixWithModuleName(ctx, msg), args...)
}

func ErrorNoCtx(msg string, args ...any) {
	args = append(args, getSourceInfo())
	if len(args) != 0 {
		if _, ok := args[0].(error); ok {
			args[0] = slog.Any("error", args[0])
		}
	}
	l.Error(msg, args...)
}

func Err(err error) slog.Attr {
	return slog.Any("error", err)
}

func prefixWithModuleName(ctx context.Context, msg string) string {
	module := contextz.ModuleName(ctx)
	if module == "" {
		return msg
	}
	return fmt.Sprintf("[%s] %s", module, msg)
}
