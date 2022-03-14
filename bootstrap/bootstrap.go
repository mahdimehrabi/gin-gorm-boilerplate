package bootstrap

import (
	"boilerplate/api/controllers"
	"boilerplate/api/middlewares"
	"boilerplate/api/repositories"
	"boilerplate/api/routes"
	"boilerplate/api/services"
	"boilerplate/api/validators"
	"boilerplate/infrastructure"
	"context"
	"fmt"

	"go.uber.org/fx"
)

var Module = fx.Options(
	infrastructure.Module,
	routes.Module,
	controllers.Module,
	services.Module,
	repositories.Module,
	middlewares.Module,
	validators.Module,
	fx.Invoke(bootstrap),
)

func bootstrap(lifecycle fx.Lifecycle, database infrastructure.Database,
	middlewares middlewares.Middlewares, router infrastructure.Router,
	routes routes.Routes, env infrastructure.Env, logger infrastructure.Logger,
	validators validators.Validators) {
	appStop := func(context.Context) error {
		logger.Zap.Info("Stopping Application 📛")
		conn, _ := database.DB.DB()
		conn.Close()
		return nil
	}

	lifecycle.Append(fx.Hook{
		OnStart: func(context.Context) error {
			logger.Zap.Info("Starting Application")
			logger.Zap.Info("------------------------")
			logger.Zap.Info(fmt.Sprintf("------ %s  ------", env.AppName))
			logger.Zap.Info("------------------------")
			go func() {
				validators.Setup()
				routes.Setup()
				middlewares.Setup()
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
