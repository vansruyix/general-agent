package main

import (
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gen"
	"gorm.io/gorm"
)

// TODO 运行前替换为实际的 MySQL 连接信息
const dsn = "root:123456@tcp(192.168.67.168:3306)/cs_plus?charset=utf8mb4&parseTime=True&loc=Local"

func main() {
	generateModel()
}

func generateModel() {
	db, err := gorm.Open(mysql.Open(dsn))
	if err != nil {
		log.Fatalf("connect mysql failed: %v", err)
	}

	g := gen.NewGenerator(gen.Config{
		ModelPkgPath: "./app/entity",
		// FieldNullable: true, // NULL 列生成指针类型
		// MergeQuery:    true,
	})
	g.UseDB(db)
	g.ApplyBasic(g.GenerateModel("user"))
	g.Execute()
}
