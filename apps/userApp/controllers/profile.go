package controllers

import (
	authServices "boilerplate/apps/authApp/services"
	"boilerplate/apps/userApp/repositories"
	"boilerplate/apps/userApp/services"
	"boilerplate/core/infrastructure"
	"boilerplate/core/models"
	"boilerplate/core/responses"
	"boilerplate/utils"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/datatypes"
)

type ProfileController struct {
	logger         infrastructure.Logger
	env            infrastructure.Env
	encryption     infrastructure.Encryption
	userService    services.UserService
	authService    authServices.AuthService
	userRepository repositories.UserRepository
	db             infrastructure.Database
}

func NewProfileController(logger infrastructure.Logger,
	env infrastructure.Env,
	encryption infrastructure.Encryption,
	userService services.UserService,
	authService authServices.AuthService,
	userRepository repositories.UserRepository,
	db infrastructure.Database,
) ProfileController {
	return ProfileController{
		logger:         logger,
		env:            env,
		encryption:     encryption,
		userService:    userService,
		authService:    authService,
		userRepository: userRepository,
		db:             db,
	}
}

// @Summary change-password
// @Schemes
// @Description Change Password , authentication required
// @Tags profile
// @Accept json
// @Produce json
// @Param currentPassword query string true "Current user password"
// @Param password query string true "password that have at least 8 length and contain an alphabet and number "
// @Param repeatPassword query string true "repeatPassword that have at least 8 length and contain an alphabet and number "
// @Success 200 {object} swagger.SuccessResponse
// @failure 422 {object} swagger.FailedValidationResponse
// @failure 401 {object} swagger.UnauthenticatedResponse
// @Router /profile/change-password [post]
func (pc ProfileController) ChangePassword(c *gin.Context) {

	// Data Parse
	var userData models.ChangeCurrentPassword
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
	currentEncryptedPassword := pc.encryption.SaltAndSha256Encrypt(userData.CurrentPassword)
	if currentEncryptedPassword != user.Password {
		fieldErrors := map[string]string{
			"currentPassword": "Current password is incorect",
		}
		responses.ManualValidationErrorsJSON(c, fieldErrors, "")
		return
	}
	if err != nil {
		pc.logger.Zap.Error("Failed to change password", err.Error())
		responses.ErrorJSON(c, http.StatusInternalServerError, gin.H{}, "Sorry an error occoured")
		return
	}
	err = pc.userService.UpdatePassword(&user, encryptedPassword)
	if err != nil {
		pc.logger.Zap.Error("Failed to change password", err.Error())
		responses.ErrorJSON(c, http.StatusInternalServerError, gin.H{}, "Sorry an error occoured")
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
// @Router /profile/devices [get]
func (pc ProfileController) LoggedInDevices(c *gin.Context) {

	user, err := pc.userRepository.GetAuthenticatedUser(c)
	if err != nil {
		pc.logger.Zap.Error("Failed to get authenticated user", err.Error())
		responses.ErrorJSON(c, http.StatusInternalServerError, gin.H{}, "Sorry an error occoured ????")
		return
	}

	devices, err := pc.userRepository.GetLoggedInDevices(user)
	if err != nil {
		pc.logger.Zap.Error("Failed to get logged in devices", err.Error())
		responses.ErrorJSON(c, http.StatusInternalServerError, gin.H{}, "Sorry an error occoured ????")
		return
	}

	responses.JSON(c, http.StatusOK, models.DeviceListResponse{
		Current: c.MustGet("deviceToken").(string),
		Devices: devices,
	}, "")
}

// @Summary terminate-device
// @Schemes
// @Description jwt terminate-device , atuhentication required
// @Tags profile
// @Accept json
// @Produce json
// @Param token query string true "token of the device that we want to remove"
// @Success 200 {object} swagger.DevicesResponse
// @failure 422 {object} swagger.FailedValidationResponse
// @failure 404 {object} swagger.NotFoundResponse
// @failure 401 {object} swagger.UnauthenticatedResponse
// @Router /profile/terminate-device [post]
func (pc ProfileController) TerminateDevice(c *gin.Context) {
	tr := models.TokenRequestNoLimit{}
	if err := c.ShouldBindJSON(&tr); err != nil {
		responses.ValidationErrorsJSON(c, err, "", map[string]string{})
		return
	}

	user, err := pc.userRepository.GetAuthenticatedUser(c)
	if err != nil {
		pc.logger.Zap.Error("Failed to change password", err.Error())
		responses.ErrorJSON(c, http.StatusInternalServerError, gin.H{}, "Sorry an error occoured")
		return
	}

	devicesBytes := []byte(user.Devices.String())
	devices, err := utils.BytesJsonToMap(devicesBytes)
	if err != nil {
		pc.logger.Zap.Error("Failed to terminate device", err.Error())
		responses.ErrorJSON(c, http.StatusInternalServerError, gin.H{}, "Sorry an error occoured ????")
		return
	}
	if devices[tr.Token] == nil {
		responses.ErrorJSON(c, http.StatusNotFound, gin.H{}, "This device already logged out or not exist at all")
		return
	}

	pc.authService.DeleteDevice(&user, tr.Token)

	resDevices, err := pc.userRepository.GetLoggedInDevices(user)
	if err != nil {
		pc.logger.Zap.Error("Failed to terminate device", err.Error())
		responses.ErrorJSON(c, http.StatusInternalServerError, gin.H{}, "Sorry an error occoured ????")
		return
	}
	responses.JSON(c, http.StatusOK, models.DeviceListResponse{
		Current: c.MustGet("deviceToken").(string),
		Devices: resDevices,
	}, "Selected device logged out succesfuly!")
}

// @Summary terminate-devices-except-me
// @Schemes
// @Description terminate all devices execpt current device , atuhentication required
// @Tags profile
// @Accept json
// @Produce json
// @Success 200 {object} swagger.DevicesResponse
// @failure 422 {object} swagger.FailedValidationResponse
// @failure 404 {object} swagger.NotFoundResponse
// @failure 401 {object} swagger.UnauthenticatedResponse
// @Router /profile/terminate-devices-except-me [post]
func (pc ProfileController) TerminateDevicesExceptMe(c *gin.Context) {
	user, err := pc.userRepository.GetAuthenticatedUser(c)
	if err != nil {
		pc.logger.Zap.Error("Failed to change password", err.Error())
		responses.ErrorJSON(c, http.StatusInternalServerError, gin.H{}, "Sorry an error occoured")
		return
	}
	token := c.MustGet("deviceToken").(string)

	devicesBytes := []byte(user.Devices.String())
	devices, err := utils.BytesJsonToMap(devicesBytes)
	if err != nil {
		pc.logger.Zap.Error("Failed to terminate all devices except me", err.Error())
		responses.ErrorJSON(c, http.StatusInternalServerError, gin.H{}, "Sorry an error occoured ????")
		return
	}
	currentDevice := devices[token]
	devices = map[string]interface{}{
		token: currentDevice,
	}
	user.Devices = datatypes.JSON(utils.MapToJsonBytesBuffer(devices).String())
	err = pc.db.DB.Save(&user).Error
	if err != nil {
		pc.logger.Zap.Error("Failed to terminate all devices except me", err.Error())
		responses.ErrorJSON(c, http.StatusInternalServerError, gin.H{}, "Sorry an error occoured ????")
		return
	}

	resDevices, err := pc.userRepository.GetLoggedInDevices(user)
	if err != nil {
		pc.logger.Zap.Error("Failed to terminate all devices except me", err.Error())
		responses.ErrorJSON(c, http.StatusInternalServerError, gin.H{}, "Sorry an error occoured ????")
		return
	}
	responses.JSON(c, http.StatusOK, models.DeviceListResponse{
		Current: c.MustGet("deviceToken").(string),
		Devices: resDevices,
	}, "All devices except current device logged out successfuly")
}

// @Summary upload-profile-picture
// @Schemes
// @Description Upload profile picture , authentication required
// @Tags profile
// @Accept json
// @Produce json
// @Param picture formData string true "file of image"
// @Success 200 {object} swagger.SuccessResponse
// @failure 422 {object} swagger.FailedValidationResponse
// @failure 401 {object} swagger.UnauthenticatedResponse
// @Router /profile/upload-profile-picture [post]
func (pc ProfileController) UploadProfilePicture(c *gin.Context) {
	user, err := pc.userRepository.GetAuthenticatedUser(c)
	if err != nil {
		pc.logger.Zap.Error("Failed to get user ", err.Error())
		responses.ErrorJSON(c, http.StatusInternalServerError, gin.H{}, "Sorry an error occoured")
		return
	}

	directoryPath := pc.env.BasePath + "/media/users/" + strconv.Itoa(int(user.ID))
	os.MkdirAll(directoryPath, os.ModePerm)
	uploadPath := directoryPath + "/picture"
	res, filePath, err := utils.UploadFile(uploadPath, c, "picture", []string{"jpg", "jpeg"})
	if err != nil {
		pc.logger.Zap.Error("Failed to save uploaded picture profile ", err.Error())
		responses.ErrorJSON(c, http.StatusInternalServerError, gin.H{}, "Sorry an error occoured")
		return
	}
	if !res {
		return
	}
	err = utils.ResizeImage(filePath, filePath, 200, 200)
	if err != nil {
		pc.logger.Zap.Error("Failed to resize user image", err.Error())
		fieldErrors := make(map[string]string, 0)
		fieldErrors["picture"] = "It seems that your image file is not standard please use another image"
		responses.ValidationErrorsJSON(c, err, "", fieldErrors)
		return
	}

	path := filePath[strings.Index(filePath, "/media"):]
	user.Picture = path
	err = pc.db.DB.Save(&user).Error
	if err != nil {
		pc.logger.Zap.Error("Failed to save user picture path in db ", err.Error())
		responses.ErrorJSON(c, http.StatusInternalServerError, gin.H{}, "Sorry an error occoured")
		return
	}

	responses.JSON(c, http.StatusOK, user.ToResponse(pc.env), "Your password changed successfuly , please login again !")
}
