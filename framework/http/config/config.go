package config

import "time"

type Options struct {
	Host         string `mapstructure:"host"`
	Port         int    `mapstructure:"port"`
	InternalPort string `mapstructure:"internal-port"`
	APIPrefix    string `mapstructure:"api-prefix" default:"/api"`
	Token        string `mapstructure:"token" default:""`
	// ReverseProxy is a flag to indicate whether this is a reverse proxy, if you
	// use a reverse proxy like nginx, traefik, etc. to proxy the request to
	// this server, you should set this flag to true, otherwise, the real ip of
	// the client may not be able to get.
	ReverseProxy bool `mapstructure:"reverse_proxy" default:"false"`

	ReadTimeout  time.Duration `mapstructure:"read_timeout" json:"read_timeout" default:"30000000000"`   // 30s
	WriteTimeout time.Duration `mapstructure:"write_timeout" json:"write_timeout" default:"30000000000"` // 30s

	LoginWhiteList      []string `mapstructure:"login-white-list"`
	PermissionWhiteList []string `mapstructure:"permission-white-list"`

	PostQueryList []string `mapstructure:"post-query-list"`
}
