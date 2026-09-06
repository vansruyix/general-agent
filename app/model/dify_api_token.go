package model

import "general-agent/extension/xtime"

type DifyApiToken struct {
	ID         string          `json:"id,omitempty" gorm:"primarykey"`
	CreatedAt  xtime.JsonTime  `json:"created_at,omitempty" gorm:"->;<-:create" swaggerignore:"true"` //创建时间
	TenantId   string          `json:"tenant_id"`                                                     // 租户ID
	AppId      string          `json:"app_id"`                                                        // 应用ID
	Type       string          `json:"type"`                                                          // 类型 app knowledge
	Token      string          `json:"token"`                                                         // token
	CopyToken  string          `json:"copy_token" gorm:"-"`                                           // 复制Token（加密）
	LastUsedAt *xtime.JsonTime `json:"last_used_at"`                                                  // 最后使用时间
}

func (DifyApiToken) TableName() string {
	return "api_tokens"
}
