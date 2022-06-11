package tasks

import (
	"boilerplate/core/infrastructure"
	"github.com/hibiken/asynq"
)

//Tasks -> Tasks Struct
type Tasks struct {
	Logger    infrastructure.Logger
	Env       infrastructure.Env
	Server    *asynq.Server
	ServerMux *asynq.ServeMux
}

//NewTasks -> return new Tasks struct,
func NewTasks(
	logger infrastructure.Logger,
	env infrastructure.Env,
) Tasks {
	return Tasks{
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
		ServerMux: asynq.NewServeMux(),
	}
}

//GetClient -> return asynq client don't forget to close it
func (t *Tasks) GetClient() *asynq.Client {
	return asynq.NewClient(asynq.RedisClientOpt{Addr: t.Env.RedisAddr})
}
