package config

// Log 定义日志级别与输出格式（console / json）。
type Log struct {
	Level    string `mapstructure:"level"`
	Encoding string `mapstructure:"encoding"`
}
