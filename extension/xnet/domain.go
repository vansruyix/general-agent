package xnet

import (
	"general-agent/extension/errorx"
	"net/url"
)

func IsValidDomain(domain string) bool {
	u, err := url.ParseRequestURI("http://" + domain)
	if err != nil {
		return false
	}

	if u.Host == "" {
		return false
	}

	return true
}
func IsHttpProtocolAssert(protocol string) {
	if !(protocol == "http" || protocol == "https") {
		panic(errorx.ErrRequestParam.WithMessage("Protocol cannot be empty"))
	}
}
