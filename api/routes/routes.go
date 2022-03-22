package routes

import "go.uber.org/fx"

// Module exports dependency to container
var Module = fx.Options(
	fx.Provide(NewRoutes),
	fx.Provide(NewGenericRoutes),
	fx.Provide(NewUserRoutes),
	fx.Provide(NewAuthRoutes),
	fx.Provide(NewProfileRoutes),
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
	userRoutes UserRoutes,
	profileRoutes ProfileRoutes,
) Routes {
	return Routes{
		genericRoutes,
		userRoutes,
		authRoutes,
		profileRoutes,
	}
}

// Setup all the route
func (r Routes) Setup() {
	for _, route := range r {
		route.Setup()
	}
}
