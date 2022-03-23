package controllers

import (
	"boilerplate/api/repositories"
	"boilerplate/api/responses"
	"boilerplate/api/services"
	"boilerplate/infrastructure"
	"boilerplate/models"
	"boilerplate/utils"
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

// @BasePath /api/auth

// @Summary register
// @Schemes
// @Description jwt register
// @Tags auth
// @Accept json
// @Produce json
// @Param email query string true "unique email"
// @Param password query string true "password that have at least 8 length and contain an alphabet and number "
// @Param repeatPassword query string true "repeatPassword that have at least 8 length and contain an alphabet and number "
// @Param firstName query string true "firstName"
// @Param lastName query string true "lastName"
// @Success 200 {object} swagger.RegisterLoginResponse
// @failure 422 {object} swagger.FailedValidationResponse
// @Router /auth/register [post]
func (ac AuthController) Register(c *gin.Context) {

	// Data Parse
	var userData models.Register
	if err := c.ShouldBindJSON(&userData); err != nil {
		fieldErrors := make(map[string]string, 0)
		if !utils.IsGoodPassword(userData.Password) {
			fieldErrors["password"] = "Password must contain at least one alphabet and one number and its length must be 8 characters or more"

		}
		responses.ValidationErrorsJSON(c, err, "", fieldErrors)
		return
	}
	if !utils.IsGoodPassword(userData.Password) {
		fieldErrors := map[string]string{
			"password": "Password must contain at least one alphabet and one number and its length must be 8 characters or more",
		}
		responses.ManualValidationErrorsJSON(c, fieldErrors, "")
		return
	}
	var user models.User
	encryptedPassword := ac.encryption.SaltAndSha256Encrypt(userData.Password)
	user.Password = encryptedPassword

	user.FirstName = userData.FirstName
	user.LastName = userData.LastName
	user.Email = userData.Email
	err := ac.userService.CreateUser(&user)
	if err != nil {
		ac.logger.Zap.Error("Failed to create registered user ", err.Error())
		responses.ErrorJSON(c, http.StatusInternalServerError, gin.H{}, "Sorry an error occoured in registering your account!")
		return
	}

	// login
	// token
	accessToken, refreshToken, err := ac.authService.CreateTokens(user)
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

// @Summary login
// @Schemes
// @Description jwt login
// @Tags auth
// @Accept json
// @Produce json
// @Param email query string true "email"
// @Param password query string true "password"
// @Success 200 {object} swagger.RegisterLoginResponse
// @failure 422 {object} swagger.FailedValidationResponse
// @failure 401 {object} swagger.FailedLoginResponse
// @Router /auth/login [post]
func (ac AuthController) Login(c *gin.Context) {
	// Data Parse
	var loginRquest models.LoginRequest
	if err := c.ShouldBindJSON(&loginRquest); err != nil {
		responses.ValidationErrorsJSON(c, err, "", map[string]string{})
		return
	}
	var user models.User
	user, err := ac.userRepository.FindByField("Email", loginRquest.Email)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		responses.ErrorJSON(c, http.StatusUnauthorized, gin.H{}, "No user found with entered credentials")
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
		accessToken, refreshToken, err = ac.authService.CreateTokens(user)
		if err != nil {
			ac.logger.Zap.Error("Failed generate jwt tokens", err.Error())
			responses.ErrorJSON(c, http.StatusInternalServerError, gin.H{}, "An error occured")
			return
		}
		var loginResult models.LoginResponse
		loginResult.AccessToken = accessToken
		loginResult.RefreshToken = refreshToken
		loginResult.User = models.UserResponse(user)
		ac.userRepository.AddDevice(&user, c, "")

		responses.JSON(c, http.StatusOK, loginResult, "Hello "+user.FirstName+" wellcome back")
		return
	} else {
		responses.ErrorJSON(c, http.StatusUnauthorized, gin.H{}, "No user found with entered credentials")
		return
	}
}

// @Summary access token verify
// @Schemes
// @Description jwt access token verify
// @Tags auth
// @Accept json
// @Produce json
// @Param accessToken query string true "accessToken"
// @Success 200 {object} swagger.SuccessVerifyAccessTokenResponse
// @failure 422 {object} swagger.FailedValidationResponse
// @failure 401 {object} swagger.FailedResponse
// @Router /auth/access-token-verify [post]
func (ac AuthController) AccessTokenVerify(c *gin.Context) {
	at := models.AccessTokenReqRes{}
	if err := c.ShouldBindJSON(&at); err != nil {
		responses.ValidationErrorsJSON(c, err, "", map[string]string{})
		return
	}

	accessToken := at.AccessToken
	accessSecret := "access" + ac.env.Secret
	valid, _, err := services.DecodeToken(accessToken, accessSecret)
	if err != nil {
		responses.ErrorJSON(c, http.StatusUnauthorized, gin.H{}, "Access token is not valid")
		return
	}

	if valid {
		responses.JSON(c, http.StatusOK, gin.H{}, "Access token is valid")
		return
	} else {
		responses.ErrorJSON(c, http.StatusUnauthorized, gin.H{}, "Access token is not valid")
		return
	}
}

// @Summary renew access token
// @Schemes
// @Description jwt renew access token
// @Tags auth
// @Accept json
// @Produce json
// @Param refreshToken query string true "accessToken"
// @Success 200 {object} swagger.SuccessVerifyAccessTokenResponse
// @failure 422 {object} swagger.FailedValidationResponse
// @failure 401 {object} swagger.FailedResponse
// @Router /auth/renew-access-token [post]
func (ac AuthController) RenewToken(c *gin.Context) {
	rtr := models.RefreshTokenRequest{}
	if err := c.ShouldBindJSON(&rtr); err != nil {
		responses.ValidationErrorsJSON(c, err, "", map[string]string{})
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
	password := atClaims["password"].(string)
	user, err := ac.userRepository.FindByField("id", strconv.Itoa(userID))

	//don't allow deleted user renew access token
	if errors.Is(err, gorm.ErrRecordNotFound) {
		responses.ErrorJSON(c, http.StatusBadRequest, gin.H{}, "Refresh token is not valid")
		return
	}
	if user.Password != password {
		responses.ErrorJSON(c, http.StatusBadRequest, gin.H{}, "Your password has changed !")
		return
	}

	if err != nil {
		ac.logger.Zap.Error("Error in finding user:", err)
		responses.ErrorJSON(c, http.StatusBadRequest, gin.H{}, "Refresh token is not valid")
		return
	}

	if valid {
		var exp int64
		accessSecret := "refresh" + ac.env.Secret
		exp = time.Now().Add(time.Hour * 2).Unix()
		accessToken, _ := ac.authService.CreateToken(user, exp, accessSecret)
		responses.JSON(c, http.StatusOK, models.AccessTokenReqRes{AccessToken: accessToken}, "")
		return
	} else {
		responses.ErrorJSON(c, http.StatusBadRequest, gin.H{}, "Refresh token is not valid")
		return
	}
}

// @Summary logout
// @Schemes
// @Description jwt logout , atuhentication required
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {object} swagger.SuccessResponse
// @failure 401 {object} swagger.UnauthenticatedResponse
// @Router /auth/logout [post]
func (ac AuthController) Logout(c *gin.Context) {
	_, err := ac.userRepository.GetAuthenticatedUser(c)
	if err != nil {
		ac.logger.Zap.Error("Failed to change password", err.Error())
		responses.ErrorJSON(c, http.StatusInternalServerError, gin.H{}, "Sorry an error occoured in changing password!")
		return
	}

	responses.JSON(c, http.StatusOK, gin.H{}, "You logged out successfuly")
}
