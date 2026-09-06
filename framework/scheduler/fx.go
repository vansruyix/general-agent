package scheduler

import (
	"go.uber.org/fx"
)

var Module = fx.Module("scheduler",
	fx.Provide(NewConfJobDao),
	fx.Provide(NewJobManager),
	fx.Invoke(registerJob),
)
