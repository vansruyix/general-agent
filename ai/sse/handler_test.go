package sse

import (
	"testing"

	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandleMessageVariant_NonStreaming_Message(t *testing.T) {
	mv := &adk.TypedMessageVariant[*schema.Message]{
		IsStreaming: false,
		Message: &schema.Message{
			Role:    schema.Assistant,
			Content: "Hello from variant",
		},
		Role: schema.Assistant,
	}

	var events []SSEvent
	consumer := func(event SSEvent) error {
		events = append(events, event)
		return nil
	}

	result, err := HandleMessageVariant(mv, consumer)
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Len(t, events, 1)
	assert.Equal(t, EventText, events[0].Type)
	assert.Equal(t, "Hello from variant", events[0].Data)
	assert.Equal(t, "Hello from variant", result.Content)
}

func TestHandleMessageVariant_Streaming_Message(t *testing.T) {
	sr, sw := schema.Pipe[*schema.Message](3)
	sw.Send(&schema.Message{Role: schema.Assistant, Content: "Hello "}, nil)
	sw.Send(&schema.Message{Role: schema.Assistant, Content: "stream"}, nil)
	sw.Close()

	mv := &adk.TypedMessageVariant[*schema.Message]{
		IsStreaming:   true,
		MessageStream: sr,
		Role:          schema.Assistant,
	}

	var events []SSEvent
	consumer := func(event SSEvent) error {
		events = append(events, event)
		return nil
	}

	result, err := HandleMessageVariant(mv, consumer)
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Len(t, events, 2)
	assert.Equal(t, "Hello ", events[0].Data)
	assert.Equal(t, "stream", events[1].Data)
	assert.Equal(t, "Hello stream", result.Content)
}

func TestHandleMessageVariant_NonStreaming_AgenticMessage(t *testing.T) {
	mv := &adk.TypedMessageVariant[*schema.AgenticMessage]{
		IsStreaming: false,
		Message: &schema.AgenticMessage{
			Role: schema.AgenticRoleTypeAssistant,
			ContentBlocks: []*schema.ContentBlock{
				schema.NewContentBlock(&schema.Reasoning{Text: "thinking..."}),
				schema.NewContentBlock(&schema.AssistantGenText{Text: "agentic answer"}),
			},
		},
	}

	var events []SSEvent
	consumer := func(event SSEvent) error {
		events = append(events, event)
		return nil
	}

	result, err := HandleMessageVariant(mv, consumer)
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Len(t, events, 2)
	assert.Equal(t, EventThinking, events[0].Type)
	assert.Equal(t, EventText, events[1].Type)
	assert.Equal(t, "agentic answer", result.ContentBlocks[1].AssistantGenText.Text)
}

func TestHandleMessageVariant_Streaming_AgenticMessage(t *testing.T) {
	sr, sw := schema.Pipe[*schema.AgenticMessage](3)
	sw.Send(&schema.AgenticMessage{
		Role: schema.AgenticRoleTypeAssistant,
		ContentBlocks: []*schema.ContentBlock{
			schema.NewContentBlock(&schema.AssistantGenText{Text: "part1 "}),
		},
	}, nil)
	sw.Send(&schema.AgenticMessage{
		Role: schema.AgenticRoleTypeAssistant,
		ContentBlocks: []*schema.ContentBlock{
			schema.NewContentBlock(&schema.AssistantGenText{Text: "part2"}),
		},
	}, nil)
	sw.Close()

	mv := &adk.TypedMessageVariant[*schema.AgenticMessage]{
		IsStreaming:   true,
		MessageStream: sr,
	}

	var events []SSEvent
	consumer := func(event SSEvent) error {
		events = append(events, event)
		return nil
	}

	result, err := HandleMessageVariant(mv, consumer)
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Len(t, events, 2)
	assert.Equal(t, "part1 ", events[0].Data)
	assert.Equal(t, "part2", events[1].Data)
	require.Len(t, result.ContentBlocks, 2)
}
