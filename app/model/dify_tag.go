package model

import "general-agent/extension/xtime"

type DifyTag struct {
	ID        string         `json:"id,omitempty" gorm:"primarykey"`
	TenantId  string         `json:"tenant_id"`                                                     // 租户ID
	Type      string         `json:"type"`                                                          // 标签类型 app knowledge
	Name      string         `json:"name"`                                                          // 标签名称
	CreatedBy string         `json:"created_by"`                                                    // 创建人ID
	CreatedAt xtime.JsonTime `json:"created_at,omitempty" gorm:"->;<-:create" swaggerignore:"true"` //创建时间
}

func (DifyTag) TableName() string {
	return "tags"
}

type DifyTagBind struct {
	ID        string         `json:"id,omitempty" gorm:"primarykey"`
	TenantId  string         `json:"tenant_id"`                                                     // 租户ID
	TagId     string         `json:"tag_id"`                                                        // 标签ID
	TargetId  string         `json:"target_id"`                                                     // 目标ID
	CreatedBy string         `json:"created_by"`                                                    // 创建人ID
	CreatedAt xtime.JsonTime `json:"created_at,omitempty" gorm:"->;<-:create" swaggerignore:"true"` //创建时间
}

func (DifyTagBind) TableName() string {
	return "tag_bindings"
}
