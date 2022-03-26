package tests

import (
	"boilerplate/api/services"
	"boilerplate/models"
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

	//test correct credentials with first device
	data := map[string]interface{}{
		"email":      user.Email,
		"password":   "m12345678",
		"deviceName": "windows10-chrome",
	}
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/auth/login", utils.MapToJsonBytesBuffer(data))
	router.ServeHTTP(w, req)
	a.Equal(http.StatusOK, w.Code, "Status code problem")
	response := ExtractResponseAsMap(w)
	dt, _ := response["data"].(map[string]interface{})
	accessTokenDevice1, _ := dt["accessToken"].(string)
	refreshTokenDevice1, _ := dt["refreshToken"].(string)
	a.True(len(accessTokenDevice1) > 7, "Access token invalid")
	a.True(len(refreshTokenDevice1) > 7, "Refresh token invalid")
	db.Find(&user)
	a.True(len(user.Devices.String()) > 7, "Devices not set")
	var atClaims jwt.MapClaims

	//login with second device
	data = map[string]interface{}{
		"email":      user.Email,
		"password":   "m12345678",
		"deviceName": "android10-chrome",
	}
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/auth/login", utils.MapToJsonBytesBuffer(data))
	router.ServeHTTP(w, req)
	a.Equal(http.StatusOK, w.Code, "Status code problem")
	response = ExtractResponseAsMap(w)
	dt, _ = response["data"].(map[string]interface{})
	accessTokenDevice2, _ := dt["accessToken"].(string)
	refreshTokenDevice2, _ := dt["refreshToken"].(string)
	a.True(len(accessTokenDevice2) > 7, "Access token invalid")
	a.True(len(refreshTokenDevice2) > 7, "Refresh token invalid")
	db.Find(&user)
	a.True(len(user.Devices.String()) > 7, "Devices not set")

	//check device1 login data
	refreshSecret := "refresh" + suite.env.Secret
	_, atClaims, _ = services.DecodeToken(refreshTokenDevice1, refreshSecret)
	device1Token := atClaims["deviceToken"].(string)

	devicesBytes := []byte(user.Devices.String())
	devices, _ := utils.BytesJsonToMap(devicesBytes)

	a.NotNil(devices[device1Token], "devices not set")
	a.Equal(devices[device1Token].(map[string]interface{})["city"], "unknown")
	a.Equal("windows10-chrome", devices[device1Token].(map[string]interface{})["deviceName"])

	//check device2 login data
	_, atClaims, _ = services.DecodeToken(refreshTokenDevice2, refreshSecret)
	device2Token := atClaims["deviceToken"].(string)

	devicesBytes = []byte(user.Devices.String())
	devices, _ = utils.BytesJsonToMap(devicesBytes)

	a.NotNil(devices[device2Token], "devices not set")
	a.Equal(devices[device2Token].(map[string]interface{})["city"], "unknown")
	a.Equal("android10-chrome", devices[device2Token].(map[string]interface{})["deviceName"])

	//test access token
	data = map[string]interface{}{
		"accessToken": accessTokenDevice1,
	}
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/auth/access-token-verify", utils.MapToJsonBytesBuffer(data))
	router.ServeHTTP(w, req)
	a.Equal(http.StatusOK, w.Code, "Status code problem")

	//test refresh token
	data = map[string]interface{}{
		"refreshToken": refreshTokenDevice1,
	}
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/auth/renew-access-token", utils.MapToJsonBytesBuffer(data))
	router.ServeHTTP(w, req)
	a.Equal(http.StatusOK, w.Code, "Status code problem")

	//test wrong email
	data = map[string]interface{}{
		"email":      "mahdi1@gmail.com",
		"password":   "m12345678",
		"deviceName": "windows10-chrome",
	}
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/auth/login", utils.MapToJsonBytesBuffer(data))
	router.ServeHTTP(w, req)
	a.Equal(http.StatusUnauthorized, w.Code, "Status code problem")

	//test wrong password
	data = map[string]interface{}{
		"email":      "mahdi1@gmail.com",
		"password":   "m123456781",
		"deviceName": "windows10-chrome",
	}
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/auth/login", utils.MapToJsonBytesBuffer(data))
	router.ServeHTTP(w, req)
	a.Equal(http.StatusUnauthorized, w.Code, "Status code problem")

	//test without email
	data = map[string]interface{}{
		"password":   "m123456781",
		"deviceName": "windows10-chrome",
	}
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/auth/login", utils.MapToJsonBytesBuffer(data))
	router.ServeHTTP(w, req)
	a.Equal(http.StatusUnprocessableEntity, w.Code, "Status code problem")

	//test without password
	data = map[string]interface{}{
		"email":      "mahdi1@gmail.com",
		"deviceName": "windows10-chrome",
	}
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/auth/login", utils.MapToJsonBytesBuffer(data))
	router.ServeHTTP(w, req)
	a.Equal(http.StatusUnprocessableEntity, w.Code, "Status code problem")
}

func (suite TestSuiteEnv) TestRegister() {
	router := suite.router.Gin
	db := suite.database.DB
	a := suite.Assert()
	var beforeUserCount int64
	db.Model(models.User{}).Count(&beforeUserCount)
	emailMock := new(SendEmailMock)

	//test with completed credentials
	data := map[string]interface{}{
		"email":          "mahdi@gmail.com",
		"password":       "m12345678",
		"repeatPassword": "m12345678",
		"firstName":      "mahdi",
		"lastName":       "mehrabi",
	}
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/auth/register", utils.MapToJsonBytesBuffer(data))
	router.ServeHTTP(w, req)
	a.Equal(http.StatusOK, w.Code, "Status code problem")

	var afterUserCount int64
	db.Model(models.User{}).Count(&afterUserCount)
	var user models.User
	db.Model(models.User{}).Last(&user)
	a.True(afterUserCount == beforeUserCount+1, "User count problem")
	a.False(user.VerifiedEmail, "Email verified must be false")
	a.False(user.LastVerifyEmailDate.IsZero(), "last verify email date must set after")
	emailMock.AssertNumberOfCalls(suite.T(), "SendEmail", 1)

	//test with duplicate email
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/auth/register", utils.MapToJsonBytesBuffer(data))
	router.ServeHTTP(w, req)
	a.Equal(http.StatusUnprocessableEntity, w.Code, "Status code problem")
	db.Model(models.User{}).Count(&afterUserCount)
	a.True(afterUserCount == beforeUserCount+1, "User count problem")

	//test with weak password
	data = map[string]interface{}{
		"email":     "mahdi1@gmail.com",
		"password":  "12345678",
		"firstName": "mahdi",
		"lastName":  "mehrabi",
	}
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/auth/register", utils.MapToJsonBytesBuffer(data))
	router.ServeHTTP(w, req)
	a.Equal(http.StatusUnprocessableEntity, w.Code, "Status code problem")
	db.Model(models.User{}).Count(&afterUserCount)
	a.True(afterUserCount == beforeUserCount+1, "User count problem")
}

func (suite TestSuiteEnv) TestLogout() {
	router := suite.router.Gin
	db := suite.database.DB
	a := suite.Assert()
	user := CreateUser("m12345678", db, suite.encryption)

	w := httptest.NewRecorder()
	data := new(bytes.Buffer)
	req, deviceToken, _ := NewAuthenticatedRequest(suite.authService, suite.database, user, "POST", "/api/auth/logout", data)
	suite.database.DB.Find(&user)
	devicesBytes := []byte(user.Devices.String())
	devices, _ := utils.BytesJsonToMap(devicesBytes)
	a.NotNil(devices[deviceToken])

	router.ServeHTTP(w, req)
	a.Equal(http.StatusOK, w.Code, "Status code problem")
	suite.database.DB.Find(&user)

	devicesBytes = []byte(user.Devices.String())
	devices, _ = utils.BytesJsonToMap(devicesBytes)
	a.Nil(devices[deviceToken])
}

func (suite TestSuiteEnv) TestRenewAccessTokenWithUnexistDeviceToken() {
	router := suite.router.Gin
	db := suite.database.DB
	a := suite.Assert()
	user := CreateUser("m12345678", db, suite.encryption)
	tokensData, _ := suite.authService.CreateTokens(user, "fake-device-token")

	data := map[string]interface{}{
		"refreshToken": tokensData["refreshToken"],
	}
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/auth/renew-access-token", utils.MapToJsonBytesBuffer(data))

	router.ServeHTTP(w, req)
	a.Equal(http.StatusBadRequest, w.Code, "Status code problem")
}
