// Package config 提供基于 viper 的配置加载能力。
// 通过 -env 命令行参数（默认 dev）指定运行环境，先加载 config/config.yaml 公共配置，
// 再加载 config/{env}.yaml 覆盖环境差异。
package config

import (
	"flag"
	"fmt"

	"github.com/spf13/viper"
	"go.uber.org/fx"
)

// env 运行环境标识，由 -env 命令行参数指定，默认值 "dev"。
var env = flag.String("env", "", "运行环境 (dev|test|prod)")

// Config 是应用全局配置的根结构体，聚合 HTTP、MySQL、Log、App 四个子配置。
type Config struct {
	HTTP  HTTP  `mapstructure:"http"`
	MySQL MySQL `mapstructure:"mysql"`
	Log   Log   `mapstructure:"log"`
	App   App   `mapstructure:"app"`
}

// HTTP 定义 HTTP 服务监听地址与 API 基础路径。
type HTTP struct {
	Addr     string `mapstructure:"addr"`
	BasePath string `mapstructure:"basePath"`
}

// MySQL 定义数据库连接与连接池参数。
type MySQL struct {
	Host            string `mapstructure:"host"`
	Port            int    `mapstructure:"port"`
	Username        string `mapstructure:"username"`
	Password        string `mapstructure:"password"`
	Database        string `mapstructure:"database"`
	MaxIdleConns    int    `mapstructure:"maxIdleConns"`
	MaxOpenConns    int    `mapstructure:"maxOpenConns"`
	ConnMaxLifetime int    `mapstructure:"connMaxLifetime"`
}

// Log 定义日志级别与输出格式（console / json）。
type Log struct {
	Level    string `mapstructure:"level"`
	Encoding string `mapstructure:"encoding"`
}

// App 定义应用元信息，Env 字段标记当前运行环境（dev / test / prod）。
type App struct {
	Env string `mapstructure:"env"`
}

// New 加载并返回应用配置。
// 加载顺序：
//  1. 读 config/config.yaml 并 unmarshal 到 Config
//  2. 确定运行环境：优先取 -env 命令行参数，未指定则取 config.yaml 中的 app.env
//  3. 读 config/{env}.yaml 合并覆盖
//  4. 重新 unmarshal 使环境配置生效
//
// 文件缺失或反序列化失败会直接返回 error，启动即失败，不静默降级。
func New() (*Config, error) {
	flag.Parse()

	v := viper.New()

	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath("../config")
	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("read config.yaml: %w", err)
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}

	// 确定运行环境：命令行参数优先，否则使用 config.yaml 中的 app.env
	activeEnv := *env
	if activeEnv == "" {
		activeEnv = cfg.App.Env
	}

	v.SetConfigName(activeEnv)
	if err := v.MergeInConfig(); err != nil {
		return nil, fmt.Errorf("merge %s.yaml: %w", activeEnv, err)
	}

	// 重新 unmarshal 使环境覆盖配置生效
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unmarshal merged config: %w", err)
	}
	return &cfg, nil
}

// Module 是 fx 模块，提供 *Config 单例。
var Module = fx.Module("config", fx.Provide(New))
