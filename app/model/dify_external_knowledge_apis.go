package model

import (
	"time"
)

// DifyExternalKnowledgeApi 对应public.external_knowledge_apis表的真实结构
type DifyExternalKnowledgeApi struct {
	Id          string    `gorm:"column:id;type:uuid;primaryKey"`                                         // 主键，uuid类型
	Name        string    `gorm:"column:name;type:varchar(255);not null;index"`                           // 名称，varchar(255)，非空，有索引
	Description string    `gorm:"column:description;type:varchar(255);not null"`                          // 描述，varchar(255)，非空
	TenantId    string    `gorm:"column:tenant_id;type:uuid;not null;index"`                              // 租户ID，uuid类型，非空，有索引
	Settings    string    `gorm:"column:settings;type:text;null"`                                         // 配置，text类型，可为空
	CreatedBy   string    `gorm:"column:created_by;type:uuid;not null"`                                   // 创建人，uuid类型，非空
	CreatedAt   time.Time `gorm:"column:created_at;type:timestamp;not null;default:CURRENT_TIMESTAMP(0)"` // 创建时间
	UpdatedBy   string    `gorm:"column:updated_by;type:uuid;null"`                                       // 更新人，uuid类型，可为空
	UpdatedAt   time.Time `gorm:"column:updated_at;type:timestamp;not null;default:CURRENT_TIMESTAMP(0)"` // 更新时间
}

// TableName 显式指定数据库表名，避免GORM自动复数化
func (DifyExternalKnowledgeApi) TableName() string {
	return "external_knowledge_apis"
}
