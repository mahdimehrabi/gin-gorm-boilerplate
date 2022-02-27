package controllers

import (
	"boilerplate/infrastructure"
	"net/http"

	"github.com/gin-gonic/gin"
)

type GenericController struct {
	logger infrastructure.Logger
	env    infrastructure.Env
}

func NewGenericController(logger infrastructure.Logger,
	env infrastructure.Env,
) GenericController {
	return GenericController{
		logger: logger,
		env:    env,
	}
}

func (uc GenericController) Ping(ctx *gin.Context) {
	ctx.String(http.StatusOK, "Pong 🏓")
}
