package model

import (
	"general-agent/extension/xtime"
)

// RAGflowSnapshot RAGflow知识库
type RAGflowSnapshot struct {
	ID         string         `gorm:"primaryKey;column:id"`
	Name       string         `gorm:"column:name"`
	CreateTime int64          `gorm:"column:create_time"`
	CreateDate xtime.JsonTime `gorm:"column:create_date"`
	TenantId   string         `gorm:"column:tenant_id"`
	Status     int            `gorm:"column:status;comment:0失败 1成功 2进行中"`
	Log        string         `gorm:"column:log"`
}

func (RAGflowSnapshot) TableName() string {
	return "snapshot"
}

// 状态常量定义
const (
	StatusFailed   = 0 // 失败
	StatusSuccess  = 1 // 成功
	StatusProgress = 2 // 进行中
)
