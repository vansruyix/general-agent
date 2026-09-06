package xconst

type MigrateType string

const (
	MigrateApp      MigrateType = "app"      // 导出应用
	MigrateTool     MigrateType = "tool"     // 导出工具
	MigrateProvider MigrateType = "provider" // 导出模型供应商
	MigrateTag      MigrateType = "tag"      // 导出标签
	MigrateTenant   MigrateType = "tenant"   // 导出租户
)

// MigrateDataFileSuffaix 迁移数据文件后缀
const MigrateDataFileSuffaix = ".dat"

// ExportZipFileName 导出压缩包文件名
const ExportZipFileName = "export-liangjie-"

// ExportZipPassword 压缩包密码
const ExportZipPassword = "T$6VEIZ5s$"

// MigrateTempDirPath 数据迁移临时目录
const MigrateTempDirPath = "/tmp/agentcc-migrate/"
