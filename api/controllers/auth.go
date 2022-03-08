package controllers

import (
	"boilerplate/api/services"
	"boilerplate/infrastructure"
	"boilerplate/models"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func CreateToken(userID int, exp int64, secret string) (string, error) {
	var err error
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["user_id"] = userID
	atClaims["exp"] = exp
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, err := at.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return token, nil
}
func CreateTokens(userID int) (string, string, error) {
	var exp int64

	os.Setenv("SECRET", "you need to set secret")
	accessSecret := "access" + os.Getenv("SECRET")
	exp = time.Now().Add(time.Hour * 2).Unix()
	accessToken, err := CreateToken(userID, exp, accessSecret)

	refreshSecret := "refresh" + os.Getenv("SECRET")
	exp = time.Now().Add(time.Hour * 24 * 14).Unix()
	refreshToken, err := CreateToken(userID, exp, refreshSecret)

	return accessToken, refreshToken, err
}

//TODO:remove
func NullStr2Str(str string) (nullStr sql.NullString) {
	if str == "" {
		nullStr.String = ""
		nullStr.Valid = false
	} else {
		nullStr.String = str
		nullStr.Valid = true
	}
	return nullStr
}

//TODO:move to utils
func EncodePassword(password string) string {
	hash := sha256.New()
	hash.Write([]byte(password))
	md := hash.Sum(nil)
	encodedPassword := hex.EncodeToString(md)
	return encodedPassword
}

type AuthController struct {
	logger      infrastructure.Logger
	env         infrastructure.Env
	userService services.UserService
}

func NewAuthController(logger infrastructure.Logger,
	env infrastructure.Env,
	userService services.UserService,
) AuthController {
	return AuthController{
		logger:      logger,
		env:         env,
		userService: userService,
	}
}

func (ac AuthController) Register(c *gin.Context) {

	// Data Parse
	var userData models.AddUserData
	err := c.ShouldBindJSON(&userData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	encodedPassword := EncodePassword(userData.Password)
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
	accessToken, refreshToken, err := CreateTokens(int(user.Base.ID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	var loginResult models.LoginResult
	loginResult.UserID = int(user.Base.ID)
	loginResult.AccessToken = accessToken
	loginResult.RefreshToken = refreshToken

	c.JSON(http.StatusOK, loginResult)
}
