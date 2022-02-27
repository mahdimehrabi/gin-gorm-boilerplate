package routes

import (
	"boilerplate/api/controllers"
	"boilerplate/infrastructure"
)

// GenericRoutes -> utility routes struct
type GenericRoutes struct {
	router            infrastructure.Router
	Logger            infrastructure.Logger
	GenericController controllers.GenericController
}

//NewGenericRoute -> returns new utility route
func NewGenericRoutes(
	logger infrastructure.Logger,
	router infrastructure.Router,
	GenericController controllers.GenericController,
) GenericRoutes {
	return GenericRoutes{
		Logger:            logger,
		router:            router,
		GenericController: GenericController,
	}
}

//Setup -> sets up route for util entities
func (gr GenericRoutes) Setup() {
	gr.router.Gin.GET("/ping", gr.GenericController.Ping)
}
