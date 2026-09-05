// Package database 提供基于 gorm 的 MySQL 数据库连接。
// 日志通过 zapGormLogger 桥接到 zap，慢 SQL（>1s）以 Warn 级别记录，
// 查询错误以 Error 级别记录。
package database

import (
	"context"
	"fmt"
	"time"

	"general-agent/internal/config"

	"go.uber.org/fx"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

// zapGormLogger 实现 gorm logger.Interface，将 gorm 日志桥接到 zap。
type zapGormLogger struct {
	log *zap.Logger
}

// LogMode 实现 gorm logger.Interface，直接返回自身（忽略级别切换）。
func (l *zapGormLogger) LogMode(level gormlogger.LogLevel) gormlogger.Interface { return l }

// Info 实现 gorm logger.Interface，以 Info 级别记录。
func (l *zapGormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	l.log.Info(fmt.Sprintf(msg, data...))
}

// Warn 实现 gorm logger.Interface，以 Warn 级别记录。
func (l *zapGormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	l.log.Warn(fmt.Sprintf(msg, data...))
}

// Error 实现 gorm logger.Interface，以 Error 级别记录。
func (l *zapGormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	l.log.Error(fmt.Sprintf(msg, data...))
}

// Trace 实现 gorm logger.Interface，记录每条 SQL 的执行耗时。
// 查询错误以 Error 级别记录；慢查询（>1s）以 Warn 级别记录。
func (l *zapGormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	elapsed := time.Since(begin)
	sql, rows := fc()
	if err != nil {
		l.log.Error("gorm query error",
			zap.Error(err),
			zap.String("sql", sql),
			zap.Int64("rows", rows),
			zap.Duration("elapsed", elapsed),
		)
	} else if elapsed > time.Second {
		l.log.Warn("slow query",
			zap.String("sql", sql),
			zap.Int64("rows", rows),
			zap.Duration("elapsed", elapsed),
		)
	}
}

// NewDB 创建 MySQL 连接并配置连接池参数。
// 开启 TranslateError 以将 MySQL 错误转为 gorm 标准错误（如 ErrRecordNotFound、ErrDuplicatedKey）。
func NewDB(cfg *config.Config, log *zap.Logger) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&loc=Asia%%2FShanghai&parseTime=true",
		cfg.MySQL.Username, cfg.MySQL.Password, cfg.MySQL.Host, cfg.MySQL.Port, cfg.MySQL.Database)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		TranslateError: true,
		Logger:         &zapGormLogger{log: log},
	})
	if err != nil {
		return nil, fmt.Errorf("open mysql: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("get sql.DB: %w", err)
	}
	sqlDB.SetMaxIdleConns(cfg.MySQL.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.MySQL.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Duration(cfg.MySQL.ConnMaxLifetime) * time.Second)

	return db, nil
}

// Module 是 fx 模块，提供 *gorm.DB 单例。
var Module = fx.Module("database", fx.Provide(NewDB))