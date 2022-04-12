package tests

import (
	"boilerplate/api/services"
	"boilerplate/infrastructure"
	"boilerplate/models"
	"boilerplate/utils"
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

func CreateUser(password string, db *gorm.DB, encryption infrastructure.Encryption) models.User {
	password = encryption.SaltAndSha256Encrypt(password)
	user := models.User{Email: password[len(password)-3:] + utils.GenerateRandomEmail(4),
		FirstName:        "mahdi",
		LastName:         "mehrabi",
		Password:         password,
		VerifyEmailToken: utils.GenerateRandomCode(40),
		VerifiedEmail:    true,
	}
	err := db.Create(&user).Error
	if err != nil {
		panic(err)
	}
	return user
}

func CreateAdmin(password string, db *gorm.DB, encryption infrastructure.Encryption) models.User {
	password = encryption.SaltAndSha256Encrypt(password)
	user := models.User{Email: password[len(password)-3:] + utils.GenerateRandomEmail(4),
		FirstName:        "mahdi",
		LastName:         "mehrabi",
		Password:         password,
		VerifyEmailToken: utils.GenerateRandomCode(40),
		VerifiedEmail:    true,
		IsAdmin:          true,
	}
	err := db.Create(&user).Error
	if err != nil {
		panic(err)
	}
	return user
}

func NewAuthenticatedRequest(as services.AuthService, db infrastructure.Database, user models.User, method string, url string, data *bytes.Buffer) (*http.Request, string, error) {
	req, err := http.NewRequest(method, url, data)
	if err != nil {
		return nil, "", err
	}
	deviceToken := utils.GenerateRandomCode(20)
	devices := make(map[string]interface{})
	devices[deviceToken] = map[string]string{
		"ip":         "1.1.1.1",
		"city":       "alaki",
		"date":       strconv.Itoa(int(time.Now().Unix())),
		"deviceName": "windows10-chrome",
	}
	user.Devices = datatypes.JSON(utils.MapToJsonBytesBuffer(devices).String())
	db.DB.Save(&user)
	tokensData, err := as.CreateTokens(user, deviceToken)
	req.Header.Add("Authorization", "Bearer "+tokensData["accessToken"])
	return req, deviceToken, err
}

func NewAuthenticatedRequestCustomDeviceToken(as services.AuthService, db infrastructure.Database, user models.User, method string, url string, data *bytes.Buffer, deviceToken string) (*http.Request, string, error) {
	req, err := http.NewRequest(method, url, data)
	if err != nil {
		return nil, "", err
	}
	devices := make(map[string]interface{})
	devices[deviceToken] = map[string]string{
		"ip":         "1.1.1.1",
		"city":       "alaki",
		"date":       strconv.Itoa(int(time.Now().Unix())),
		"deviceName": "windows10-chrome",
	}
	user.Devices = datatypes.JSON(utils.MapToJsonBytesBuffer(devices).String())
	db.DB.Save(&user)
	tokensData, err := as.CreateTokens(user, deviceToken)
	req.Header.Add("Authorization", "Bearer "+tokensData["accessToken"])
	return req, deviceToken, err
}

func ExtractResponseAsMap(w *httptest.ResponseRecorder) map[string]interface{} {
	response, err := utils.BytesJsonToMap(w.Body.Bytes())
	if err != nil {
		panic(err)
	}
	return response
}

func CreateFileRequestBody(key string, file *os.File) (*bytes.Buffer, *multipart.Writer, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(key, filepath.Base(file.Name()))

	if err != nil {
		return body, writer, err
	}

	io.Copy(part, file)
	writer.Close()

	return body, writer, nil
}
