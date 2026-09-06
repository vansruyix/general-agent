package user

type User struct {
	ID                             string `gorm:"primaryKey;column:id" json:"-"`
	Username                       string `gorm:"column:username" json:"username"`         // 用户名
	Password                       string `gorm:"column:password" json:"password"`         // 密码，MD5加密
	PasswordSm3                    string `gorm:"column:password_sm3" json:"password_sm3"` // 密码，SM3加密
	Name                           string `gorm:"column:name" json:"name"`
	CreateTime                     int64  `gorm:"column:create_time" json:"create_time"`                                               // 创建时间
	LatestLoginTime                int64  `gorm:"column:latest_login_time" json:"latest_login_time"`                                   // 最近一次登录时间
	PasswordChangeTime             int64  `gorm:"column:password_change_time" json:"password_change_time"`                             // 修改密码时间
	Status                         int8   `gorm:"column:status" json:"status"`                                                         // 0:正常|1:第一次登录|2:密码过期
	RoleID                         int    `gorm:"column:role_id" json:"role_id"`                                                       // 角色id
	Classification                 int8   `gorm:"column:classification" json:"classification"`                                         // 用户分类 0：系统内置用户|1：自定义用户
	GatherType                     int    `gorm:"column:gather_type" json:"gather_type"`                                               // 告警视角，默认按照事件名称聚合
	PlatformSha256Pass             string `gorm:"column:platform_sha256_pass" json:"platform_sha256_pass"`                             // 平台共用的sha256加密后的密码
	IsLock                         int8   `gorm:"column:is_lock" json:"is_lock"`                                                       // 当前用户是否被锁定
	OtpSn                          string `gorm:"column:otp_sn" json:"otp_sn"`                                                         // otp令牌的序列号
	MdsPass                        string `gorm:"column:mds_pass" json:"mds_pass"`                                                     // MDS平台密码
	CloudPass                      string `gorm:"column:cloud_pass" json:"cloud_pass"`                                                 // 云平台规范格式的密码
	RadiusServerID                 int    `gorm:"column:radius_server_id" json:"radius_server_id"`                                     // radius认证服务器ID
	AllowOperationModuleIDs        string `gorm:"column:allow_operation_module_ids" json:"allow_operation_module_ids"`                 // 允许操作的功能模块id，多个逗号分割
	AllowOperationModuleExpireTime int64  `gorm:"column:allow_operation_module_expire_time" json:"allow_operation_module_expire_time"` // 允许操作的功能模块限期过期时间，为0则代表永不过期
}

// TableName get sql table name.获取数据库表名
func (m *User) TableName() string {
	return "user"
}
