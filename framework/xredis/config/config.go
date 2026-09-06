package config

// 自动绑定 REDIS_XX 的环境变量
type Redis struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
	Mode     string `mapstructure:"mode"  default:"single"` // single/sentinel
	Master   string `mapstructure:"master" default:"mymaster"`
	Nodes    string `mapstructure:"nodes" default:"" env:"REDIS_NODES"`
	PoolSize int    `mapstructure:"pool-size" default:"10"`
}
