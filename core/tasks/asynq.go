package tasks

import (
	"boilerplate/core/infrastructure"
	"github.com/hibiken/asynq"
)

//TaskAsynq -> TaskAsynq Struct
type TaskAsynq struct {
	Logger infrastructure.Logger
	Env    infrastructure.Env
	Server *asynq.Server
}

//NewTaskAsynq -> return new TaskAsynq struct,
func NewTaskAsynq(
	logger infrastructure.Logger,
	env infrastructure.Env,
) TaskAsynq {
	return TaskAsynq{
		Logger: logger,
		Env:    env,
		Server: asynq.NewServer(asynq.RedisClientOpt{Addr: env.RedisAddr},
			asynq.Config{
				Concurrency: 10,
				Queues: map[string]int{
					"critical": 6,
					"default":  3,
					"info":     1,
				},
			},
		),
	}
}

//NewClient -> return asynq client don't forget to close it
func (t *TaskAsynq) NewClient() *asynq.Client {
	return asynq.NewClient(asynq.RedisClientOpt{Addr: t.Env.RedisAddr})
}
