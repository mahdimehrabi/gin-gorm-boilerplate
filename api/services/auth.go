package services

import (
	"boilerplate/models"
	"errors"
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

func (as AuthService) CreateToken(user models.User, exp int64, secret string) (string, error) {
	var err error
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["userId"] = user.ID
	atClaims["password"] = user.Password
	atClaims["exp"] = exp
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, err := at.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return token, nil
}

func (as AuthService) CreateTokens(user models.User) (string, string, error) {
	var exp int64

	accessSecret := "access" + os.Getenv("Secret")
	exp = time.Now().Add(time.Hour * 2).Unix()
	accessToken, err := as.CreateToken(user, exp, accessSecret)

	refreshSecret := "refresh" + os.Getenv("Secret")
	exp = time.Now().Add(time.Hour * 24 * 14).Unix()
	refreshToken, err := as.CreateToken(user, exp, refreshSecret)

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
