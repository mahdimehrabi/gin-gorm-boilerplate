package controllers

import (
	"boilerplate/api/services"
	"boilerplate/infrastructure"
	"boilerplate/models"
	"boilerplate/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	logger      infrastructure.Logger
	env         infrastructure.Env
	userService services.UserService
	authService services.AuthService
}

func NewAuthController(logger infrastructure.Logger,
	env infrastructure.Env,
	userService services.UserService,
	authService services.AuthService,
) AuthController {
	return AuthController{
		logger:      logger,
		env:         env,
		userService: userService,
		authService: authService,
	}
}

func (ac AuthController) Register(c *gin.Context) {

	// Data Parse
	var userData models.CreateUser
	err := c.ShouldBindJSON(&userData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	encodedPassword := utils.Sha256Encrypt(userData.Password)
	user.Password = encodedPassword
	user.FullName = userData.FullName
	user.Email = userData.Email
	err = ac.userService.CreateUser(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// login
	// token
	accessToken, refreshToken, err := ac.authService.CreateTokens(int(user.Base.ID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	var loginResult models.LoginResult
	loginResult.AccessToken = accessToken
	loginResult.RefreshToken = refreshToken

	c.JSON(http.StatusOK, loginResult)
}
