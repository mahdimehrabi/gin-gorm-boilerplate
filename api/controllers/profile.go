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

type ProfileController struct {
	logger         infrastructure.Logger
	env            infrastructure.Env
	encryption     infrastructure.Encryption
	userService    services.UserService
	authService    services.AuthService
	userRepository repositories.UserRepository
}

func NewProfileController(logger infrastructure.Logger,
	env infrastructure.Env,
	encryption infrastructure.Encryption,
	userService services.UserService,
	authService services.AuthService,
	userRepository repositories.UserRepository,
) ProfileController {
	return ProfileController{
		logger:         logger,
		env:            env,
		encryption:     encryption,
		userService:    userService,
		authService:    authService,
		userRepository: userRepository,
	}
}

// @Summary change-password
// @Schemes
// @Description Change Password , authentication required
// @Tags profile
// @Accept json
// @Produce json
// @Param password query string true "password that have at least 8 length and contain an alphabet and number "
// @Param repeatPassword query string true "repeatPassword that have at least 8 length and contain an alphabet and number "
// @Success 200 {object} swagger.SuccessResponse
// @failure 422 {object} swagger.FailedValidationResponse
// @failure 401 {object} swagger.UnauthenticatedResponse
// @Router /profile/change-password [post]
func (pc ProfileController) ChangePassword(c *gin.Context) {

	// Data Parse
	var userData models.ChangePassword
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

	encryptedPassword := pc.encryption.SaltAndSha256Encrypt(userData.Password)
	user, err := pc.userRepository.GetAuthenticatedUser(c)
	if err != nil {
		pc.logger.Zap.Error("Failed to change password", err.Error())
		responses.ErrorJSON(c, http.StatusInternalServerError, gin.H{}, "Sorry an error occoured in changing password!")
		return
	}
	err = pc.userService.UpdatePassword(&user, encryptedPassword)
	if err != nil {
		pc.logger.Zap.Error("Failed to change password", err.Error())
		responses.ErrorJSON(c, http.StatusInternalServerError, gin.H{}, "Sorry an error occoured in changing password!")
		return
	}

	responses.JSON(c, http.StatusOK, gin.H{}, "Your password changed successfuly , please login again !")
}

// @Summary devices
// @Schemes
// @Description return logged in devices in user's account , authentication required
// @Tags profile
// @Accept json
// @Produce json
// @Success 200 {object} swagger.DevicesResponse
// @failure 401 {object} swagger.UnauthenticatedResponse
// @Router /profile/devices [post]
func (pc ProfileController) LoggedInDevices(c *gin.Context) {

	user, err := pc.userRepository.GetAuthenticatedUser(c)
	if err != nil {
		pc.logger.Zap.Error("Failed to get authenticated user", err.Error())
		responses.ErrorJSON(c, http.StatusInternalServerError, gin.H{}, "Sorry an error occoured 😢")
		return
	}

	devices, err := pc.userRepository.GetLoggedInDevices(user)
	if err != nil {
		pc.logger.Zap.Error("Failed to get logged in devices", err.Error())
		responses.ErrorJSON(c, http.StatusInternalServerError, gin.H{}, "Sorry an error occoured 😢")
		return
	}

	responses.JSON(c, http.StatusOK, devices, "")
}

// @Summary terminate-device
// @Schemes
// @Description jwt terminate-device , atuhentication required
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {object} swagger.SuccessResponse
// @failure 422 {object} swagger.FailedValidationResponse
// @failure 401 {object} swagger.UnauthenticatedResponse
// @Router /auth/terminate-device [post]
func (pc ProfileController) TerminateDevice(c *gin.Context) {
	tr := models.TokenRequest{}
	if err := c.ShouldBindJSON(&tr); err != nil {
		responses.ValidationErrorsJSON(c, err, "", map[string]string{})
		return
	}

	user, err := pc.userRepository.GetAuthenticatedUser(c)
	if err != nil {
		pc.logger.Zap.Error("Failed to change password", err.Error())
		responses.ErrorJSON(c, http.StatusInternalServerError, gin.H{}, "Sorry an error occoured in changing password!")
		return
	}
	pc.authService.DeleteDevice(&user, tr.Token)

	responses.JSON(c, http.StatusOK, gin.H{}, "Selected device logged out succesfuly!")
}
