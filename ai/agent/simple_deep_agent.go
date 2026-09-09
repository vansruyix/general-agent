package agent

import (
	"context"
	"general-agent/ai/chat"

	"github.com/cloudwego/eino-ext/adk/backend/local"
	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/adk/prebuilt/deep"
)

func NewSimpleDeepAgent() *adk.Runner {
	ctx := context.Background()
	cm, _ := chat.NewSimpleChatModel(ctx)
	// 创建 LocalBackend
	backend, _ := local.NewBackend(ctx, &local.Config{})

	// 创建 DeepAgent,自动注册文件系统工具
	agent, _ := deep.New(ctx, &deep.Config{
		Name:           "GeneralToolAgent",
		Description:    "ChatWithDoc agent with filesystem access via LocalBackend.",
		ChatModel:      cm,
		Instruction:    "你是一个通用Agent智能体，擅长编码、逻辑推理等任务，根据用户提问你可以调用工具来帮助你完成任务！",
		Backend:        backend, // 提供文件系统操作能力
		StreamingShell: backend, // 提供命令执行能力
		MaxIteration:   50,
	})

	runner := adk.NewRunner(ctx, adk.RunnerConfig{
		Agent:           agent,
		EnableStreaming: true,
	})
	return runner

}
