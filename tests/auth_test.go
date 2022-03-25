package tests

import (
	"boilerplate/api/services"
	"boilerplate/utils"
	"bytes"
	"net/http"
	"net/http/httptest"

	"github.com/golang-jwt/jwt"
)

func (suite TestSuiteEnv) TestLogin() {
	router := suite.router.Gin
	db := suite.database.DB
	a := suite.Assert()
	user := CreateUser("m12345678", db, suite.encryption)

	//test correct credentials
	data := map[string]interface{}{
		"email":    user.Email,
		"password": "m12345678",
		"device":   "windows10-chrome",
	}
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/auth/login", utils.MapToJsonBytesBuffer(data))
	router.ServeHTTP(w, req)
	a.Equal(http.StatusOK, w.Code, "Status code problem")
	response := ExtractResponseAsMap(w)
	dt, _ := response["data"].(map[string]interface{})
	accessToken, _ := dt["accessToken"].(string)
	refreshToken, _ := dt["refreshToken"].(string)
	a.True(len(accessToken) > 7, "Access token invalid")
	a.True(len(refreshToken) > 7, "Refresh token invalid")
	db.Find(&user)
	a.True(len(user.Devices.String()) > 7, "Devices not set")

	//check device token
	var atClaims jwt.MapClaims
	refreshSecret := "refresh" + suite.env.Secret
	_, atClaims, _ = services.DecodeToken(refreshToken, refreshSecret)
	deviceToken := atClaims["deviceToken"].(string)

	devicesBytes := []byte(user.Devices.String())
	devices, _ := utils.BytesJsonToMap(devicesBytes)

	a.NotNil(devices[deviceToken], "devices not set")
	a.Equal(devices[deviceToken].(map[string]interface{})["city"], "alaki")

	//test access token
	data = map[string]interface{}{
		"accessToken": accessToken,
	}
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/auth/access-token-verify", utils.MapToJsonBytesBuffer(data))
	router.ServeHTTP(w, req)
	a.Equal(http.StatusOK, w.Code, "Status code problem")

	//test refresh token
	data = map[string]interface{}{
		"refreshToken": refreshToken,
	}
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/auth/renew-access-token", utils.MapToJsonBytesBuffer(data))
	router.ServeHTTP(w, req)
	a.Equal(http.StatusOK, w.Code, "Status code problem")

	//test wrong email
	data = map[string]interface{}{
		"email":    "mahdi1@gmail.com",
		"password": "m12345678",
		"device":   "windows10-chrome",
	}
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/auth/login", utils.MapToJsonBytesBuffer(data))
	router.ServeHTTP(w, req)
	a.Equal(http.StatusUnauthorized, w.Code, "Status code problem")

	//test wrong password
	data = map[string]interface{}{
		"email":    "mahdi1@gmail.com",
		"password": "m123456781",
		"device":   "windows10-chrome",
	}
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/auth/login", utils.MapToJsonBytesBuffer(data))
	router.ServeHTTP(w, req)
	a.Equal(http.StatusUnauthorized, w.Code, "Status code problem")

	//test without email
	data = map[string]interface{}{
		"password": "m123456781",
		"device":   "windows10-chrome",
	}
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/auth/login", utils.MapToJsonBytesBuffer(data))
	router.ServeHTTP(w, req)
	a.Equal(http.StatusUnprocessableEntity, w.Code, "Status code problem")

	//test without password
	data = map[string]interface{}{
		"email":  "mahdi1@gmail.com",
		"device": "windows10-chrome",
	}
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/auth/login", utils.MapToJsonBytesBuffer(data))
	router.ServeHTTP(w, req)
	a.Equal(http.StatusUnprocessableEntity, w.Code, "Status code problem")
}

// func (suite TestSuiteEnv) TestRegister() {
// 	router := suite.router.Gin
// 	db := suite.database.DB
// 	a := suite.Assert()
// 	var beforeUserCount int64
// 	db.Model(models.User{}).Count(&beforeUserCount)

// 	//test with completed credentials
// 	data := map[string]interface{}{
// 		"email":          "mahdi@gmail.com",
// 		"password":       "m12345678",
// 		"repeatPassword": "m12345678",
// 		"firstName":      "mahdi",
// 		"lastName":       "mehrabi",
// 	}
// 	w := httptest.NewRecorder()
// 	req, _ := http.NewRequest("POST", "/api/auth/register", utils.MapToJsonBytesBuffer(data))
// 	router.ServeHTTP(w, req)
// 	a.Equal(http.StatusOK, w.Code, "Status code problem")
// 	response := ExtractResponseAsMap(w)
// 	dt, _ := response["data"].(map[string]interface{})
// 	accessToken, _ := dt["accessToken"].(string)
// 	refreshToken, _ := dt["refreshToken"].(string)
// 	a.True(len(accessToken) > 7, "Access token invalid")
// 	a.True(len(refreshToken) > 7, "Refresh token invalid")

// 	//test access token
// 	data = map[string]interface{}{
// 		"accessToken": accessToken,
// 	}
// 	w = httptest.NewRecorder()
// 	req, _ = http.NewRequest("POST", "/api/auth/access-token-verify", utils.MapToJsonBytesBuffer(data))
// 	router.ServeHTTP(w, req)
// 	a.Equal(http.StatusOK, w.Code, "Status code problem")

// 	//test refresh token
// 	data = map[string]interface{}{
// 		"refreshToken": refreshToken,
// 	}
// 	w = httptest.NewRecorder()
// 	req, _ = http.NewRequest("POST", "/api/auth/renew-access-token", utils.MapToJsonBytesBuffer(data))
// 	router.ServeHTTP(w, req)
// 	a.Equal(http.StatusOK, w.Code, "Status code problem")

// 	var afterUserCount int64
// 	db.Model(models.User{}).Count(&afterUserCount)
// 	a.True(afterUserCount == beforeUserCount+1, "User count problem")

// 	//test with duplicate email
// 	w = httptest.NewRecorder()
// 	req, _ = http.NewRequest("POST", "/api/auth/register", utils.MapToJsonBytesBuffer(data))
// 	router.ServeHTTP(w, req)
// 	a.Equal(http.StatusUnprocessableEntity, w.Code, "Status code problem")
// 	db.Model(models.User{}).Count(&afterUserCount)
// 	a.True(afterUserCount == beforeUserCount+1, "User count problem")

// 	//test with weak password
// 	data = map[string]interface{}{
// 		"email":     "mahdi1@gmail.com",
// 		"password":  "12345678",
// 		"firstName": "mahdi",
// 		"lastName":  "mehrabi",
// 	}
// 	w = httptest.NewRecorder()
// 	req, _ = http.NewRequest("POST", "/api/auth/register", utils.MapToJsonBytesBuffer(data))
// 	router.ServeHTTP(w, req)
// 	a.Equal(http.StatusUnprocessableEntity, w.Code, "Status code problem")
// 	db.Model(models.User{}).Count(&afterUserCount)
// 	a.True(afterUserCount == beforeUserCount+1, "User count problem")
// }

func (suite TestSuiteEnv) TestLogout() {
	router := suite.router.Gin
	db := suite.database.DB
	a := suite.Assert()
	user := CreateUser("m12345678", db, suite.encryption)

	w := httptest.NewRecorder()
	data := new(bytes.Buffer)
	req, _, _ := NewAuthenticatedRequest(suite.authService, suite.database, user, "POST", "/api/auth/logout", data)

	router.ServeHTTP(w, req)
	a.Equal(http.StatusOK, w.Code, "Status code problem")
	suite.database.DB.Find(&user)
}
