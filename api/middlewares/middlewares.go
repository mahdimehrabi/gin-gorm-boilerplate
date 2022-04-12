package middlewares

import "go.uber.org/fx"

// Module Middleware exported
var Module = fx.Options(
	fx.Provide(NewMiddlewares),
	fx.Provide(NewDBTransactionMiddleware),
	fx.Provide(NewAuthMiddleware),
	fx.Provide(NewAdminMiddleware),
)

// IMiddleware middleware interface
type Middleware interface {
	Setup()
}

// Middlewares contains multiple middleware
type Middlewares []Middleware

// NewMiddlewares creates new middlewares
// Register the middleware that should be applied directly (globally)
func NewMiddlewares() Middlewares {
	return Middlewares{}
}

// Setup sets up middlewares
func (m Middlewares) Setup() {
	for _, middleware := range m {
		middleware.Setup()
	}
}
