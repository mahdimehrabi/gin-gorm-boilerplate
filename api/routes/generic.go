package routes

import (
	"boilerplate/api/controllers"
	docs "boilerplate/docs"
	"boilerplate/infrastructure"

	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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
	docs.SwaggerInfo.BasePath = "/api"
	gr.router.Gin.GET("/api/ping", gr.GenericController.Ping)
	gr.router.Gin.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
}
