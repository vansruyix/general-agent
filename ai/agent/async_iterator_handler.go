package agent

import (
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/schema"
)

func AsyncIteraorHandler[M adk.MessageType](events *adk.AsyncIterator[*adk.TypedAgentEvent[M]]) {
	for {
		event, ok := events.Next()
		if !ok {
			break
		}
		if event.Err != nil {
			// 处理错误
		}
		if event.Output != nil && event.Output.MessageOutput != nil {
			// 处理消息输出（可能是流式）
			mo := event.Output.MessageOutput
			if mo.Role == schema.Tool {
				content, err := drainMessageVariant(mo)
				if err != nil {
					return
				}
				fmt.Printf("[Tool Result]\n%s\n\n", content)
				continue
			}
			if mo.Role != schema.Assistant && mo.Role != "" {
				continue
			}

			if mo.IsStreaming && mo.MessageStream != nil {
				defer mo.MessageStream.Close()
				var tools []schema.ToolCall
				for {
					frame, err := mo.MessageStream.Recv()
					if errors.Is(err, io.EOF) {
						break
					}
					if err != nil {
						return
					}
					if frame != nil {
						// fmt.Println(frame.Content)
						fmt.Println(extractMessageContent(frame))
					}
					switch f := any(frame).(type) {
					case *schema.Message:
						if len(f.ToolCalls) > 0 {
							tools = append(tools, f.ToolCalls...)
						}
					}
					// if len(frame.ToolCalls) > 0 {
					// 	tools = append(tools, frame.ToolCalls...)
					// }
				}
				fmt.Println()
				for _, tc := range tools {
					fmt.Printf("[Tool Result] %s(%s)\n", tc.Function.Name, tc.Function.Arguments)
				}
				continue
			}
		}
	}
}

func drainMessageVariant[M adk.MessageType](mv *adk.TypedMessageVariant[M]) (string, error) {
	msg, err := mv.GetMessage()
	if err != nil {
		return "", nil
	}
	return extractMessageContent(msg), nil
}

func extractMessageContent[M adk.MessageType](msg M) string {
	switch m := any(msg).(type) {
	case *schema.Message:
		if m == nil {
			return ""
		}
		return m.Content
	case *schema.AgenticMessage:
		if m == nil {
			return ""
		}
		var texts []string
		for _, block := range m.ContentBlocks {
			if block != nil && block.Type == schema.ContentBlockTypeAssistantGenText && block.AssistantGenText != nil {
				texts = append(texts, block.AssistantGenText.Text)
			}
		}
		return strings.Join(texts, "\n")
	}
	return ""
}

// func drainMessageVariant[M adk.MessageType](mv *adk.TypedMessageVariant[M]) (string, error) {
// 	if mv.Message != nil {
// 		return mv.Message.Content, nil
// 	}
// 	if !mv.IsStreaming || mv.MessageStream == nil {
// 		return "", nil
// 	}
// 	var sb strings.Builder
// 	for {
// 		chunk, err := mv.MessageStream.Recv()
// 		if errors.Is(err, io.EOF) {
// 			break
// 		}
// 		if err != nil {
// 			return "", nil
// 		}
// 		if chunk != nil && len(chunk.Content) > 0 {
// 			sb.WriteString(chunk.Content)
// 		}
// 	}
// 	return sb.String(), nil
// }
