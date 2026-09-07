package sse

import (
	"fmt"
	"io"

	"github.com/cloudwego/eino/schema"
)

// ParseMessage 解析非流式的 schema.Message，通过 consumer 消费每个 SSEvent。
// 仅处理 Assistant 和 Tool 角色的消息，User/System 消息不产生事件。
func ParseMessage(msg *schema.Message, consumer Consumer) error {
	if msg == nil {
		return nil
	}

	switch msg.Role {
	case schema.Assistant:
		return parseAssistantMessage(msg, consumer)
	case schema.Tool:
		return parseToolMessage(msg, consumer)
	default:
		// User、System 消息不产生 SSE 事件
		return nil
	}
}

// parseAssistantMessage 解析 Assistant 角色的消息，按顺序产生：
// 思考过程 → 多模态输出 → 文本内容 → 工具调用
func parseAssistantMessage(msg *schema.Message, consumer Consumer) error {
	// 1. 思考过程（推理模型的 chain-of-thought）
	if msg.ReasoningContent != "" {
		if err := consumer(SSEvent{
			Type: EventThinking,
			Data: msg.ReasoningContent,
		}); err != nil {
			return err
		}
	}

	// 2. 多模态输出（AssistantGenMultiContent），有此字段时不再处理 Content
	if len(msg.AssistantGenMultiContent) > 0 {
		for _, part := range msg.AssistantGenMultiContent {
			event, err := outputPartToEvent(part)
			if err != nil {
				return err
			}
			if err := consumer(event); err != nil {
				return err
			}
		}
		return nil
	}

	// 3. 文本内容
	if msg.Content != "" {
		if err := consumer(SSEvent{
			Type: EventText,
			Data: msg.Content,
		}); err != nil {
			return err
		}
	}

	// 4. 工具调用
	for _, tc := range msg.ToolCalls {
		if err := consumer(SSEvent{
			Type: EventToolCall,
			Data: tc.Function.Arguments,
			Meta: map[string]any{
				"tool_name": tc.Function.Name,
				"call_id":   tc.ID,
			},
		}); err != nil {
			return err
		}
	}

	return nil
}

// parseToolMessage 解析 Tool 角色的消息，产生工具调用结果事件
func parseToolMessage(msg *schema.Message, consumer Consumer) error {
	return consumer(SSEvent{
		Type: EventToolResult,
		Data: msg.Content,
		Meta: map[string]any{
			"call_id":   msg.ToolCallID,
			"tool_name": msg.ToolName,
		},
	})
}

// outputPartToEvent 将单个 MessageOutputPart 转换为对应的 SSEvent
func outputPartToEvent(part schema.MessageOutputPart) (SSEvent, error) {
	switch part.Type {
	case schema.ChatMessagePartTypeText:
		return SSEvent{Type: EventText, Data: part.Text}, nil
	case schema.ChatMessagePartTypeImageURL:
		event := SSEvent{Type: EventImage}
		if part.Image != nil {
			if part.Image.URL != nil {
				event.URL = *part.Image.URL
			}
			if part.Image.Base64Data != nil {
				event.Base64 = *part.Image.Base64Data
			}
			event.MIMEType = part.Image.MIMEType
		}
		return event, nil
	case schema.ChatMessagePartTypeAudioURL:
		event := SSEvent{Type: EventAudio}
		if part.Audio != nil {
			if part.Audio.URL != nil {
				event.URL = *part.Audio.URL
			}
			if part.Audio.Base64Data != nil {
				event.Base64 = *part.Audio.Base64Data
			}
			event.MIMEType = part.Audio.MIMEType
		}
		return event, nil
	case schema.ChatMessagePartTypeVideoURL:
		event := SSEvent{Type: EventVideo}
		if part.Video != nil {
			if part.Video.URL != nil {
				event.URL = *part.Video.URL
			}
			if part.Video.Base64Data != nil {
				event.Base64 = *part.Video.Base64Data
			}
			event.MIMEType = part.Video.MIMEType
		}
		return event, nil
	case schema.ChatMessagePartTypeReasoning:
		event := SSEvent{Type: EventThinking}
		if part.Reasoning != nil {
			event.Data = part.Reasoning.Text
		}
		return event, nil
	default:
		return SSEvent{}, fmt.Errorf("unknown output part type: %s", part.Type)
	}
}

// ParseAgenticMessage 解析非流式的 schema.AgenticMessage，通过 consumer 消费每个 SSEvent。
func ParseAgenticMessage(msg *schema.AgenticMessage, consumer Consumer) error {
	if msg == nil {
		return nil
	}

	for _, block := range msg.ContentBlocks {
		if block == nil {
			continue
		}
		event, err := contentBlockToEvent(block)
		if err != nil {
			return err
		}
		if err := consumer(event); err != nil {
			return err
		}
	}

	return nil
}

// contentBlockToEvent 将单个 ContentBlock 转换为对应的 SSEvent
func contentBlockToEvent(block *schema.ContentBlock) (SSEvent, error) {
	switch block.Type {
	case schema.ContentBlockTypeReasoning:
		event := SSEvent{Type: EventThinking}
		if block.Reasoning != nil {
			event.Data = block.Reasoning.Text
		}
		return event, nil
	case schema.ContentBlockTypeAssistantGenText:
		event := SSEvent{Type: EventText}
		if block.AssistantGenText != nil {
			event.Data = block.AssistantGenText.Text
		}
		return event, nil
	case schema.ContentBlockTypeAssistantGenImage:
		event := SSEvent{Type: EventImage}
		if block.AssistantGenImage != nil {
			event.URL = block.AssistantGenImage.URL
			event.Base64 = block.AssistantGenImage.Base64Data
			event.MIMEType = block.AssistantGenImage.MIMEType
		}
		return event, nil
	case schema.ContentBlockTypeAssistantGenAudio:
		event := SSEvent{Type: EventAudio}
		if block.AssistantGenAudio != nil {
			event.URL = block.AssistantGenAudio.URL
			event.Base64 = block.AssistantGenAudio.Base64Data
			event.MIMEType = block.AssistantGenAudio.MIMEType
		}
		return event, nil
	case schema.ContentBlockTypeAssistantGenVideo:
		event := SSEvent{Type: EventVideo}
		if block.AssistantGenVideo != nil {
			event.URL = block.AssistantGenVideo.URL
			event.Base64 = block.AssistantGenVideo.Base64Data
			event.MIMEType = block.AssistantGenVideo.MIMEType
		}
		return event, nil
	case schema.ContentBlockTypeFunctionToolCall:
		event := SSEvent{Type: EventToolCall}
		if block.FunctionToolCall != nil {
			event.Data = block.FunctionToolCall.Arguments
			event.Meta = map[string]any{
				"tool_name": block.FunctionToolCall.Name,
				"call_id":   block.FunctionToolCall.CallID,
			}
		}
		return event, nil
	case schema.ContentBlockTypeFunctionToolResult:
		event := SSEvent{Type: EventToolResult}
		if block.FunctionToolResult != nil {
			event.Meta = map[string]any{
				"tool_name": block.FunctionToolResult.Name,
				"call_id":   block.FunctionToolResult.CallID,
			}
		}
		return event, nil
	default:
		// 其他类型（UserInput*、Server/MCP tool calls 等）暂不产生事件
		return SSEvent{}, nil
	}
}

// ParseMessageStream 从流式 schema.Message reader 中逐块读取，
// 每块通过 consumer 消费为 SSEvent，最后返回完整聚合后的 schema.Message
// 供多轮对话缓存使用。
func ParseMessageStream(stream *schema.StreamReader[*schema.Message], consumer Consumer) (*schema.Message, error) {
	defer stream.Close()

	var msgs []*schema.Message
	for {
		msg, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		msgs = append(msgs, msg)
		// 每个 chunk 实时通过 consumer 消费
		if err := ParseMessage(msg, consumer); err != nil {
			return nil, err
		}
	}

	// 聚合所有 chunk 为完整消息，供多轮对话缓存
	return schema.ConcatMessages(msgs)
}

// ParseAgenticMessageStream 从流式 schema.AgenticMessage reader 中逐块读取，
// 每块通过 consumer 消费为 SSEvent，最后返回完整聚合后的 schema.AgenticMessage
// 供多轮对话缓存使用。
func ParseAgenticMessageStream(stream *schema.StreamReader[*schema.AgenticMessage], consumer Consumer) (*schema.AgenticMessage, error) {
	defer stream.Close()

	var msgs []*schema.AgenticMessage
	for {
		msg, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		msgs = append(msgs, msg)
		// 每个 chunk 实时通过 consumer 消费
		if err := ParseAgenticMessage(msg, consumer); err != nil {
			return nil, err
		}
	}

	// 聚合所有 chunk 为完整消息，供多轮对话缓存
	return schema.ConcatAgenticMessages(msgs)
}
