package sse

type MetaKey string

const (
	CallId    MetaKey = "call_id"    //工具id
	ToolName  MetaKey = "tool_name"  //工具名称
	ToolParam MetaKey = "tool_param" //工具调用参数
)

// EventType 表示 SSE 事件的类型
type EventType string

const (
	EventThinking   EventType = "thinking"    // EventThinking 思考过程（推理模型的 chain-of-thought）
	EventText       EventType = "text"        // EventText 文本输出
	EventToolCall   EventType = "tool_call"   // EventToolCall 工具调用请求
	EventToolResult EventType = "tool_result" // EventToolResult 工具调用结果
	EventImage      EventType = "image"       // EventImage 图片输出
	EventAudio      EventType = "audio"       // EventAudio 音频输出
	EventVideo      EventType = "video"       // EventVideo 视频输出
	EventFile       EventType = "file"        // EventFile 文件输出
	EventError      EventType = "error"       // EventError 错误信息
	EventDone       EventType = "done"        // EventDone 流结束标记，可附带 token 使用统计等元数据
)

// SSEvent 是统一的 SSE 流式响应结构体，兼容思考过程、文本输出、
// 工具调用、图片、音频、视频等多模态输出。
type SSEvent struct {
	Type     EventType       `json:"type"`                // Type 事件类型
	Data     string          `json:"data,omitempty"`      // Data 文本内容或 JSON 序列化后的数据
	MIMEType string          `json:"mime_type,omitempty"` // MIMEType 多模态内容的 MIME 类型，例如 "image/png"
	Base64   string          `json:"base64,omitempty"`    // Base64 多模态内容的 Base64 编码数据
	URL      string          `json:"url,omitempty"`       // URL 多模态内容的 URL 地址
	Meta     map[MetaKey]any `json:"meta,omitempty"`      // Meta 附加元数据，如 tool_name、call_id、token_usage 等
}
