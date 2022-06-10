package infrastructure

import (
	"github.com/hibiken/asynq"
	"go.uber.org/fx"
)

var TasksModules = fx.Options(
	fx.Provide(NewTasks),
	fx.Provide(NewEmail),
)

//Tasks -> Tasks Struct
type Tasks struct {
	Logger Logger
	Env    Env
	Client *asynq.Client
	Server *asynq.Server
}

//NewTasks -> return new Tasks struct
func NewTasks(
	logger Logger,
	env Env,
) Tasks {
	return Tasks{
		Logger: logger,
		Env:    env,
		Client: asynq.NewClient(asynq.RedisClientOpt{Addr: env.RedisAddr}),
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
