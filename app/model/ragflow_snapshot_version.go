package model

import (
	"general-agent/extension/xtime"
)

// RAGflowSnapshotVersion 快照版本表
type RAGflowSnapshotVersion struct {
	ID         string         `gorm:"primaryKey;column:id"`                     // 主键ID
	SnapID     string         `gorm:"column:snap_id;index:idx_snap_id"`         // 快照ID，关联snapshot表
	TenantId   string         `gorm:"column:tenant_id;index:idx_tenant_id"`     // 租户ID
	Status     int            `gorm:"column:status;comment:0失败 1成功 2进行中"`       // 状态
	Version    string         `gorm:"column:version;comment:快照版本号"`             // 版本号
	CreateTime int64          `gorm:"column:create_time"`                       // 创建时间戳
	CreateDate xtime.JsonTime `gorm:"column:create_date;index:idx_create_date"` // 创建日期
}

func (RAGflowSnapshotVersion) TableName() string {
	return "snapshot_version"
}
