package core

import (
	authServices "boilerplate/apps/authApp/services"
	userServices "boilerplate/apps/userApp/services"

	"go.uber.org/fx"
)

var SeviceModule = fx.Options(
	fx.Provide(userServices.NewUserService),
	fx.Provide(authServices.NewAuthService),
)
