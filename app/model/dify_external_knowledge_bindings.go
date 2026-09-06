package model

import (
	"time"
)

// DifyExternalKnowledgeBinding 对应public.external_knowledge_bindings表结构
type DifyExternalKnowledgeBinding struct {
	Id                     string    `gorm:"column:id;type:uuid;primaryKey"`
	TenantId               string    `gorm:"column:tenant_id;type:uuid;index"`
	ExternalKnowledgeApiId string    `gorm:"column:external_knowledge_api_id;type:uuid;index"`
	DatasetId              string    `gorm:"column:dataset_id;type:uuid;index"`
	ExternalKnowledgeId    string    `gorm:"column:external_knowledge_id;type:text;index"` // 已发布的知识库ID
	CreatedBy              string    `gorm:"column:created_by;type:uuid"`
	CreatedAt              time.Time `gorm:"column:created_at;type:timestamp;default:CURRENT_TIMESTAMP(0)"`
	UpdatedBy              string    `gorm:"column:updated_by;type:uuid;null"`
	UpdatedAt              time.Time `gorm:"column:updated_at;type:timestamp;default:CURRENT_TIMESTAMP(0)"`
}

// TableName 指定数据库表名
func (DifyExternalKnowledgeBinding) TableName() string {
	return "external_knowledge_bindings"
}
