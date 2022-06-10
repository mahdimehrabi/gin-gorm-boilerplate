package services

import (
	"boilerplate/apps/userApp/repositories"
	"boilerplate/core/infrastructure"
	"boilerplate/core/models"
	"boilerplate/utils"
	"errors"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/ip2location/ip2location-go/v9"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// UserService -> struct
type AuthService struct {
	db             infrastructure.Database
	env            infrastructure.Env
	email          infrastructure.Email
	logger         infrastructure.Logger
	userRepository repositories.UserRepository
}

// NewAuthService -> creates a new AuthService
func NewAuthService(
	db infrastructure.Database,
	env infrastructure.Env,
	email infrastructure.Email,
	logger infrastructure.Logger,
	userRepository repositories.UserRepository) AuthService {
	return AuthService{
		db:             db,
		env:            env,
		email:          email,
		logger:         logger,
		userRepository: userRepository,
	}
}

func (as AuthService) CreateAccessToken(user models.User, exp int64, secret string, deviceToken string) (string, error) {
	var err error
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["userId"] = user.ID
	atClaims["exp"] = exp
	atClaims["deviceToken"] = deviceToken
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, err := at.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return token, nil
}

func (as AuthService) CreateRefreshToken(user models.User, exp int64, secret string, deviceToken string) (string, error) {
	var err error
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["userId"] = user.ID
	atClaims["exp"] = exp
	atClaims["deviceToken"] = deviceToken
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, err := at.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return token, nil
}

func (as AuthService) CreateTokens(user models.User, deviceToken string) (map[string]string, error) {

	accessSecret := "access" + os.Getenv("Secret")
	expAccessToken := time.Now().Add(time.Minute * 15).Unix()
	accessToken, err := as.CreateAccessToken(user, expAccessToken, accessSecret, deviceToken)

	refreshSecret := "refresh" + os.Getenv("Secret")
	expRefreshToken := time.Now().Add(time.Hour * 24 * 60).Unix()
	refreshToken, err := as.CreateRefreshToken(user, expRefreshToken, refreshSecret, deviceToken)

	return map[string]string{
		"refreshToken":    refreshToken,
		"accessToken":     accessToken,
		"expRefreshToken": strconv.Itoa(int(expRefreshToken)),
		"expAccessToken":  strconv.Itoa(int(expAccessToken)),
	}, err
}

func DecodeToken(tokenString string, secret string) (bool, jwt.MapClaims, error) {

	Claims := jwt.MapClaims{}

	key := func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			ErrUnexpectedSigningMethod := errors.New("unexpected signing method")
			return nil, ErrUnexpectedSigningMethod
		}
		return []byte(secret), nil
	}

	token, err := jwt.ParseWithClaims(tokenString, Claims, key)
	var valid bool
	if token == nil {
		valid = false
	} else {
		valid = token.Valid
	}
	return valid, Claims, err
}

//Add device information on login and set deviceToken that used as jwt claim in refreshToken
func (as AuthService) AddDevice(user *models.User, c *gin.Context, deviceName string) (string, error) {
	deviceToken := utils.GenerateRandomCode(20)
	devices := make(map[string]interface{})
	ip := c.ClientIP()
	city, _ := as.GetCityByIp(ip)
	if user.Devices != nil {
		devicesBytes := []byte(user.Devices.String())
		devices, err := utils.BytesJsonToMap(devicesBytes)
		if err != nil {
			return deviceToken, err
		}
		devices[deviceToken] = map[string]string{
			"ip":         ip,
			"city":       city,
			"date":       strconv.Itoa(int(time.Now().Unix())),
			"deviceName": deviceName,
		}
		user.Devices = datatypes.JSON(utils.MapToJsonBytesBuffer(devices).String())
		as.db.DB.Save(&user)
		return deviceToken, nil
	}
	devices[deviceToken] = map[string]string{
		"ip":         ip,
		"city":       city,
		"date":       strconv.Itoa(int(time.Now().Unix())),
		"deviceName": deviceName,
	}
	user.Devices = datatypes.JSON(utils.MapToJsonBytesBuffer(devices).String())
	as.db.DB.Save(&user)
	return deviceToken, nil
}

//delete specefic deviceToken so he/she will logedout with that device
func (as AuthService) DeleteDevice(user *models.User, deviceToken string) (bool, error) {
	if user.Devices == nil {
		return false, errors.New("user devices is nil")
	}
	devicesBytes := []byte(user.Devices.String())
	devices, err := utils.BytesJsonToMap(devicesBytes)
	if err != nil {
		return false, err
	}
	delete(devices, deviceToken)
	user.Devices = datatypes.JSON(utils.MapToJsonBytesBuffer(devices).String())
	err = as.db.DB.Save(&user).Error
	if err != nil {
		return false, err
	}
	return true, nil
}

//get city by ip , if cant identified city it will return unknown string
func (as AuthService) GetCityByIp(ip string) (string, error) {

	filePath := as.env.BasePath + "/vendors/IP2LOCATION-LITE-DB3.BIN"
	db, err := ip2location.OpenDB(filePath)

	if err != nil {
		return "unknown", err
	}
	results, err := db.Get_city(ip)

	if err != nil {
		return "unknown", err
	}

	if results.City == "Invalid IP address." || results.City == "-" {
		return "unknown", nil
	}

	return results.City, nil
}

//find user by id and deviceToken
func (as AuthService) FindUserByIdDeviceToken(id string, deviceToken string) (models.User, error) {
	var user models.User
	err := as.db.DB.Model(&models.User{}).Where("id", id).
		Find(&user, datatypes.JSONQuery("devices").HasKey(deviceToken)).Error
	if user.ID == 0 {
		return user, gorm.ErrRecordNotFound
	}
	return user, err
}

func (as AuthService) SendRegisterationEmail(user *models.User) error {
	ch := make(chan error)
	htmlFile := as.env.BasePath + "/vendors/templates/mail/auth/register.tmpl"

	data := map[string]string{
		"name": user.FirstName,
		"link": as.env.SiteUrl + "/verify-email?token=" + user.VerifyEmailToken,
	}
	go as.email.SendEmail(ch, user.Email, "Verify email", htmlFile, data)
	err := <-ch
	if err != nil {
		as.logger.Zap.Error(err)
		return err
	}
	err = as.userRepository.UpdateColumn(user, "last_verify_email_date", time.Now())
	if err != nil {
		as.logger.Zap.Error(err)
		return err
	}
	return nil
}
