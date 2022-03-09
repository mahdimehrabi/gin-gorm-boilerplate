package routes

import (
	"boilerplate/api/controllers"
	"boilerplate/api/middlewares"
	"boilerplate/infrastructure"
)

// UserRoutes -> struct
type UserRoutes struct {
	logger         infrastructure.Logger
	router         infrastructure.Router
	userController controllers.UserController
	trxMiddleware  middlewares.DBTransactionMiddleware
	authMiddleware middlewares.AuthMiddleware
}

// Setup user routes
func (ur UserRoutes) Setup() {
	ur.logger.Zap.Info("Setting up user routes 👷")
	users := ur.router.Gin.Group("/api/users").Use(ur.authMiddleware.AuthHandle())
	{
		users.GET("", ur.userController.GetAllUsers)
		users.POST("", ur.trxMiddleware.DBTransactionHandle(), ur.userController.CreateUser)
		users.DELETE(":id", ur.userController.DeleteUser)
	}
}

// NewUserRoutes -> creates new user controller
func NewUserRoutes(
	logger infrastructure.Logger,
	router infrastructure.Router,
	userController controllers.UserController,
	trxMiddleware middlewares.DBTransactionMiddleware,
	authMiddleware middlewares.AuthMiddleware,
) UserRoutes {
	return UserRoutes{
		router:         router,
		logger:         logger,
		userController: userController,
		trxMiddleware:  trxMiddleware,
		authMiddleware: authMiddleware,
	}
}
