package core

import (
	adminControllers "boilerplate/apps/adminApp/controllers"
	authControllers "boilerplate/apps/authApp/controllers"
	genericControllers "boilerplate/apps/genericApp/controllers"
	profileControllers "boilerplate/apps/userApp/controllers"

	"go.uber.org/fx"
)

var ControllerModule = fx.Options(
	fx.Provide(genericControllers.NewGenericController),
	fx.Provide(authControllers.NewAuthController),
	fx.Provide(profileControllers.NewProfileController),
	fx.Provide(adminControllers.NewUserController),
)
