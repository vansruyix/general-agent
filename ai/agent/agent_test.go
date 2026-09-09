package agent

import (
	"context"
	"fmt"
	"general-agent/ai/chat"
	"testing"

	"github.com/cloudwego/eino-ext/adk/backend/local"
	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/adk/prebuilt/deep"
	"github.com/cloudwego/eino/schema"
)

func TestDeepAgents(t *testing.T) {
	//NewSimpleDeepAgent(uuid.NewString(), "当前项目目录的结构是什么样的,全程使用中文交流")
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
	messages := []*schema.Message{schema.UserMessage("当前项目的目录结构是什么样")}
	events := runner.Run(ctx, messages)
	var lastMsg string
	detail := make([]string, 0)
	for {
		event, OK := events.Next()
		if !OK {
			break
		}
		if event.Err != nil {
			fmt.Errorf("%w", event.Err)
		}
		if event.Output != nil {
			msg, err := event.Output.MessageOutput.GetMessage()
			if err != nil {
				fmt.Errorf("%w", err)
			}
			lastMsg = msg.Content
			detail = append(detail, msg.Content)
		}
	}
	fmt.Print(lastMsg)
}
