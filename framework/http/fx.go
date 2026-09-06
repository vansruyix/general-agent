// Package http
// @author: fengyi
// @date: 2024/7/9
// @note:
package http

import (
	"go.uber.org/fx"
)

var Module = fx.Module("http",
	fx.Provide(NewDefaultEngine),
	fx.Provide(NewDefaultRouterGroup),
	fx.Provide(NewEngineServer),
)
