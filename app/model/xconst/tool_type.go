package xconst

// ToolType 工具类型枚举
type ToolType string

const (
	BuiltinTool  ToolType = "builtin"  // 内置工具
	ModelTool    ToolType = "model"    // 模型工具
	APITool      ToolType = "api"      // API工具
	WorkflowTool ToolType = "workflow" // 工作流工具
)

// GetToolType 根据工具类型字符串获取工具类型枚举值
func GetToolType(toolType string) *ToolType {
	switch toolType {
	case "builtin":
		t := BuiltinTool
		return &t
	case "model":
		t := ModelTool
		return &t
	case "api":
		t := APITool
		return &t
	case "workflow":
		t := WorkflowTool
		return &t
	default:
		return nil
	}
}
