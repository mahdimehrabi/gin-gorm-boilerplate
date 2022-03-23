package tests

import (
	"boilerplate/api/services"
	"boilerplate/infrastructure"
	"boilerplate/models"
	"boilerplate/utils"
	"bytes"
	"net/http"
	"net/http/httptest"

	"gorm.io/gorm"
)

func CreateUser(password string, db *gorm.DB, encryption infrastructure.Encryption) models.User {
	password = encryption.SaltAndSha256Encrypt(password)
	user := models.User{Email: utils.GenerateRandomEmail(7), FirstName: "mahdi", LastName: "mehrabi", Password: password}
	err := db.Create(&user).Error
	if err != nil {
		panic(err)
	}
	return user
}

func NewAuthenticatedRequest(as services.AuthService, user models.User, method string, url string, data *bytes.Buffer) (*http.Request, error) {
	req, err := http.NewRequest(method, url, data)
	if err != nil {
		return nil, err
	}
	accessToken, _, err := as.CreateTokens(user)
	req.Header.Add("Authorization", "Bearer "+accessToken)
	return req, err
}

func ExtractResponseAsMap(w *httptest.ResponseRecorder) map[string]interface{} {
	response, err := utils.BytesJsonToMap(w.Body.Bytes())
	if err != nil {
		panic(err)
	}
	return response
}
