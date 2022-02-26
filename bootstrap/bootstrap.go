package bootstrap

import (
	"boilerplate/infrastructure"
	"context"

	"go.uber.org/fx"
)

var Module = fx.Options(
	infrastructure.Module,
	fx.Invoke(bootstrap),
)

func bootstrap(lifecycle fx.Lifecycle, router infrastructure.Router, env infrastructure.Env, logger infrastructure.Logger) {
	appStop := func(context.Context) error {
		logger.Zap.Info("Stopping Application")
		/*conn, _ := database.DB.DB()
		conn.Close()*/
		return nil
	}

	lifecycle.Append(fx.Hook{
		OnStart: func(context.Context) error {
			logger.Zap.Info("Starting Application")
			logger.Zap.Info("------------------------")
			logger.Zap.Info("------ Boilerplate 📺 ------")
			logger.Zap.Info("------------------------")
			go func() {
				if env.ServerPort == "" {
					router.Gin.Run(":5000")
				} else {
					router.Gin.Run(":" + env.ServerPort)
				}
			}()
			return nil
		},
		OnStop: appStop,
	})
}
