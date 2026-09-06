package model

// User 用户表
type User struct {
	BaseModel
	TenantId       string `json:"tenant_id"`        // 租户ID
	Username       string `json:"username"`         // 用户名
	ManagerUserId  string `json:"manager_user_id"`  // 管理侧用户ID
	PlatformUserId string `json:"platform_user_id"` // 平台侧用户ID
	PlatformName   string `json:"platform_name"`    // 平台名称
	Email          string `json:"email"`            // 平台邮箱
	Password       string `json:"password"`         // 平台密码
	Role           string `json:"role"`             // 平台角色 owner|admin|editor|normal|dataset_operator
	Status         string `json:"status"`           // 平台状态
}

func (User) TableName() string {
	return "acc_user"
}
