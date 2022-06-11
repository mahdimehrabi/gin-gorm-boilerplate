package core

import (
	"boilerplate/core/infrastructure"
	"boilerplate/core/responses"
	"boilerplate/core/tasks"
	"boilerplate/core/validators"
	"boilerplate/docs"
	"context"
	"fmt"
	"net/http"
	"runtime"

	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-gonic/gin"
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
	tasks.Modules,
	fx.Invoke(bootstrap),
)

func bootstrap(lifecycle fx.Lifecycle, database infrastructure.Database,
	middlewares Middlewares, router infrastructure.Router,
	routes Routes, env infrastructure.Env,
	logger infrastructure.Logger, validators validators.Validators,
	taskAsynq tasks.TaskAsynq,
	emailTask tasks.EmailTask,
) {
	appStop := func(context.Context) error {
		logger.Zap.Info("Stopping Application 📛")
		conn, _ := database.DB.DB()
		conn.Close()
		return nil
	}

	//recover unwanted 500 errors
	router.Gin.Use(gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		if recovered != nil {
			switch e := recovered.(type) {
			case string:
				logger.Zap.Warn("recovered (string) panic:", e)
			case runtime.Error:
				logger.Zap.Warn("recovered (runtime.Error) panic:", e.Error())
			case error:
				logger.Zap.Warn("recovered (error) panic:", e.Error())
			default:
				logger.Zap.Warn("recovered (default) panic:", e)
			}
			responses.ErrorJSON(c, http.StatusInternalServerError, gin.H{}, "Sorry an error occured!")
			return
		}
		c.AbortWithStatus(http.StatusInternalServerError)
	}))

	//recover unwanted 500 errors to sentery
	router.Gin.Use(sentrygin.New(sentrygin.Options{
		Repanic: true,
	}))

	tasksStruct := tasks.NewTasks(logger, taskAsynq, emailTask)
	err := tasksStruct.HandleTasks()
	if err != nil {
		logger.Zap.Error("Failed to run asynq handlers:", err)
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
