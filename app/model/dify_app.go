package model

import "general-agent/app/model/xconst"

// DifyApp 应用
type DifyApp struct {
	BaseModel
	PlatformTenantId string          `json:"platform_tenant_id" gorm:"column:tenant_id"` // 平台租户ID
	Name             string          `json:"name"`                                       // 应用名称
	Mode             xconst.ModeType `json:"mode"`                                       // 应用模式
	Status           string          `json:"status"`                                     // 应用状态
	Description      string          `json:"description"`                                // 应用描述
	WorkflowId       string          `json:"workflow_id"`                                // 应用工作流ID
}

func (DifyApp) TableName() string {
	return "apps"
}

type DifyAppVO struct {
	DifyApp
	Tags []AppTag `json:"tags"` // 应用标签
}

type AppTag struct {
	AppId string `json:"-"`    // 应用ID
	Id    string `json:"id"`   // 标签ID
	Name  string `json:"name"` // 标签名称
	Type  string `json:"type"` // 标签类型
}
