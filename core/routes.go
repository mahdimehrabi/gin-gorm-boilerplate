package core

import (
	admin "boilerplate/apps/adminApp"
	"boilerplate/apps/authApp"
	"boilerplate/apps/genericApp"
	userApp "boilerplate/apps/userApp"

	"go.uber.org/fx"
)

// Module exports dependency to container
var RoutesModule = fx.Options(
	fx.Provide(NewRoutes),
	fx.Provide(genericApp.NewGenericRoutes),
	fx.Provide(authApp.NewAuthRoutes),
	fx.Provide(userApp.NewUserRoutes),
	fx.Provide(admin.NewAdminRoutes),
)

// Routes contains multiple routes
type Routes []Route

// Route interface
type Route interface {
	Setup()
}

// NewRoutes sets up routes
func NewRoutes(
	genericRoutes genericApp.GenericRoutes,
	authRoutes authApp.AuthRoutes,
	profileRoutes userApp.UserRoutes,
	adminRoutes admin.AdminRoutes,
) Routes {
	return Routes{
		genericRoutes,
		authRoutes,
		profileRoutes,
		adminRoutes,
	}
}

// Setup all the route
func (r Routes) Setup() {
	for _, route := range r {
		route.Setup()
	}
}
