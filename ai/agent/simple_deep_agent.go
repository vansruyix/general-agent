package agent

import (
	"context"
	"fmt"
	"general-agent/ai/chat"
	"general-agent/ai/memory"
	"general-agent/ai/sse"

	"github.com/cloudwego/eino-ext/adk/backend/local"
	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/adk/prebuilt/deep"
	"github.com/cloudwego/eino/schema"
)

func NewSimpleDeepAgent(id, question string) {
	ctx := context.Background()
	cm, _ := chat.NewSimpleChatModel(ctx)
	// 创建 LocalBackend
	backend, _ := local.NewBackend(ctx, &local.Config{})

	// 创建 DeepAgent,自动注册文件系统工具
	agent, _ := deep.New(ctx, &deep.Config{
		Name:           "Ch04ToolAgent",
		Description:    "ChatWithDoc agent with filesystem access via LocalBackend.",
		ChatModel:      cm,
		Instruction:    "agentInstruction",
		Backend:        backend, // 提供文件系统操作能力
		StreamingShell: backend, // 提供命令执行能力
		MaxIteration:   50,
	})

	runner := adk.NewRunner(ctx, adk.RunnerConfig{
		Agent:           agent,
		EnableStreaming: true,
	})

	userMsg := &schema.Message{
		Role:    schema.User,
		Content: question,
	}

	userMemory := memory.GetMemory(id)
	history := make([]*schema.Message, len(userMemory.Messages))
	history = append(history, userMsg)
	events := runner.Run(ctx, history)
	var eventType sse.EventType
	AsyncIteraorHandler(events, func(event sse.SSEvent) error {
		if eventType == "" {
			eventType = event.Type
			switch event.Type {
			case sse.EventThinking:
				fmt.Printf("<think>%s", event.Data)
			case sse.EventToolCall:
				fmt.Printf("[Tool Caller: %s]", event.Data)
			}
			return nil
		}
		if event.Type == eventType {
			fmt.Print(event.Data)
		} else {
			switch eventType {
			case sse.EventThinking:
				fmt.Printf("</think>")
			}
			fmt.Println()
			eventType = event.Type
			switch event.Type {
			case sse.EventThinking:
				fmt.Printf("<think>%s", event.Data)
			case sse.EventToolCall:
				fmt.Printf("[Tool Caller: %s]", event.Data)
			case sse.EventToolResult:
				fmt.Printf("[Tool Result: %s]", event.Data)
			}
		}
		return nil
	})

}
