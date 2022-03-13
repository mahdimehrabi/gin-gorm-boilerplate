package controllers

import (
	"boilerplate/api/repositories"
	"boilerplate/api/responses"
	"boilerplate/api/services"
	"boilerplate/infrastructure"
	"boilerplate/models"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"gorm.io/gorm"
)

type AuthController struct {
	logger         infrastructure.Logger
	env            infrastructure.Env
	encryption     infrastructure.Encryption
	userService    services.UserService
	authService    services.AuthService
	userRepository repositories.UserRepository
}

func NewAuthController(logger infrastructure.Logger,
	env infrastructure.Env,
	encryption infrastructure.Encryption,
	userService services.UserService,
	authService services.AuthService,
	userRepository repositories.UserRepository,
) AuthController {
	return AuthController{
		logger:         logger,
		env:            env,
		encryption:     encryption,
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
	encryptedPassword := ac.encryption.SaltAndSha256Encrypt(userData.Password)
	user.Password = encryptedPassword
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
	if errors.Is(err, gorm.ErrRecordNotFound) {
		responses.ErrorJSON(c, http.StatusBadRequest, gin.H{}, "No user found with entered credentials")
		return
	}
	if err != nil {
		ac.logger.Zap.Error("Error to find user", err.Error())
		responses.ErrorJSON(c, http.StatusInternalServerError, gin.H{}, "An error occured")
		return
	}
	encryptedPassword := ac.encryption.SaltAndSha256Encrypt(loginRquest.Password)
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

		//make must_logout false
		if err = ac.userRepository.UpdateColumn(&user, "must_logout", false); err != nil {
			ac.logger.Zap.Error("Failed make must_logout false", err.Error())
			responses.ErrorJSON(c, http.StatusInternalServerError, gin.H{}, "An error occured")
			return
		}
		responses.JSON(c, http.StatusOK, loginResult, "Hello "+user.FullName+" wellcome back")
		return
	} else {
		responses.ErrorJSON(c, http.StatusBadRequest, gin.H{}, "No user found with entered credentials")
		return
	}
}

type accessTokenReqRes struct {
	AccessToken string `json:"accessToken" binding:"required"`
}

func (ac AuthController) AccessTokenVerify(c *gin.Context) {
	at := accessTokenReqRes{}
	if err := c.ShouldBindJSON(&at); err != nil {
		responses.ValidationErrorsJSON(c, err, "")
		return
	}

	accessToken := at.AccessToken
	accessSecret := "access" + ac.env.Secret
	valid, _, err := services.DecodeToken(accessToken, accessSecret)
	if err != nil {
		responses.ErrorJSON(c, http.StatusBadRequest, gin.H{}, "Access token is not valid")
		return
	}

	if valid {
		responses.JSON(c, http.StatusOK, gin.H{}, "Access token is valid")
		return
	} else {
		responses.ErrorJSON(c, http.StatusBadRequest, gin.H{}, "Access token is not valid")
		return
	}
}

type refreshTokenRequest struct {
	RefreshToken string `json:"refreshToken" binding:"required"`
}

func (ac AuthController) RenewToken(c *gin.Context) {
	rtr := refreshTokenRequest{}
	if err := c.ShouldBindJSON(&rtr); err != nil {
		responses.ValidationErrorsJSON(c, err, "")
		return
	}

	//Parse and extract claims
	refreshToken := rtr.RefreshToken
	var valid bool
	var atClaims jwt.MapClaims
	refreshSecret := "refresh" + ac.env.Secret
	valid, atClaims, err := services.DecodeToken(refreshToken, refreshSecret)
	if err != nil {
		responses.ErrorJSON(c, http.StatusBadRequest, gin.H{}, "Refresh token is not valid")
		return
	}

	userID := int(atClaims["userId"].(float64))
	user, err := ac.userRepository.FindByField("id", strconv.Itoa(userID))

	//don't allow deleted user renew access token
	if errors.Is(err, gorm.ErrRecordNotFound) {
		responses.ErrorJSON(c, http.StatusBadRequest, gin.H{}, "Refresh token is not valid")
		return
	}

	if err != nil {
		ac.logger.Zap.Error("Error in finding user:", err)
		responses.ErrorJSON(c, http.StatusBadRequest, gin.H{}, "Refresh token is not valid")
		return
	}

	//if user must_logout field be true it can't refresh its token
	if user.MustLogout {
		responses.ErrorJSON(c, http.StatusBadRequest, gin.H{}, "Refresh token is not valid")
		return
	}

	if valid {
		var exp int64
		accessSecret := "refresh" + ac.env.Secret
		exp = time.Now().Add(time.Hour * 2).Unix()
		accessToken, _ := ac.authService.CreateToken(int(userID), exp, accessSecret)
		responses.JSON(c, http.StatusOK, accessTokenReqRes{AccessToken: accessToken}, "")
		return
	} else {
		responses.ErrorJSON(c, http.StatusBadRequest, gin.H{}, "Refresh token is not valid")
		return
	}
}
