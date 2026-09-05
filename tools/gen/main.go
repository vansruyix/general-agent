// Package main 是 gorm gen 实体生成器入口。
// 从 MySQL 数据库读取表结构，自动生成 gorm 实体文件到 app/<domain>/ 目录。
// 使用方式：make gen 或 go run ./tools/gen
// 需要先配置 config/dev.yaml 中的 MySQL 连接信息。
package main

import (
	"fmt"
	"log"

	"general-agent/internal/config"

	"gorm.io/driver/mysql"
	"gorm.io/gen"
	"gorm.io/gorm"
)

// targets 定义表名到输出目录的映射，新增表只需在此添加一行。
var targets = map[string]string{
	"user": "./app/user",
}

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.MySQL.Username, cfg.MySQL.Password, cfg.MySQL.Host, cfg.MySQL.Port, cfg.MySQL.Database)

	db, err := gorm.Open(mysql.Open(dsn))
	if err != nil {
		log.Fatalf("connect mysql: %v", err)
	}

	for tableName := range targets {
		g := gen.NewGenerator(gen.Config{
			ModelPkgPath: "./model",
		})
		g.UseDB(db)
		g.ApplyBasic(g.GenerateModel(tableName))
		g.Execute()
	}
}
