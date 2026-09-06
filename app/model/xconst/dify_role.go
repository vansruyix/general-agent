package xconst

type DifyRole string

const (
	DifyDatasetOperator DifyRole = "dataset_operator" // 数据集操作员，能够查看和操作数据集，但不能修改数据集的权限
	DifyNormal          DifyRole = "normal"           // 成员，只能够使用应用程序，不能建立应用程序
	DifyEditor          DifyRole = "editor"           // 编辑，能够建立并编辑应用程序，不能管理团队设置
	DifyAdmin           DifyRole = "admin"            // 管理员，能够建立应用程序和管理团队设置
	DifyOwner           DifyRole = "owner"            // 所有者
)

// GetDifyRole 根据role字符串获取对应的DifyRole类型
func GetDifyRole(role string) DifyRole {
	switch role {
	case "dataset_operator":
		return DifyDatasetOperator
	case "normal":
		return DifyNormal
	case "editor":
		return DifyEditor
	case "admin":
		return DifyAdmin
	case "owner":
		return DifyOwner
	default:
		return DifyNormal
	}
}
