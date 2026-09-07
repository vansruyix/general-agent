package sse

import (
	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/schema"
)

// HandleMessageVariant 是处理 ADK agent 事件中 TypedMessageVariant 的顶层入口。
// 自动分发流式/非流式模式，以及 schema.Message / schema.AgenticMessage 类型。
//
// 返回完整聚合后的消息，供多轮对话缓存使用。
func HandleMessageVariant[M adk.MessageType](mv *adk.TypedMessageVariant[M], consumer Consumer) (M, error) {
	if mv == nil {
		var zero M
		return zero, nil
	}

	// 根据是否流式分发到不同处理分支
	if !mv.IsStreaming {
		return handleNonStreaming(mv, consumer)
	}
	return handleStreaming(mv, consumer)
}

// handleNonStreaming 处理非流式消息，根据具体类型分发到 ParseMessage 或 ParseAgenticMessage
func handleNonStreaming[M adk.MessageType](mv *adk.TypedMessageVariant[M], consumer Consumer) (M, error) {
	switch msg := any(mv.Message).(type) {
	case *schema.Message:
		if err := ParseMessage(msg, consumer); err != nil {
			var zero M
			return zero, err
		}
		return any(msg).(M), nil
	case *schema.AgenticMessage:
		if err := ParseAgenticMessage(msg, consumer); err != nil {
			var zero M
			return zero, err
		}
		return any(msg).(M), nil
	default:
		var zero M
		return zero, nil
	}
}

// handleStreaming 处理流式消息，根据 stream 的具体类型分发到对应的流式解析函数
func handleStreaming[M adk.MessageType](mv *adk.TypedMessageVariant[M], consumer Consumer) (M, error) {
	switch stream := any(mv.MessageStream).(type) {
	case *schema.StreamReader[*schema.Message]:
		result, err := ParseMessageStream(stream, consumer)
		if err != nil {
			var zero M
			return zero, err
		}
		return any(result).(M), nil
	case *schema.StreamReader[*schema.AgenticMessage]:
		result, err := ParseAgenticMessageStream(stream, consumer)
		if err != nil {
			var zero M
			return zero, err
		}
		return any(result).(M), nil
	default:
		var zero M
		return zero, nil
	}
}