package model

// Provider 模型供应商
type Provider struct {
	BaseModel
	TenantId            string   `json:"tenant_id"`             // 租户ID
	PlatformTenantId    string   `json:"platform_tenant_id"`    // 平台租户ID
	Name                string   `json:"name"`                  // 供应商名称
	Provider            string   `json:"provider"`              // 供应商标识
	SupportedModelTypes []string `json:"supported_model_types"` // 支持的模型类型
	Status              string   `json:"status"`                // 状态
}

var OpenaiApiCompatible = Provider{
	Name:                "OpenAI-API-compatible",
	Provider:            "langgenius/openai_api_compatible/openai_api_compatible",
	SupportedModelTypes: []string{"llm", "text-embedding", "rerank", "speech2text", "tts"},
	Status:              "active",
}

var TongYi = Provider{
	Name:                "通义千问",
	Provider:            "langgenius/tongyi/tongyi",
	SupportedModelTypes: []string{"llm", "text-embedding", "rerank", "speech2text", "tts"},
	Status:              "active",
}

var XorbitsInference = Provider{
	Name:                "Xorbits Inference",
	Provider:            "langgenius/xinference/xinference",
	SupportedModelTypes: []string{"llm", "text-embedding", "rerank", "speech2text", "tts"},
	Status:              "active",
}

// GetProvider 获取模型供应商
func GetProvider(provider string) *Provider {
	switch provider {
	case OpenaiApiCompatible.Provider:
		return &OpenaiApiCompatible
	case TongYi.Provider:
		return &TongYi
	case XorbitsInference.Provider:
		return &XorbitsInference
	default:
		return nil
	}
}

type DifyProviderModels struct {
	BaseModel
	TenantId        string `json:"tenant_id"`        // 租户ID
	ProviderName    string `json:"provider_name"`    // 供应商名称
	ModelName       string `json:"model_name"`       // 模型名称
	ModelType       string `json:"model_type"`       // 模型类型
	EncryptedConfig string `json:"encrypted_config"` // 加密配置
	IsValid         bool   `json:"is_valid"`         // 是否有效
}

func (DifyProviderModels) TableName() string {
	return "provider_models"
}
