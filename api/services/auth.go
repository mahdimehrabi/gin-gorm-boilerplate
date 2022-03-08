package services

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

// UserService -> struct
type AuthService struct {
}

// NewAuthService -> creates a new AuthService
func NewAuthService() AuthService {
	return AuthService{}
}

func (as AuthService) CreateToken(userID int, exp int64, secret string) (string, error) {
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

func (as AuthService) CreateTokens(userID int) (string, string, error) {
	var exp int64

	os.Setenv("SECRET", "you need to set secret")
	accessSecret := "access" + os.Getenv("SECRET")
	exp = time.Now().Add(time.Hour * 2).Unix()
	accessToken, err := as.CreateToken(userID, exp, accessSecret)

	refreshSecret := "refresh" + os.Getenv("SECRET")
	exp = time.Now().Add(time.Hour * 24 * 14).Unix()
	refreshToken, err := as.CreateToken(userID, exp, refreshSecret)

	return accessToken, refreshToken, err
}
