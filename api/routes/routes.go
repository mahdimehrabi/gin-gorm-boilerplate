package routes

import "go.uber.org/fx"

// Module exports dependency to container
var Module = fx.Options(
	fx.Provide(NewRoutes),
	fx.Provide(NewGenericRoutes),
	fx.Provide(NewAuthRoutes),
	fx.Provide(NewProfileRoutes),
	fx.Provide(NewAdminRoutes),
)

// Routes contains multiple routes
type Routes []Route

// Route interface
type Route interface {
	Setup()
}

// NewRoutes sets up routes
func NewRoutes(
	genericRoutes GenericRoutes,
	authRoutes AuthRoutes,
	profileRoutes ProfileRoutes,
	adminRoutes AdminRoutes,
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
