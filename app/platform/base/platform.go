package base

import (
	"general-agent/app/model"
	"general-agent/app/model/xconst"
	"time"
)

// Platform 平台接口定义
type Platform interface {
	Register(username, email, password string) error
	Login(tu *model.TenantUser, inviteToken *string) (*LoginRes, error)
	Invoke(tu *model.TenantUser, endpoint *xconst.Endpoint, req interface{}, resPtr interface{}) error
	InvokeWithVars(tu *model.TenantUser, endpoint *xconst.Endpoint, pathVars map[string]string, req interface{}, resPtr interface{}) error
}

// LoginRes 登录返回结果
type LoginRes struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// Config 平台配置
type Config struct {
	Type     xconst.PlatformType `json:"type"`
	BaseUrl  string              `json:"base_url"` // 平台地址
	Timeout  time.Duration       `json:"timeout"`
	Settings map[string]string   `json:"settings,omitempty"`
}

// Error 平台错误
type Error struct {
	Platform xconst.PlatformType `json:"platform"`
	Code     string              `json:"code"`
	Message  string              `json:"message"`
}

func (e *Error) Error() string {
	return string(e.Platform) + ": " + e.Message
}
