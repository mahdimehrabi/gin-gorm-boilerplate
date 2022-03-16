package controllers

import (
	"boilerplate/api/responses"
	_ "boilerplate/api/responses/swagger"
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

// @basePath /api

// @Summary ping
// @Schemes
// @Description do ping
// @Tags generic
// @Accept json
// @Produce json
// @Success 200 {object} swagger.PingResponse
// @Router /ping [get]
func (uc GenericController) Ping(ctx *gin.Context) {
	responses.JSON(ctx, http.StatusOK, gin.H{"pingpong": "🏓🏓🏓🏓🏓🏓"}, "pong")
}
