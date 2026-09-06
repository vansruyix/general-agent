package chat

import (
	"context"
	"fmt"
	"io"

	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/schema"
)

func NewSimpleChatModel(ctx context.Context) (*openai.ChatModel, error) {

	cm, err := openai.NewChatModel(ctx, &openai.ChatModelConfig{
		// 配置参数
		APIKey:  "",
		BaseURL: "",
		Model:   "",
	})
	return cm, err
}

func SimpleChat() {
	// 初始化模型 (以openai为例)
	ctx := context.Background()
	cm, _ := NewSimpleChatModel(ctx)

	// 准备输入消息
	messages := []*schema.Message{
		{
			Role:    schema.System,
			Content: "你是一个有帮助的助手。",
		},
		{
			Role:    schema.User,
			Content: "你好！",
		},
	}

	// 生成响应
	response, _ := cm.Generate(ctx, messages, model.WithTemperature(0.8))

	// 响应处理
	fmt.Print(response.Content)

	// 流式生成
	streamResult, _ := cm.Stream(ctx, messages)

	defer streamResult.Close()

	for {
		chunk, err := streamResult.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			// 错误处理
		}
		// 响应片段处理
		fmt.Print(chunk.Content)
	}
}
