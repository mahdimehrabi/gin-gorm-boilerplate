package core

import (
	"boilerplate/core/infrastructure"
	"boilerplate/core/validators"
	"boilerplate/docs"
	"context"
	"fmt"

	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"go.uber.org/fx"
)

var BootstrapModule = fx.Options(
	infrastructure.Module,
	RoutesModule,
	ControllerModule,
	SeviceModule,
	RepositoryModule,
	MiddlewaresModule,
	validators.Module,
	fx.Invoke(bootstrap),
)

func bootstrap(lifecycle fx.Lifecycle, database infrastructure.Database,
	middlewares Middlewares, router infrastructure.Router,
	routes Routes, env infrastructure.Env, logger infrastructure.Logger,
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
				docs.SwaggerInfo.BasePath = "/api"
				router.Gin.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

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
