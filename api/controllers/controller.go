package controllers

import (
	"boilerplate/api/controllers/dashboard/admin"

	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(NewGenericController),
	fx.Provide(NewAuthController),
	fx.Provide(NewProfileController),
	fx.Provide(admin.NewUserController),
)
