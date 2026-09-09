package agent

import (
	"general-agent/ai/memory"
	"general-agent/ai/sse"
	"general-agent/internal/response"

	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/schema"
	"github.com/gin-gonic/gin"

	msgHandler "general-agent/ai/agent"
)

type ChatAgentController struct {
	runner *adk.Runner
}

func NewChatAgentController(runner *adk.Runner) *ChatAgentController {
	return &ChatAgentController{runner: runner}
}

func (c *ChatAgentController) Chat(ctx *gin.Context) {
	// 设置SSE必要响应头
	ctx.Header("Content-Type", "text/event-stream")
	ctx.Header("Cache-Control", "no-cache")
	ctx.Header("Connection", "keep-alive")

	var req ChatCommonReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Fail(ctx, err)
		return
	}
	userMsg := &schema.Message{
		Role:    schema.User,
		Content: req.Question,
	}

	userMemory := memory.GetMemory(req.ID)
	history := make([]*schema.Message, len(userMemory.Messages))
	copy(history, userMemory.Messages)
	history = append(history, userMsg)
	events := c.runner.Run(ctx, history)
	msg := msgHandler.MessageHandler(events, func(event sse.SSEvent) error {
		ctx.SSEvent(string(event.Type), event)
		ctx.Writer.Flush()
		return nil
	})
	memory.AppendMessage(req.ID, userMsg)
	memory.AppendMessage(req.ID, msg)
	// var eventType sse.EventType
	// MessageHandler(events, ConsoleFmtConsumer(eventType))
}
