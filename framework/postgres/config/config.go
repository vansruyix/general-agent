// Package model
// @author: liang_jinxin
// @date: 2024/6/27
// @note:
package config

type Options struct {
	Primary    string       `mapstructure:"primary"`
	Datasource DynamicDbMap `mapstructure:"datasource"`
}

type DynamicDbMap map[string]Postgres

// DBConfig 数据库配置
type Postgres struct {
	Database string `mapstructure:"database"` // 数据库名称
	Host     string `mapstructure:"host"`     // 地址IP
	Port     int    `mapstructure:"port"`     // 数据库端口
	Username string `mapstructure:"username"` // 账号
	Password string `mapstructure:"password"` // 密码
}
