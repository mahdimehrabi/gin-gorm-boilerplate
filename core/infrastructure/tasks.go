package infrastructure

import (
	"github.com/hibiken/asynq"
	"go.uber.org/fx"
)

var TasksModules = fx.Options(
	fx.Provide(NewTasks),
	fx.Provide(NewEmail),
)

type Task interface {
	HandlesToMux() error
}

//Tasks -> Tasks Struct
type Tasks struct {
	Logger    Logger
	Env       Env
	Server    *asynq.Server
	ServerMux *asynq.ServeMux
}

//NewTasks -> return new Tasks struct
func NewTasks(
	logger Logger,
	env Env,
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
