package tasks

import (
	"boilerplate/core/infrastructure"
	"go.uber.org/fx"
)

var TasksModules = fx.Options(
	fx.Provide(NewTasks),
	fx.Provide(infrastructure.NewEmail),
)

type Task interface {
	HandlesToMux() error
}
