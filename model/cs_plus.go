package model

// OpLog [...]
type OpLog struct {
	ID        int    `gorm:"primaryKey;column:id" json:"-"`
	Username  string `gorm:"column:username" json:"username"`     // 用户名
	ClientIP  string `gorm:"column:client_ip" json:"client_ip"`   // 操作IP
	Module    string `gorm:"column:module" json:"module"`         // 操作模块
	Topic     string `gorm:"column:topic" json:"topic"`           // 操作来源
	ReqLog    string `gorm:"column:req_log" json:"req_log"`       // 请求记录
	Detail    string `gorm:"column:detail" json:"detail"`         // 操作详情
	Success   int8   `gorm:"column:success" json:"success"`       // 是否成功
	TimeStamp int64  `gorm:"column:time_stamp" json:"time_stamp"` // 时间戳
}

// TableName get sql table name.获取数据库表名
func (m *OpLog) TableName() string {
	return "op_log"
}
