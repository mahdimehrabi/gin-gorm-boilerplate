package routes

import (
	"boilerplate/api/controllers"
	"boilerplate/api/middlewares"
	"boilerplate/infrastructure"
)

// AuthRoutes -> utility routes struct
type AuthRoutes struct {
	router         infrastructure.Router
	Logger         infrastructure.Logger
	AuthController controllers.AuthController
	authMiddleware middlewares.AuthMiddleware
}

//NewAuthRoute -> returns new utility route
func NewAuthRoutes(
	logger infrastructure.Logger,
	router infrastructure.Router,
	AuthController controllers.AuthController,
	authMiddleware middlewares.AuthMiddleware,

) AuthRoutes {
	return AuthRoutes{
		Logger:         logger,
		router:         router,
		AuthController: AuthController,
		authMiddleware: authMiddleware,
	}
}

//Setup -> sets up route for util entities
func (ar AuthRoutes) Setup() {
	ar.router.Gin.POST("/api/auth/register", ar.AuthController.Register)
	ar.router.Gin.POST("/api/auth/login", ar.AuthController.Login)
	ar.router.Gin.POST("/api/auth/access-token-verify", ar.AuthController.AccessTokenVerify)
	ar.router.Gin.POST("/api/auth/renew-access-token", ar.AuthController.RenewToken)
	ar.router.Gin.POST("/api/auth/logout", ar.authMiddleware.AuthHandle(), ar.AuthController.Logout)
	ar.router.Gin.POST("/api/auth/verify-email", ar.AuthController.VerifyEmail)
	ar.router.Gin.POST("/api/auth/forgot-password", ar.AuthController.ForgotPassword)
	ar.router.Gin.POST("/api/auth/recover-password", ar.AuthController.RecoverPassword)
}
