package controllers

import "go.uber.org/fx"

var Module = fx.Options(
	fx.Provide(NewGenericController),
	fx.Provide(NewAuthController),
	fx.Provide(NewProfileController),
)
