package controllers

import (
	"boilerplate/api/repositories"
	"boilerplate/models"

	"github.com/gin-gonic/gin"
)

func GetUser(c gin.Context, ur repositories.UserRepository) (models.User, error) {
	userId := c.MustGet("userId")
	return ur.FindByField("id", userId.(string))
}
