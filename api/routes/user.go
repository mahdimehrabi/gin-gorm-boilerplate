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
}

// Setup user routes
func (i UserRoutes) Setup() {
	i.logger.Zap.Info(" Setting up user routes")
	users := i.router.Gin.Group("/api/users")
	{
		users.GET("", i.userController.GetAllUsers)
		users.POST("", i.trxMiddleware.DBTransactionHandle(), i.userController.CreateUser)
		users.DELETE(":id", i.userController.DeleteUser)
	}
}

// NewUserRoutes -> creates new user controller
func NewUserRoutes(
	logger infrastructure.Logger,
	router infrastructure.Router,
	userController controllers.UserController,
	trxMiddleware middlewares.DBTransactionMiddleware,
) UserRoutes {
	return UserRoutes{
		router:         router,
		logger:         logger,
		userController: userController,
		trxMiddleware:  trxMiddleware,
	}
}
