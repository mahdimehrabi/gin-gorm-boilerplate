package infrastructure

import "go.uber.org/fx"

var Module = fx.Options(
	fx.Provide(NewLogger),
	fx.Provide(NewEnv),
	fx.Provide(NewRouter),
)
