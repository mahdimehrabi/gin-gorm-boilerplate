package core

import (
	"boilerplate/apps/userApp/repositories"

	"go.uber.org/fx"
)

// Module exports dependency
var RepositoryModule = fx.Options(
	fx.Provide(repositories.NewUserRepository),
)
