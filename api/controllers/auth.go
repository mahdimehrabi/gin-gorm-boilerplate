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
	email          infrastructure.Email
}

func NewAuthController(logger infrastructure.Logger,
	env infrastructure.Env,
	encryption infrastructure.Encryption,
	userService services.UserService,
	authService services.AuthService,
	userRepository repositories.UserRepository,
	email infrastructure.Email,
) AuthController {
	return AuthController{
		logger:         logger,
		env:            env,
		encryption:     encryption,
		userService:    userService,
		authService:    authService,
		userRepository: userRepository,
		email:          email,
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
// @Success 200 {object} swagger.SuccessResponse
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
	user.VerifyEmailToken = utils.GenerateRandomCode(40)
	err := ac.userService.CreateUser(&user)
	if err != nil {
		ac.logger.Zap.Error("Failed to create registered user ", err.Error())
		responses.ErrorJSON(c, http.StatusInternalServerError, gin.H{}, "Sorry an error occoured in registering your account!")
		return
	}
	responses.JSON(c, http.StatusOK, gin.H{}, "Your account created successfuly, an verification link sent to your email use that to verify your account")

	go ac.sendRegisterationEmail(&user)
}

// @Summary login
// @Schemes
// @Description jwt login
// @Tags auth
// @Accept json
// @Produce json
// @Param email query string true "email"
// @Param deviceName query string true "send user operating system + browser name in this param"
// @Param password query string true "password"
// @Success 200 {object} swagger.LoginResponse
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
		if !user.VerifiedEmail {
			responses.ErrorJSON(c, http.StatusBadRequest, gin.H{}, "You must verify your email first")
			return
		}
		deviceToken, err := ac.authService.AddDevice(&user, c, loginRquest.DeviceName)
		if err != nil {
			ac.logger.Zap.Error("Failed to add device", err.Error())
			responses.ErrorJSON(c, http.StatusInternalServerError, gin.H{}, "An error occured")
			return
		}

		tokensData, err := ac.authService.CreateTokens(user, deviceToken)
		if err != nil {
			ac.logger.Zap.Error("Failed generate jwt tokens", err.Error())
			responses.ErrorJSON(c, http.StatusInternalServerError, gin.H{}, "An error occured")
			return
		}
		var loginResult models.LoginResponse
		loginResult.AccessToken = tokensData["accessToken"]
		loginResult.RefreshToken = tokensData["refreshToken"]
		loginResult.ExpRefreshToken = tokensData["expRefreshToken"]
		loginResult.ExpAccessToken = tokensData["expAccessToken"]
		loginResult.User = models.UserResponse(user)

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
// @Success 200 {object} swagger.SuccessResponse
// @failure 422 {object} swagger.FailedValidationResponse
// @failure 401 {object} swagger.FailedResponse
// @Router /auth/access-token-verify [post]
func (ac AuthController) AccessTokenVerify(c *gin.Context) {
	at := models.AccessTokenReq{}
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
// @Param refreshToken query string true "refreshToken"
// @Success 200 {object} swagger.SuccessVerifyAccessTokenResponse
// @failure 422 {object} swagger.FailedValidationResponse
// @failure 400 {object} swagger.FailedResponse
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

	uid, ok := atClaims["userId"].(float64)
	if !ok {
		responses.ErrorJSON(c, http.StatusBadRequest, gin.H{}, "Refresh token is not valid")
		return
	}
	userID := int(uid)

	deviceToken, ok := atClaims["deviceToken"].(string)
	if !ok {
		responses.ErrorJSON(c, http.StatusBadRequest, gin.H{}, "Refresh token is not valid")
		return
	}

	user, err := ac.authService.FindUserByIdDeviceToken(strconv.Itoa(userID), deviceToken)

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

	if valid {
		var exp int64
		accessSecret := "access" + ac.env.Secret
		exp = time.Now().Add(time.Hour * 2).Unix()
		accessToken, _ := ac.authService.CreateAccessToken(user, exp, accessSecret, deviceToken)
		responses.JSON(c, http.StatusOK, models.AccessTokenRes{AccessToken: accessToken, ExpAccessToken: strconv.Itoa(int(exp))}, "")
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
	user, err := ac.userRepository.GetAuthenticatedUser(c)
	if err != nil {
		ac.logger.Zap.Error("Failed to change password", err.Error())
		responses.ErrorJSON(c, http.StatusInternalServerError, gin.H{}, "Sorry an error occoured in changing password!")
		return
	}
	deviceToken := c.MustGet("deviceToken").(string)
	ac.authService.DeleteDevice(&user, deviceToken)

	responses.JSON(c, http.StatusOK, gin.H{}, "You logged out successfuly")
}

// @Summary verify-email
// @Schemes
// @Description verify-email
// @Tags auth
// @Accept json
// @Produce json
// @Param token query string true "token"
// @Success 200 {object} swagger.SuccessResponse
// @failure 422 {object} swagger.FailedValidationResponse
// @failure 400 {object} swagger.FailedResponse
// @Router /auth/verify-email [post]
func (ac AuthController) VerifyEmail(c *gin.Context) {
	tr := models.TokenRequest{}
	if err := c.ShouldBindJSON(&tr); err != nil {
		responses.ValidationErrorsJSON(c, err, "", map[string]string{})
		return
	}

	user, err := ac.userRepository.FindByField("verify_email_token", tr.Token)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		responses.ErrorJSON(c, http.StatusBadRequest, gin.H{}, "No user found with your token")
		return
	}

	if err != nil {
		ac.logger.Zap.Error("Failed to search for verify email token", err.Error())
		responses.ErrorJSON(c, http.StatusInternalServerError, gin.H{}, "Sorry an error occoured in changing password!")
		return
	}
	err = ac.userRepository.UpdateColumn(&user, "verified_email", true)
	if err != nil {
		ac.logger.Zap.Error("Failed to update verified_email column", err.Error())
		responses.ErrorJSON(c, http.StatusInternalServerError, gin.H{}, "Sorry an error occoured in changing password!")
		return
	}

	responses.JSON(c, http.StatusOK, gin.H{}, "Your email verified successfuly you can login now")
}

// @Summary forgot password
// @Schemes
// @Description forgot password
// @Tags auth
// @Accept json
// @Produce json
// @Param email query string true "unique email"
// @Success 200 {object} swagger.SuccessResponse
// @failure 422 {object} swagger.FailedValidationResponse
// @failure 400 {object} swagger.FailedResponse
// @Router /auth/forgot-password [post]
func (ac AuthController) ForgotPassword(c *gin.Context) {

	// Data Parse
	var er models.EmailRequest
	if err := c.ShouldBindJSON(&er); err != nil {
		responses.ValidationErrorsJSON(c, err, "", map[string]string{})
		return
	}

	user, err := ac.userRepository.FindByField("email", er.Email)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		responses.JSON(c, http.StatusOK, gin.H{}, "An email contain password recovery link sent to your email")
		return
	}
	if err != nil {
		ac.logger.Zap.Error("Failed to find user by email  ", err.Error())
		responses.ErrorJSON(c, http.StatusInternalServerError, gin.H{}, "Sorry an error occoured in registering your account!")
		return
	}
	if user.LastForgotEmailDate.After(time.Now().Add(time.Duration(-15) * time.Minute)) {
		responses.ErrorJSON(c, http.StatusBadRequest, gin.H{}, "You can request for recovering password after 15 minutes of your last reqeust")
		return
	}

	ac.userRepository.UpdateColumn(&user, "forgot_password_token", utils.GenerateRandomCode(40))

	responses.JSON(c, http.StatusOK, gin.H{}, "An email contain password recovery link sent to your email")
	go ac.sendForgotPassowrdEmail(&user)
}

func (ac AuthController) sendForgotPassowrdEmail(user *models.User) error {

	ch := make(chan error)
	htmlFile := ac.env.BasePath + "/vendors/templates/mail/auth/forgot.tmpl"

	data := map[string]string{
		"name": user.FirstName,
		"link": ac.env.SiteUrl + "/forgot-password?token=" + user.VerifyEmailToken,
	}
	go ac.email.SendEmail(ch, user.Email, "Recover password", htmlFile, data)
	err := <-ch
	if err != nil {
		ac.logger.Zap.Error(err)
		return err
	}
	err = ac.userRepository.UpdateColumn(user, "last_forgot_email_date", time.Now())
	if err != nil {
		ac.logger.Zap.Error(err)
		return err
	}
	return nil
}

// @Summary recover-password
// @Schemes
// @Description Let user change it password with forgot token
// @Tags auth
// @Accept json
// @Produce json
// @Param token query string true "token"
// @Param password query string true "password that have at least 8 length and contain an alphabet and number "
// @Param repeatPassword query string true "repeatPassword that have at least 8 length and contain an alphabet and number "
// @Success 200 {object} swagger.SuccessResponse
// @failure 422 {object} swagger.FailedValidationResponse
// @failure 404 {object} swagger.NotFoundResponse
// @Router /auth/recover-password [post]
func (ac AuthController) RecoverPassword(c *gin.Context) {
	// Data Parse
	var userData models.RecoverPassword
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

	user, err := ac.userRepository.FindByField("forgot_password_token", userData.Token)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		responses.JSON(c, http.StatusNotFound, gin.H{}, "404 not found!")
		return
	}
	if err != nil {
		ac.logger.Zap.Error("Failed to find user by forgot password token  ", err.Error())
		responses.ErrorJSON(c, http.StatusInternalServerError, gin.H{}, "Sorry an error occoured in registering your account!")
		return
	}
	encryptedPassword := ac.encryption.SaltAndSha256Encrypt(userData.Password)

	err = ac.userService.UpdatePassword(&user, encryptedPassword)
	if err != nil {
		ac.logger.Zap.Error("Failed to change password", err.Error())
		responses.ErrorJSON(c, http.StatusInternalServerError, gin.H{}, "Sorry an error occoured in changing password!")
		return
	}

	err = ac.userRepository.UpdateColumn(&user, "forgot_password_token", "")
	if err != nil {
		ac.logger.Zap.Error("Failed to change password", err.Error())
		responses.ErrorJSON(c, http.StatusInternalServerError, gin.H{}, "Sorry an error occoured in changing password!")
		return
	}

	responses.JSON(c, http.StatusOK, gin.H{}, "Your password changed successfuly")
}

// @Summary resend-verify-email
// @Schemes
// @Description resend-verify-email
// @Tags auth
// @Accept json
// @Produce json
// @Param email query string true "unique email"
// @Success 200 {object} swagger.SuccessResponse
// @failure 422 {object} swagger.FailedValidationResponse
// @Router /auth/resend-verify-email [post]
func (ac AuthController) ResendVerifyEmail(c *gin.Context) {

	// Data Parse
	var er models.EmailRequest
	if err := c.ShouldBindJSON(&er); err != nil {
		responses.ValidationErrorsJSON(c, err, "", map[string]string{})
		return
	}

	user, err := ac.userRepository.FindByField("email", er.Email)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		responses.JSON(c, http.StatusOK, gin.H{}, "Verification email sent!")
		return
	}
	if err != nil {
		ac.logger.Zap.Error("Failed to find user by email  ", err.Error())
		responses.ErrorJSON(c, http.StatusInternalServerError, gin.H{}, "Sorry an error occoured in registering your account!")
		return
	}
	if user.LastVerifyEmailDate.After(time.Now().Add(time.Duration(-15) * time.Minute)) {
		responses.ErrorJSON(c, http.StatusBadRequest, gin.H{}, "You can request for resending verification email after 15 minutes of your last reqeust")
		return
	}

	responses.JSON(c, http.StatusOK, gin.H{}, "Verification email sent!")
	go ac.sendRegisterationEmail(&user)
}

func (ac AuthController) sendRegisterationEmail(user *models.User) error {
	ch := make(chan error)
	htmlFile := ac.env.BasePath + "/vendors/templates/mail/auth/register.tmpl"

	data := map[string]string{
		"name": user.FirstName,
		"link": ac.env.SiteUrl + "/verify-email?token=" + user.VerifyEmailToken,
	}
	go ac.email.SendEmail(ch, user.Email, "Verify email", htmlFile, data)
	err := <-ch
	if err != nil {
		ac.logger.Zap.Error(err)
		return err
	}
	err = ac.userRepository.UpdateColumn(user, "last_verify_email_date", time.Now())
	if err != nil {
		ac.logger.Zap.Error(err)
		return err
	}
	return nil
}
