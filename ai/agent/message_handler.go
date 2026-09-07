package agent

import (
	"fmt"

	"general-agent/ai/sse"

	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/schema"
)

// ConsoleConsumer 返回一个将 SSEvent 打印到控制台的 Consumer 实现。
// 作为默认消费函数，当调用方未传入自定义 Consumer 时使用。
func ConsoleConsumer() sse.Consumer {
	return func(event sse.SSEvent) error {
		switch event.Type {
		case sse.EventThinking:
			fmt.Printf("🤔 thinking: %s\n", event.Data)
		case sse.EventText:
			fmt.Print(event.Data)
		case sse.EventToolCall:
			fmt.Printf("\n🔧 [Tool Call] %s(%s)\n",
				event.Meta["tool_name"], event.Data)
		case sse.EventToolResult:
			fmt.Printf("📋 [Tool Result] %s\n", event.Data)
		case sse.EventImage:
			fmt.Printf("\n🖼️ [Image] %s (mime: %s)\n", event.URL, event.MIMEType)
		case sse.EventAudio:
			fmt.Printf("\n🎵 [Audio] %s (mime: %s)\n", event.URL, event.MIMEType)
		case sse.EventVideo:
			fmt.Printf("\n🎬 [Video] %s (mime: %s)\n", event.URL, event.MIMEType)
		case sse.EventError:
			fmt.Printf("\n❌ [Error] %s\n", event.Data)
		case sse.EventDone:
			fmt.Println()
		}
		return nil
	}
}

// StreamOutput 处理 TypedMessageVariant 并通过 consumer 消费。
// 返回完整聚合后的消息，供对话历史缓存使用。
func StreamOutput[M adk.MessageType](mv *adk.TypedMessageVariant[M], consumer sse.Consumer) (M, error) {
	return sse.HandleMessageVariant(mv, consumer)
}

// MessageOutput 解析单个 schema.Message 并通过 consumer 消费。
func MessageOutput(msg *schema.Message, consumer sse.Consumer) error {
	return sse.ParseMessage(msg, consumer)
}

// AsyncIteraorHandler 遍历 agent 事件流，通过 consumer 消费每个事件。
// 若 consumer 为 nil，默认使用 ConsoleConsumer 输出到控制台。
func AsyncIteraorHandler[M adk.MessageType](events *adk.AsyncIterator[*adk.TypedAgentEvent[M]], consumer sse.Consumer) {
	if consumer == nil {
		consumer = ConsoleConsumer()
	}

	for {
		event, ok := events.Next()
		if !ok {
			break
		}
		// 处理错误事件
		if event.Err != nil {
			consumer(sse.SSEvent{
				Type: sse.EventError,
				Data: event.Err.Error(),
			})
			continue
		}
		// 处理消息输出事件（流式或非流式）
		if event.Output != nil && event.Output.MessageOutput != nil {
			msg, err := sse.HandleMessageVariant(event.Output.MessageOutput, consumer)
			if err != nil {
				consumer(sse.SSEvent{
					Type: sse.EventError,
					Data: err.Error(),
				})
				continue
			}
			_ = msg // 调用方可使用返回的完整消息进行对话历史缓存
		}
	}
}
