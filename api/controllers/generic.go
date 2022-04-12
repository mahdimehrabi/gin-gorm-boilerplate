package controllers

import (
	"boilerplate/api/responses"
	_ "boilerplate/api/responses/swagger"
	"boilerplate/core/infrastructure"
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

// @BasePath /api/auth
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
