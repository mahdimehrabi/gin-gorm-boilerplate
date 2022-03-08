package routes

import (
	"boilerplate/api/controllers"
	"boilerplate/infrastructure"
)

// AuthRoutes -> utility routes struct
type AuthRoutes struct {
	router         infrastructure.Router
	Logger         infrastructure.Logger
	AuthController controllers.AuthController
}

//NewAuthRoute -> returns new utility route
func NewAuthRoutes(
	logger infrastructure.Logger,
	router infrastructure.Router,
	AuthController controllers.AuthController,
) AuthRoutes {
	return AuthRoutes{
		Logger:         logger,
		router:         router,
		AuthController: AuthController,
	}
}

//Setup -> sets up route for util entities
func (gr AuthRoutes) Setup() {
	gr.router.Gin.POST("/api/auth/register", gr.AuthController.Register)
	gr.router.Gin.POST("/api/auth/login", gr.AuthController.Login)
	gr.router.Gin.POST("/api/auth/access-token-verify", gr.AuthController.AccessTokenVerify)
	gr.router.Gin.POST("/api/auth/renew-access-token", gr.AuthController.RenewToken)
}
