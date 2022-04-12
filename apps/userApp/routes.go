package userApp

import (
	"boilerplate/apps/genericApp/middlewares"
	"boilerplate/apps/userApp/controllers"
	"boilerplate/core/infrastructure"
)

// UserRoutes -> utility routes struct
type UserRoutes struct {
	router            infrastructure.Router
	Logger            infrastructure.Logger
	ProfileController controllers.ProfileController
	authMiddleware    middlewares.AuthMiddleware
}

//NewProfileRoute -> returns new utility route
func NewUserRoutes(
	logger infrastructure.Logger,
	router infrastructure.Router,
	ProfileController controllers.ProfileController,
	authMiddleware middlewares.AuthMiddleware,
) UserRoutes {
	return UserRoutes{
		Logger:            logger,
		router:            router,
		ProfileController: ProfileController,
		authMiddleware:    authMiddleware,
	}
}

//Setup -> sets up route for util entities
func (pr UserRoutes) Setup() {
	g := pr.router.Gin.Group("/api/profile").Use(pr.authMiddleware.AuthHandle())
	{
		g.POST("/change-password", pr.ProfileController.ChangePassword)
		g.GET("/devices", pr.ProfileController.LoggedInDevices)
		g.POST("/terminate-device", pr.ProfileController.TerminateDevice)
		g.POST("/terminate-devices-except-me", pr.ProfileController.TerminateDevicesExceptMe)
		g.POST("/upload-profile-picture", pr.ProfileController.UploadProfilePicture)
	}
}
