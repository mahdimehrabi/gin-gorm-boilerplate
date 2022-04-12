package routes

import (
	"boilerplate/api/controllers/dashboard/admin"
	"boilerplate/api/middlewares"
	"boilerplate/core/infrastructure"
)

// AdminRoutes -> utility routes struct
type AdminRoutes struct {
	router          infrastructure.Router
	logger          infrastructure.Logger
	userController  admin.UserController
	authMiddleware  middlewares.AuthMiddleware
	AdminMiddleware middlewares.AdminMiddleware
}

//NewAdminRoute -> returns new utility route
func NewAdminRoutes(
	logger infrastructure.Logger,
	router infrastructure.Router,
	userController admin.UserController,
	authMiddleware middlewares.AuthMiddleware,
	adminMiddleware middlewares.AdminMiddleware,
) AdminRoutes {
	return AdminRoutes{
		logger:          logger,
		router:          router,
		userController:  userController,
		authMiddleware:  authMiddleware,
		AdminMiddleware: adminMiddleware,
	}
}

//Setup -> sets up route for util entities
func (er AdminRoutes) Setup() {
	g := er.router.Gin.Group("/api/admin").
		Use(er.authMiddleware.AuthHandle()).
		Use(er.AdminMiddleware.AdminHandle())
	{
		g.POST("/users", er.userController.CreateUser)
		g.GET("/users", er.userController.ListUser)
		g.PUT("/users/:id", er.userController.UpdateUser)
		g.DELETE("/users/:id", er.userController.DeleteUser)
	}
}
