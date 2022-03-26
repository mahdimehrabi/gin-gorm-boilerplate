package tests

import (
	"boilerplate/utils"
	"net/http"
	"net/http/httptest"
)

func (suite TestSuiteEnv) TestChangePassword() {
	router := suite.router.Gin
	db := suite.database.DB
	a := suite.Assert()
	user := CreateUser("m12345678", db, suite.encryption)

	data := map[string]interface{}{
		"password":       "m987654321",
		"repeatPassword": "m987654321",
	}
	w := httptest.NewRecorder()

	req, _, _ := NewAuthenticatedRequest(suite.authService, suite.database, user, "POST", "/api/profile/change-password", utils.MapToJsonBytesBuffer(data))

	router.ServeHTTP(w, req)
	a.Equal(http.StatusOK, w.Code, "Status code problem")
	suite.database.DB.Find(&user)
	encryptedPassword := suite.encryption.SaltAndSha256Encrypt("m987654321")
	a.Equal(encryptedPassword, user.Password, "encrypt password problem")
	devicesBytes := []byte(user.Devices.String())
	devices, _ := utils.BytesJsonToMap(devicesBytes)
	a.Equal(len(devices), 0, "devices is not empty")

	//test with weak password
	user = CreateUser("m12345678", db, suite.encryption)
	data = map[string]interface{}{
		"password":       "12345678",
		"repeatPassword": "12345678",
	}
	w = httptest.NewRecorder()
	req, _, _ = NewAuthenticatedRequest(suite.authService, suite.database, user, "POST", "/api/profile/change-password", utils.MapToJsonBytesBuffer(data))
	router.ServeHTTP(w, req)
	a.Equal(http.StatusUnprocessableEntity, w.Code, "Status code problem")
}

func (suite TestSuiteEnv) TestGetLoggedInDevices() {
	router := suite.router.Gin
	db := suite.database.DB
	a := suite.Assert()
	user := CreateUser("m12345678", db, suite.encryption)

	data := map[string]interface{}{
		"password":       "m987654321",
		"repeatPassword": "m987654321",
	}
	w := httptest.NewRecorder()

	req, _, _ := NewAuthenticatedRequest(suite.authService, suite.database, user, "POST", "/api/profile/devices", utils.MapToJsonBytesBuffer(data))

	router.ServeHTTP(w, req)
	a.Equal(http.StatusOK, w.Code, "Status code problem")
	suite.database.DB.Find(&user)

	response := ExtractResponseAsMap(w)
	devices, _ := response["data"].([]interface{})
	a.Equal("alaki", devices[0].(map[string]interface{})["city"].(string))
}
