package tasks

import (
	"boilerplate/core/infrastructure"
	"github.com/hibiken/asynq"
	"go.uber.org/fx"
)

var Modules = fx.Options(
	fx.Provide(NewTaskAsynq),
	fx.Provide(NewEmailTask),
)

type Task interface {
	HandlesToMux() error
}

type Tasks struct {
	logger    infrastructure.Logger
	taskAsynq TaskAsynq
	emailTask EmailTask
}

func NewTasks(
	logger infrastructure.Logger,
	taskAsynq TaskAsynq,
	emailTask EmailTask) Tasks {
	return Tasks{
		logger:    logger,
		taskAsynq: taskAsynq,
	}
}

func (t *Tasks) HandleTasks() error {
	serverMux := asynq.NewServeMux()
	serverMux.HandleFunc(
		TypeSendVerifyEmail,
		t.emailTask.HandleVerifyEmailTask,
	)
	return t.taskAsynq.Server.Run(serverMux)
}
