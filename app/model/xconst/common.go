package xconst

type EnableEnum int8 // 启用

const (
	Enable  EnableEnum = 1 // 启用
	Disable EnableEnum = 2 // 停用、禁用
)

type YesEnum int8

const (
	YES YesEnum = 1 //是
	NO  YesEnum = 2 //否
)

type EnableStrEnum string

const (
	ON  EnableStrEnum = "on"  //开启
	OFF EnableStrEnum = "off" //关闭
)
