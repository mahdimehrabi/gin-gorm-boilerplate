package controllers

import (
	"boilerplate/core/infrastructure"
	_ "boilerplate/core/models"
	"boilerplate/core/responses"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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
	uc.logger.Zap.Error(gorm.ErrInvalidDB)
	responses.JSON(ctx, http.StatusOK, gin.H{"pingpong": "🏓🏓🏓🏓🏓🏓"}, "pong")
}
