package controllers

import (
	"boilerplate/api/repositories"
	"boilerplate/api/responses"
	"boilerplate/api/services"
	"boilerplate/infrastructure"
	"boilerplate/models"
	"boilerplate/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	logger         infrastructure.Logger
	env            infrastructure.Env
	userService    services.UserService
	authService    services.AuthService
	userRepository repositories.UserRepository
}

func NewAuthController(logger infrastructure.Logger,
	env infrastructure.Env,
	userService services.UserService,
	authService services.AuthService,
	userRepository repositories.UserRepository,
) AuthController {
	return AuthController{
		logger:         logger,
		env:            env,
		userService:    userService,
		authService:    authService,
		userRepository: userRepository,
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
	var loginResult models.LoginResponse
	loginResult.AccessToken = accessToken
	loginResult.RefreshToken = refreshToken
	loginResult.User = models.UserResponse(user)

	responses.JSON(c, http.StatusOK, loginResult, "Your account created successfuly!")
}

func (ac AuthController) Login(c *gin.Context) {
	// Data Parse
	var loginRquest models.LoginRequest
	if err := c.ShouldBindJSON(&loginRquest); err != nil {
		responses.ValidationErrorsJSON(c, err, "")
		return
	}
	var user models.User
	user, err := ac.userRepository.FindByField("Email", loginRquest.Email)
	if err != nil {
		ac.logger.Zap.Error("Failed to find user", err.Error())
		responses.ErrorJSON(c, http.StatusInternalServerError, gin.H{}, "An error occured")
		return
	}
	encryptedPassword := utils.Sha256Encrypt(loginRquest.Password)
	if user.Password == encryptedPassword {
		// login
		// token
		var accessToken string
		var refreshToken string
		accessToken, refreshToken, err = ac.authService.CreateTokens(int(user.Base.ID))
		if err != nil {
			ac.logger.Zap.Error("Failed generate jwt tokens", err.Error())
			responses.ErrorJSON(c, http.StatusInternalServerError, gin.H{}, "An error occured")
			return
		}
		var loginResult models.LoginResponse
		loginResult.AccessToken = accessToken
		loginResult.RefreshToken = refreshToken
		loginResult.User = models.UserResponse(user)

		responses.JSON(c, http.StatusOK, loginResult, "Hello "+user.FullName+" wellcome back")
		return
	} else {
		responses.ErrorJSON(c, http.StatusBadRequest, gin.H{}, "No user found with entered credentials")
		return
	}
}
