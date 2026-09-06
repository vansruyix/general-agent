package model

// DifyTool 工具API提供者
type DifyTool struct {
	BaseModel
	Name        string `json:"name"`        // 工具名称
	Label       string `json:"label"`       // 工具标签
	Type        string `json:"type"`        // 工具类型
	Description string `json:"description"` // 工具描述
}

func (DifyTool) TableName() string {
	return "tool_api_providers"
}

// DifyToolVO 工具API提供者VO
type DifyToolVO struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Label       I18n   `json:"label"`
	Description I18n   `json:"description"`
	Type        string `json:"type"`
}
