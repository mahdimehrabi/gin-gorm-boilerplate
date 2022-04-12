package core

import (
	"boilerplate/apps/genericApp/middlewares"

	"go.uber.org/fx"
)

// Module Middleware exported
var MiddlewaresModule = fx.Options(
	fx.Provide(NewMiddlewares),
	fx.Provide(middlewares.NewDBTransactionMiddleware),
	fx.Provide(middlewares.NewAuthMiddleware),
	fx.Provide(middlewares.NewAdminMiddleware),
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
