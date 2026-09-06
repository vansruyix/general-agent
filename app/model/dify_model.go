package model

import "general-agent/app/model/xconst"

// DifyTenantDefaultModel dify租户默认模型
type DifyTenantDefaultModel struct {
	BaseModel
	TenantID     string           `json:"tenant_id"`     // 租户ID
	ProviderName string           `json:"provider_name"` // 供应商标识
	ModelName    string           `json:"model_name"`    // 模型名称
	ModelType    xconst.ModelType `json:"model_type"`    // 模型类型
}

func (DifyTenantDefaultModel) TableName() string {
	return "tenant_default_models"
}

type DifyTenantDefaultModelVO struct {
	DifyTenantDefaultModel
	Enabled bool `json:"enabled"` // 是否启用
}
