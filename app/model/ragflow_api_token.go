package model

import "general-agent/extension/xtime"

// RAGflowApiToken rag_flow.api_token表模型（与RAGflowKnowledge风格对齐）
type RAGflowApiToken struct {
	CreateTime int64          `json:"create_time" gorm:"column:create_time"`        // 创建时间戳
	CreateDate xtime.JsonTime `json:"create_date" gorm:"column:create_date"`        // 创建时间
	UpdateTime int64          `json:"update_time" gorm:"column:update_time"`        // 更新时间戳
	UpdateDate xtime.JsonTime `json:"update_date" gorm:"column:update_date"`        // 更新时间
	TenantId   string         `json:"tenant_id" gorm:"column:tenant_id;primaryKey"` // 租户ID（复合主键）
	Token      string         `json:"token" gorm:"column:token;primaryKey"`         // 令牌（复合主键）
	DialogId   *string        `json:"dialog_id,omitempty" gorm:"column:dialog_id"`  // 对话ID（可为空）
	Source     *string        `json:"source,omitempty" gorm:"column:source"`        // 来源（可为空）
	Beta       *string        `json:"beta,omitempty" gorm:"column:beta"`            // 测试标识（可为空）
}

// TableName 指定数据库表名（与RAGflowKnowledge保持一致的命名规范）
func (RAGflowApiToken) TableName() string {
	return "api_token"
}
