package app

import (
	"general-agent/app/user"

	"go.uber.org/fx"
)

var Module = fx.Module("app",
	user.Module,
)
