package ai

import (
	"general-agent/ai/agent"

	"go.uber.org/fx"
)

var Module = fx.Module("ai",
	fx.Provide(agent.NewSimpleDeepAgent),
)
