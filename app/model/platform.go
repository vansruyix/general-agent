package model

import "general-agent/app/model/xconst"

type Platform struct {
	BaseModel
	Name    string            `json:"name"`     // 平台名称 dify coze
	BaseUrl string            `json:"base_url"` // 平台地址
	Enabled xconst.EnableEnum `json:"enabled"`  // 平台启用状态 1-启用 2-禁用
}

func (Platform) TableName() string {
	return "acc_platform"
}
