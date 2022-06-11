package tasks

import (
	"boilerplate/core/infrastructure"
	"github.com/hibiken/asynq"
	"go.uber.org/fx"
)

var TasksModules = fx.Options(
	fx.Provide(NewTasks),
)

type Task interface {
	HandlesToMux() error
}

type Tasks struct {
	logger    infrastructure.Logger
	taskAsynq TaskAsynq
}

func NewTasks(logger infrastructure.Logger, taskAsynq TaskAsynq) Tasks {
	return Tasks{
		logger:    logger,
		taskAsynq: taskAsynq,
	}
}

func (t *Tasks) HandleTasks() error {
	emailTask := NewEmailTask(t.logger)
	serverMux := asynq.NewServeMux()
	serverMux.HandleFunc(
		TypeSendVerifyEmail,
		emailTask.HandleVerifyEmailTask,
	)
	return t.taskAsynq.Server.Run(serverMux)
}
