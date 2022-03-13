package tests

import (
	"boilerplate/models"
	"boilerplate/utils"
	"bytes"
	"encoding/json"
	"net/http/httptest"

	"gorm.io/gorm"
)

func MapToJsonBytesBuffer(mp map[string]interface{}) *bytes.Buffer {
	j, err := json.Marshal(mp)
	if err != nil {
		panic(err)
	}
	return bytes.NewBuffer(j)
}

func CreateUser(db *gorm.DB) models.User {
	password := utils.Sha256Encrypt("m12345678")
	user := models.User{Email: "mahdi@gmail.com", FullName: "mahdi mehrabi", Password: password}
	err := db.Create(&user).Error
	if err != nil {
		panic(err)
	}
	return user
}

func ExtractResponseAsMap(w *httptest.ResponseRecorder) map[string]interface{} {
	response := map[string]interface{}{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		panic(err)
	}
	return response
}
