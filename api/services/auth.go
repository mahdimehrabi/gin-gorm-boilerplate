package services

import (
	"boilerplate/infrastructure"
	"boilerplate/models"
	"boilerplate/utils"
	"errors"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"gorm.io/datatypes"
)

// UserService -> struct
type AuthService struct {
	db infrastructure.Database
}

// NewAuthService -> creates a new AuthService
func NewAuthService(db infrastructure.Database) AuthService {
	return AuthService{
		db: db,
	}
}

func (as AuthService) CreateAccessToken(user models.User, exp int64, secret string, deviceToken string) (string, error) {
	var err error
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["userId"] = user.ID
	atClaims["password"] = user.Password
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
	atClaims["password"] = user.Password
	atClaims["exp"] = exp
	atClaims["deviceToken"] = deviceToken
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, err := at.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return token, nil
}

func (as AuthService) CreateTokens(user models.User, deviceToken string) (string, string, error) {
	var exp int64

	accessSecret := "access" + os.Getenv("Secret")
	exp = time.Now().Add(time.Hour * 2).Unix()
	accessToken, err := as.CreateAccessToken(user, exp, accessSecret, deviceToken)

	refreshSecret := "refresh" + os.Getenv("Secret")
	exp = time.Now().Add(time.Hour * 24 * 14).Unix()
	refreshToken, err := as.CreateRefreshToken(user, exp, refreshSecret, deviceToken)

	return accessToken, refreshToken, err
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
	if user.Devices != nil {
		devicesBytes := []byte(user.Devices.String())
		devices, err := utils.BytesJsonToMap(devicesBytes)
		if err != nil {
			return deviceToken, err
		}
		devices[deviceToken] = map[string]string{
			"ip":         c.ClientIP(),
			"city":       "alaki",
			"date":       strconv.Itoa(int(time.Now().Unix())),
			"deviceName": deviceName,
		}
		user.Devices = datatypes.JSON(utils.MapToJsonBytesBuffer(devices).String())
		as.db.DB.Save(&user)
		return deviceToken, nil
	}
	devices[deviceToken] = map[string]string{
		"ip":         c.ClientIP(),
		"city":       "alaki",
		"date":       strconv.Itoa(int(time.Now().Unix())),
		"deviceName": deviceName,
	}
	user.Devices = datatypes.JSON(utils.MapToJsonBytesBuffer(devices).String())
	as.db.DB.Save(&user)
	return deviceToken, nil
}
