package model

import "general-agent/extension/xtime"

type DifyConversation struct {
	BaseModel
	AppId            string          `json:"app_id"`              // 应用ID
	AppModelConfigId string          `json:"app_model_config_id"` // 应用模型配置ID
	ModelProvider    string          `json:"model_provider"`      // 模型提供者
	ModelId          string          `json:"model_id"`            // 模型ID
	Mode             string          `json:"mode"`                // 模式
	Name             string          `json:"name"`                // 名称
	Summary          string          `json:"summary"`             // 摘要
	Status           string          `json:"status"`              // 状态
	FromSource       string          `json:"from_source"`         // 来源
	FromEndUserId    string          `json:"from_end_user_id"`    // 来源用户ID
	FromAccountId    string          `json:"from_account_id"`     // 来源账户ID
	ReadAt           *xtime.JsonTime `json:"read_at"`             // 阅读时间
	ReadAccountId    string          `json:"read_account_id"`     // 阅读账户ID
	InvokeFrom       string          `json:"invoke_from"`         // 调用来源
}

func (DifyConversation) TableName() string {
	return "conversations"
}
