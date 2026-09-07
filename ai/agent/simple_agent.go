package agent

import (
	"context"
	"general-agent/ai/chat"

	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/schema"
)

func SimpleAgent() {
	ctx := context.Background()
	cm, _ := chat.NewSimpleChatModel(ctx)
	agent, _ := adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
		Name:        "",
		Description: "",
		Instruction: "",
		Model:       cm,
	})

	runner := adk.NewRunner(ctx, adk.RunnerConfig{
		Agent:           agent,
		EnableStreaming: true,
	})
	history := make([]*schema.Message, 0, 16)
	events := runner.Run(ctx, history)
	AsyncIteraorHandler(events, nil)
}
