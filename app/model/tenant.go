package model

import "general-agent/app/model/xconst"

// Tenant 租户表
type Tenant struct {
	BaseModel
	Username         string `json:"username"`           // 租户名称
	ManagerTenantId  string `json:"manager_tenant_id"`  // 管理侧租户ID
	PlatformTenantId string `json:"platform_tenant_id"` // 平台租户ID
	PlatformName     string `json:"platform_name"`      // 平台名称
	LicenseCode      string `json:"license_code"`       // 平台授权码
	Status           string `json:"status"`             // 平台状态
}

func (Tenant) TableName() string {
	return "acc_tenant"
}

type TenantUser struct {
	TenantId         string `json:"tenant_id"`
	UserId           string `json:"user_id"`
	ManagerTenantId  string `json:"manager_tenant_id"`
	ManagerUserId    string `json:"manager_user_id"`
	PlatformTenantId string `json:"platform_tenant_id"`
	PlatformUserId   string `json:"platform_user_id"`
	PlatformName     string `json:"platform_name"`
	LicenseCode      string `json:"license_code"`
	Email            string `json:"email"`
	Username         string `json:"username"`
	Password         string `json:"password"`
	Role             string `json:"role"`
	InvokePlatform   xconst.PlatformType
}

// SuperTenantId 超级租户ID，可查看所有租户数据
const SuperTenantId = "10000"

// IsSuper 是否是超级管理员
func (tu *TenantUser) IsSuper() bool {
	return tu != nil && tu.ManagerTenantId == SuperTenantId
}

// IsOwner 是否为工作空间所有者
func (tu *TenantUser) IsOwner() bool {
	return tu != nil && tu.Role == string(xconst.DifyOwner)
}

// IsManager 是否具有管理权限
func (tu *TenantUser) IsManager() bool {
	role := xconst.GetDifyRole(tu.Role)
	return tu != nil && (role == xconst.DifyOwner || role == xconst.DifyAdmin)
}

// IsEditor 是否具有编辑权限
func (tu *TenantUser) IsEditor() bool {
	role := xconst.GetDifyRole(tu.Role)
	return tu != nil && (role == xconst.DifyOwner || role == xconst.DifyAdmin || role == xconst.DifyEditor)
}

type RAGflowTenantUser struct {
	AgentTenantUser *TenantUser `gorm:"-"`
	TenantId        string
	UserId          string
	Email           string
	Nickname        string
	Role            string
}
