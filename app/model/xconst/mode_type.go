package xconst

// ModeType APP模式枚举
type ModeType string

const (
	ModeTypeCompletion   ModeType = "completion"    // 文本生成应用
	ModeTypeChat         ModeType = "chat"          // 聊天助手
	ModeTypeAdvancedChat ModeType = "advanced-chat" // 对话流
	ModeTypeWorkflow     ModeType = "workflow"      // 工作流
	ModeTypeAgentChat    ModeType = "agent-chat"    // Agent
	ModeTypeChannel      ModeType = "channel"
)

// GetModeType 获取APP模式枚举
func GetModeType(modeType string) *ModeType {
	switch modeType {
	case "completion":
		m := ModeTypeCompletion
		return &m
	case "chat":
		m := ModeTypeChat
		return &m
	case "advanced-chat":
		m := ModeTypeAdvancedChat
		return &m
	case "workflow":
		m := ModeTypeWorkflow
		return &m
	case "agent-chat":
		m := ModeTypeAgentChat
		return &m
	case "channel":
		m := ModeTypeChannel
		return &m
	default:
		return nil
	}
}
