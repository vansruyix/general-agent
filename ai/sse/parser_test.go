package sse

import (
	"testing"

	"github.com/cloudwego/eino/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseMessageStream(t *testing.T) {
	sr, sw := schema.Pipe[*schema.Message](3)
	sw.Send(&schema.Message{Role: schema.Assistant, Content: "Hello "}, nil)
	sw.Send(&schema.Message{Role: schema.Assistant, Content: "World"}, nil)
	sw.Close()

	var events []SSEvent
	consumer := func(event SSEvent) error {
		events = append(events, event)
		return nil
	}

	result, err := ParseMessageStream(sr, consumer)
	require.NoError(t, err)
	require.NotNil(t, result)
	// Each chunk is consumed individually
	require.Len(t, events, 2)
	assert.Equal(t, "Hello ", events[0].Data)
	assert.Equal(t, "World", events[1].Data)
	// The returned result is the concatenated message
	assert.Equal(t, "Hello World", result.Content)
	assert.Equal(t, schema.Assistant, result.Role)
}

func TestParseMessageStream_Reasoning(t *testing.T) {
	sr, sw := schema.Pipe[*schema.Message](3)
	sw.Send(&schema.Message{Role: schema.Assistant, ReasoningContent: "thinking..."}, nil)
	sw.Send(&schema.Message{Role: schema.Assistant, Content: "answer"}, nil)
	sw.Close()

	var events []SSEvent
	consumer := func(event SSEvent) error {
		events = append(events, event)
		return nil
	}

	result, err := ParseMessageStream(sr, consumer)
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Len(t, events, 2)
	assert.Equal(t, EventThinking, events[0].Type)
	assert.Equal(t, "thinking...", events[0].Data)
	assert.Equal(t, EventText, events[1].Type)
	assert.Equal(t, "answer", events[1].Data)
	assert.Equal(t, "thinking...", result.ReasoningContent)
	assert.Equal(t, "answer", result.Content)
}

func TestParseAgenticMessageStream(t *testing.T) {
	sr, sw := schema.Pipe[*schema.AgenticMessage](3)
	sw.Send(&schema.AgenticMessage{
		Role: schema.AgenticRoleTypeAssistant,
		ContentBlocks: []*schema.ContentBlock{
			schema.NewContentBlock(&schema.AssistantGenText{Text: "Hello "}),
		},
	}, nil)
	sw.Send(&schema.AgenticMessage{
		Role: schema.AgenticRoleTypeAssistant,
		ContentBlocks: []*schema.ContentBlock{
			schema.NewContentBlock(&schema.AssistantGenText{Text: "agentic world"}),
		},
	}, nil)
	sw.Close()

	var events []SSEvent
	consumer := func(event SSEvent) error {
		events = append(events, event)
		return nil
	}

	result, err := ParseAgenticMessageStream(sr, consumer)
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Len(t, events, 2)
	assert.Equal(t, "Hello ", events[0].Data)
	assert.Equal(t, "agentic world", events[1].Data)
	// The result should have the content blocks (non-streaming blocks are not merged)
	require.Len(t, result.ContentBlocks, 2)
	assert.Equal(t, "Hello ", result.ContentBlocks[0].AssistantGenText.Text)
	assert.Equal(t, "agentic world", result.ContentBlocks[1].AssistantGenText.Text)
}

func TestParseMessage_AssistantText(t *testing.T) {
	msg := &schema.Message{
		Role:    schema.Assistant,
		Content: "你好，有什么可以帮助你的？",
	}

	var events []SSEvent
	consumer := func(event SSEvent) error {
		events = append(events, event)
		return nil
	}

	err := ParseMessage(msg, consumer)
	require.NoError(t, err)
	require.Len(t, events, 1)
	assert.Equal(t, EventText, events[0].Type)
	assert.Equal(t, "你好，有什么可以帮助你的？", events[0].Data)
}

func TestParseMessage_ReasoningContent(t *testing.T) {
	msg := &schema.Message{
		Role:             schema.Assistant,
		ReasoningContent: "让我思考一下这个问题...",
		Content:          "最终答案是42。",
	}

	var events []SSEvent
	consumer := func(event SSEvent) error {
		events = append(events, event)
		return nil
	}

	err := ParseMessage(msg, consumer)
	require.NoError(t, err)
	require.Len(t, events, 2)
	assert.Equal(t, EventThinking, events[0].Type)
	assert.Equal(t, "让我思考一下这个问题...", events[0].Data)
	assert.Equal(t, EventText, events[1].Type)
	assert.Equal(t, "最终答案是42。", events[1].Data)
}

func TestParseMessage_ToolCalls(t *testing.T) {
	msg := &schema.Message{
		Role: schema.Assistant,
		ToolCalls: []schema.ToolCall{
			{
				ID:   "call_1",
				Type: "function",
				Function: schema.FunctionCall{
					Name:      "search",
					Arguments: `{"query": "golang"}`,
				},
			},
		},
	}

	var events []SSEvent
	consumer := func(event SSEvent) error {
		events = append(events, event)
		return nil
	}

	err := ParseMessage(msg, consumer)
	require.NoError(t, err)
	require.Len(t, events, 1)
	assert.Equal(t, EventToolCall, events[0].Type)
	assert.Equal(t, "search", events[0].Meta["tool_name"])
	assert.Equal(t, "call_1", events[0].Meta["call_id"])
}

func TestParseMessage_ToolResult(t *testing.T) {
	msg := &schema.Message{
		Role:       schema.Tool,
		Content:    `{"result": "found 10 items"}`,
		ToolCallID: "call_1",
		ToolName:   "search",
	}

	var events []SSEvent
	consumer := func(event SSEvent) error {
		events = append(events, event)
		return nil
	}

	err := ParseMessage(msg, consumer)
	require.NoError(t, err)
	require.Len(t, events, 1)
	assert.Equal(t, EventToolResult, events[0].Type)
	assert.Equal(t, `{"result": "found 10 items"}`, events[0].Data)
	assert.Equal(t, "call_1", events[0].Meta["call_id"])
	assert.Equal(t, "search", events[0].Meta["tool_name"])
}

func TestParseMessage_UserMessage(t *testing.T) {
	msg := &schema.Message{
		Role:    schema.User,
		Content: "用户的问题",
	}

	var events []SSEvent
	consumer := func(event SSEvent) error {
		events = append(events, event)
		return nil
	}

	err := ParseMessage(msg, consumer)
	require.NoError(t, err)
	// User and System messages should not generate events
	assert.Len(t, events, 0)
}

func TestParseMessage_MultiModalOutput(t *testing.T) {
	msg := &schema.Message{
		Role: schema.Assistant,
		AssistantGenMultiContent: []schema.MessageOutputPart{
			{
				Type: schema.ChatMessagePartTypeText,
				Text: "这是生成的图片：",
			},
			{
				Type: schema.ChatMessagePartTypeImageURL,
				Image: &schema.MessageOutputImage{
					MessagePartCommon: schema.MessagePartCommon{
						URL:      stringPtr("https://example.com/image.png"),
						MIMEType: "image/png",
					},
				},
			},
		},
	}

	var events []SSEvent
	consumer := func(event SSEvent) error {
		events = append(events, event)
		return nil
	}

	err := ParseMessage(msg, consumer)
	require.NoError(t, err)
	require.Len(t, events, 2)
	assert.Equal(t, EventText, events[0].Type)
	assert.Equal(t, "这是生成的图片：", events[0].Data)
	assert.Equal(t, EventImage, events[1].Type)
	assert.Equal(t, "https://example.com/image.png", events[1].URL)
	assert.Equal(t, "image/png", events[1].MIMEType)
}

func TestParseAgenticMessage_AssistantGenText(t *testing.T) {
	msg := &schema.AgenticMessage{
		Role: schema.AgenticRoleTypeAssistant,
		ContentBlocks: []*schema.ContentBlock{
			schema.NewContentBlock(&schema.AssistantGenText{Text: "Hello from agentic"}),
		},
	}

	var events []SSEvent
	consumer := func(event SSEvent) error {
		events = append(events, event)
		return nil
	}

	err := ParseAgenticMessage(msg, consumer)
	require.NoError(t, err)
	require.Len(t, events, 1)
	assert.Equal(t, EventText, events[0].Type)
	assert.Equal(t, "Hello from agentic", events[0].Data)
}

func TestParseAgenticMessage_Reasoning(t *testing.T) {
	msg := &schema.AgenticMessage{
		Role: schema.AgenticRoleTypeAssistant,
		ContentBlocks: []*schema.ContentBlock{
			schema.NewContentBlock(&schema.Reasoning{Text: "逐步推理..."}),
			schema.NewContentBlock(&schema.AssistantGenText{Text: "推理结果"}),
		},
	}

	var events []SSEvent
	consumer := func(event SSEvent) error {
		events = append(events, event)
		return nil
	}

	err := ParseAgenticMessage(msg, consumer)
	require.NoError(t, err)
	require.Len(t, events, 2)
	assert.Equal(t, EventThinking, events[0].Type)
	assert.Equal(t, "逐步推理...", events[0].Data)
	assert.Equal(t, EventText, events[1].Type)
	assert.Equal(t, "推理结果", events[1].Data)
}

func TestParseAgenticMessage_ToolCall(t *testing.T) {
	msg := &schema.AgenticMessage{
		Role: schema.AgenticRoleTypeAssistant,
		ContentBlocks: []*schema.ContentBlock{
			schema.NewContentBlock(&schema.FunctionToolCall{
				CallID:    "call_2",
				Name:      "get_weather",
				Arguments: `{"city": "Beijing"}`,
			}),
		},
	}

	var events []SSEvent
	consumer := func(event SSEvent) error {
		events = append(events, event)
		return nil
	}

	err := ParseAgenticMessage(msg, consumer)
	require.NoError(t, err)
	require.Len(t, events, 1)
	assert.Equal(t, EventToolCall, events[0].Type)
	assert.Equal(t, "get_weather", events[0].Meta["tool_name"])
	assert.Equal(t, "call_2", events[0].Meta["call_id"])
}

func TestParseAgenticMessage_MultiModal(t *testing.T) {
	msg := &schema.AgenticMessage{
		Role: schema.AgenticRoleTypeAssistant,
		ContentBlocks: []*schema.ContentBlock{
			schema.NewContentBlock(&schema.AssistantGenText{Text: "这是一段音频："}),
			schema.NewContentBlock(&schema.AssistantGenAudio{
				URL:      "https://example.com/audio.wav",
				MIMEType: "audio/wav",
			}),
			schema.NewContentBlock(&schema.AssistantGenVideo{
				Base64Data: "base64encodedvideo",
				MIMEType:   "video/mp4",
			}),
		},
	}

	var events []SSEvent
	consumer := func(event SSEvent) error {
		events = append(events, event)
		return nil
	}

	err := ParseAgenticMessage(msg, consumer)
	require.NoError(t, err)
	require.Len(t, events, 3)
	assert.Equal(t, EventText, events[0].Type)
	assert.Equal(t, EventAudio, events[1].Type)
	assert.Equal(t, "https://example.com/audio.wav", events[1].URL)
	assert.Equal(t, "audio/wav", events[1].MIMEType)
	assert.Equal(t, EventVideo, events[2].Type)
	assert.Equal(t, "base64encodedvideo", events[2].Base64)
	assert.Equal(t, "video/mp4", events[2].MIMEType)
}

func stringPtr(s string) *string {
	return &s
}
