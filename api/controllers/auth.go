package controllers

import (
	"boilerplate/api/responses"
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
	if err := c.ShouldBindJSON(&userData); err != nil {
		responses.ValidationErrorsJSON(c, err, "")
		return
	}

	var user models.User
	encodedPassword := utils.Sha256Encrypt(userData.Password)
	user.Password = encodedPassword
	user.FullName = userData.FullName
	user.Email = userData.Email
	err := ac.userService.CreateUser(&user)
	if err != nil {
		ac.logger.Zap.Error("Failed to create registered user", err.Error())
		responses.ErrorJSON(c, http.StatusInternalServerError, gin.H{}, "Sorry an error occoured in registering your account!")
		return
	}

	// login
	// token
	accessToken, refreshToken, err := ac.authService.CreateTokens(int(user.Base.ID))
	if err != nil {
		ac.logger.Zap.Error("Failed to generate registered user token", err.Error())
		responses.ErrorJSON(c, http.StatusInternalServerError, gin.H{}, "your account registerd but failed to make you login")
		return
	}
	var loginResult models.LoginResult
	loginResult.AccessToken = accessToken
	loginResult.RefreshToken = refreshToken

	responses.JSON(c, http.StatusOK, loginResult, "Your account created successfuly!")
}
