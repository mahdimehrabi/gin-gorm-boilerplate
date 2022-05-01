package genericApp

import (
	"boilerplate/apps/genericApp/controllers"
	"boilerplate/core/infrastructure"
)

// GenericRoutes -> utility routes struct
type GenericRoutes struct {
	Router            infrastructure.Router
	Logger            infrastructure.Logger
	Env               infrastructure.Env
	GenericController controllers.GenericController
}

//NewGenericRoute -> returns new utility route
func NewGenericRoutes(
	logger infrastructure.Logger,
	router infrastructure.Router,
	env infrastructure.Env,
	genericController controllers.GenericController,
) GenericRoutes {
	return GenericRoutes{
		Logger:            logger,
		Router:            router,
		GenericController: genericController,
		Env:               env,
	}
}

//Setup -> sets up route for util entities
func (gr GenericRoutes) Setup() {
	gr.Router.Gin.GET("/api/ping", gr.GenericController.Ping)
	gr.Router.Gin.Static("/media", gr.Env.BasePath+"/media")
}
