package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewMysql(config *Config) (*gorm.DB, error) {
	option := config.Mysql
	param := "charset=utf8mb4&loc=Asia%2FShanghai&parseTime=true"
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?%s", option.Username, option.Password, option.Host, option.Port, option.Database, param)
	options := gorm.Config{
		TranslateError: true,
	}
	options.Logger = logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second,
			LogLevel:      TernaryOperation(true, logger.Info, logger.Silent),
			Colorful:      true,
		},
	)
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN: dsn,
	}), &options)
	if err != nil {
		return nil, err
	}
	sqlDb, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDb.SetConnMaxIdleTime(time.Duration(option.ConnMaxLifetime))
	sqlDb.SetMaxOpenConns(option.MaxOpenConns)
	sqlDb.SetMaxIdleConns(option.MaxIdleConns)
	return db, nil
}

func TernaryOperation[T any](operation bool, v1, v2 T) T {
	if operation {
		return v1
	} else {
		return v2
	}
}
