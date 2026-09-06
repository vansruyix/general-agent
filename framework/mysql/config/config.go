// Package model
// @author: liang_jinxin
// @date: 2024/6/27
// @note:
package config

type Options struct {
	Primary    string       `mapstructure:"primary"`
	Datasource DynamicDbMap `mapstructure:"datasource"`
}

type DynamicDbMap map[string]MySQL

// MySQL 定义数据库连接与连接池参数。
type MySQL struct {
	Host            string `mapstructure:"host"`     // 地址IP
	Port            int    `mapstructure:"port"`     // 数据库端口
	Username        string `mapstructure:"username"` // 账号
	Password        string `mapstructure:"password"` // 密码
	Database        string `mapstructure:"database"` // 数据库名称
	MaxIdleConns    int    `mapstructure:"maxIdleConns"`
	MaxOpenConns    int    `mapstructure:"maxOpenConns"`
	ConnMaxLifetime int    `mapstructure:"connMaxLifetime"`
}
