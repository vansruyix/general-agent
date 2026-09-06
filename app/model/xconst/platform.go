package xconst

// PlatformType 平台类型枚举
type PlatformType string

const (
	PlatformDify    PlatformType = "dify"
	PlatformCoze    PlatformType = "coze"
	PlatformRAGflow PlatformType = "ragflow"
)

type Endpoint struct {
	Method   string // 请求方法
	Endpoint string // 请求地址
}

func NewEndpoint(method, endpoint string) *Endpoint {
	return &Endpoint{
		Method:   method,
		Endpoint: endpoint,
	}
}
