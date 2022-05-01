package tests

import (
	"boilerplate/apps/authApp/services"
	"boilerplate/core/models"
	"boilerplate/utils"
	"bytes"
	"net/http"
	"net/http/httptest"
	"time"

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

	//test not verified email account
	suite.database.DB.Model(&user).Update("verified_email", false)
	data = map[string]interface{}{
		"email":      user.Email,
		"password":   "m12345678",
		"deviceName": "windows10-chrome",
	}
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/auth/login", utils.MapToJsonBytesBuffer(data))
	router.ServeHTTP(w, req)
	a.Equal(http.StatusBadRequest, w.Code, "Status code problem")

}

func (suite TestSuiteEnv) TestRegister() {
	router := suite.router.Gin
	db := suite.database.DB
	a := suite.Assert()
	var beforeUserCount int64
	db.Model(models.User{}).Count(&beforeUserCount)

	/*test with completed credentials*/
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
	a.True(len(user.VerifyEmailToken) > 35)

	/*test with duplicate email*/
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/auth/register", utils.MapToJsonBytesBuffer(data))
	router.ServeHTTP(w, req)
	a.Equal(http.StatusUnprocessableEntity, w.Code, "Status code problem")
	db.Model(models.User{}).Count(&afterUserCount)
	a.True(afterUserCount == beforeUserCount+1, "User count problem")

	/*test with weak password*/
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
	a.True(1 == 1)
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

func (suite TestSuiteEnv) TestVerifyEmail() {
	router := suite.router.Gin
	db := suite.database.DB
	a := suite.Assert()
	user := CreateUser("m12345678", db, suite.encryption)
	suite.database.DB.Model(&user).Update("verified_email", false)

	//test with wrong token
	w := httptest.NewRecorder()
	data := map[string]interface{}{
		"token": utils.GenerateRandomCode(40),
	}
	req, _ := http.NewRequest("POST", "/api/auth/verify-email", utils.MapToJsonBytesBuffer(data))
	router.ServeHTTP(w, req)
	a.Equal(http.StatusBadRequest, w.Code, "status code problem")
	suite.database.DB.Find(&user)
	a.False(user.VerifiedEmail)

	//test without token
	w = httptest.NewRecorder()
	data = map[string]interface{}{}
	req, _ = http.NewRequest("POST", "/api/auth/verify-email", utils.MapToJsonBytesBuffer(data))
	router.ServeHTTP(w, req)
	a.Equal(http.StatusUnprocessableEntity, w.Code, "status code problem")
	suite.database.DB.Find(&user)
	a.False(user.VerifiedEmail)

	//test with right token
	w = httptest.NewRecorder()
	data = map[string]interface{}{
		"token": user.VerifyEmailToken,
	}
	req, _ = http.NewRequest("POST", "/api/auth/verify-email", utils.MapToJsonBytesBuffer(data))
	router.ServeHTTP(w, req)
	a.Equal(http.StatusOK, w.Code, "status code problem")
	suite.database.DB.Find(&user)
	a.True(user.VerifiedEmail)
}

func (suite TestSuiteEnv) TestForgotEmail() {
	router := suite.router.Gin
	db := suite.database.DB
	a := suite.Assert()
	user := CreateUser("m12345678", db, suite.encryption)

	//test with right email
	w := httptest.NewRecorder()
	data := map[string]interface{}{
		"email": user.Email,
	}
	req, _ := http.NewRequest("POST", "/api/auth/forgot-password", utils.MapToJsonBytesBuffer(data))
	router.ServeHTTP(w, req)
	a.Equal(http.StatusOK, w.Code, "status code problem")
	suite.database.DB.Find(&user)
	oldToken := user.ForgotPasswordToken
	a.True(len(oldToken) > 20)

	//test with wrong email
	w = httptest.NewRecorder()
	data = map[string]interface{}{
		"email": "masfasf@gmail.com",
	}
	req, _ = http.NewRequest("POST", "/api/auth/forgot-password", utils.MapToJsonBytesBuffer(data))
	router.ServeHTTP(w, req)
	a.Equal(http.StatusOK, w.Code, "status code problem")
	suite.database.DB.Find(&user)
	a.Equal(oldToken, user.ForgotPasswordToken)

	//test with right email but twice in row
	w = httptest.NewRecorder()
	data = map[string]interface{}{
		"email": user.Email,
	}
	suite.database.DB.Model(&user).UpdateColumn("last_forgot_email_date", time.Now())
	req, _ = http.NewRequest("POST", "/api/auth/forgot-password", utils.MapToJsonBytesBuffer(data))
	router.ServeHTTP(w, req)
	a.Equal(http.StatusBadRequest, w.Code, "status code problem")
	suite.database.DB.Find(&user)
	a.Equal(oldToken, user.ForgotPasswordToken)
}

func (suite TestSuiteEnv) TestRecoveryPassword() {
	router := suite.router.Gin
	db := suite.database.DB
	a := suite.Assert()
	user := CreateUser("m12345678", db, suite.encryption)
	suite.database.DB.Model(&user).UpdateColumn("forgot_password_token", utils.GenerateRandomCode(40))

	//test with right info
	w := httptest.NewRecorder()
	data := map[string]interface{}{
		"token":          user.ForgotPasswordToken,
		"password":       "m987654321",
		"repeatPassword": "m987654321",
	}

	req, _, _ := NewAuthenticatedRequest(suite.authService, suite.database, user, "POST", "/api/auth/recover-password", utils.MapToJsonBytesBuffer(data))

	router.ServeHTTP(w, req)
	a.Equal(http.StatusOK, w.Code, "Status code problem")
	suite.database.DB.Find(&user)
	encryptedPassword := suite.encryption.SaltAndSha256Encrypt("m987654321")
	a.Equal(encryptedPassword, user.Password, "encrypt password problem")
	a.Empty(user.ForgotPasswordToken)

	devicesBytes := []byte(user.Devices.String())
	devices, _ := utils.BytesJsonToMap(devicesBytes)
	a.Equal(len(devices), 0, "devices is not empty")

	//test with weak password
	user = CreateUser("m123456777", db, suite.encryption)
	data = map[string]interface{}{
		"token":          user.ForgotPasswordToken,
		"password":       "12345678",
		"repeatPassword": "12345678",
	}
	w = httptest.NewRecorder()
	req, _, _ = NewAuthenticatedRequest(suite.authService, suite.database, user, "POST", "/api/auth/recover-password", utils.MapToJsonBytesBuffer(data))
	router.ServeHTTP(w, req)
	a.Equal(http.StatusUnprocessableEntity, w.Code, "Status code problem")

	//test with wrong token
	user = CreateUser("m123456765", db, suite.encryption)
	data = map[string]interface{}{
		"token":          utils.GenerateRandomCode(40),
		"password":       "m987654321",
		"repeatPassword": "m987654321",
	}
	w = httptest.NewRecorder()
	req, _, _ = NewAuthenticatedRequest(suite.authService, suite.database, user, "POST", "/api/auth/recover-password", utils.MapToJsonBytesBuffer(data))
	router.ServeHTTP(w, req)
	a.Equal(http.StatusNotFound, w.Code, "Status code problem")

}

func (suite TestSuiteEnv) TestResendEmail() {
	router := suite.router.Gin
	db := suite.database.DB
	a := suite.Assert()
	user := CreateUser("m12345678", db, suite.encryption)
	suite.database.DB.Model(&user).Update("verified_email", false)

	//test with right email
	w := httptest.NewRecorder()
	data := map[string]interface{}{
		"email": user.Email,
	}
	req, _ := http.NewRequest("POST", "/api/auth/resend-verify-email", utils.MapToJsonBytesBuffer(data))
	router.ServeHTTP(w, req)
	a.Equal(http.StatusOK, w.Code, "status code problem")
	suite.database.DB.Find(&user)
	token := user.VerifyEmailToken
	a.True(len(token) > 20)

	//test with wrong email
	w = httptest.NewRecorder()
	data = map[string]interface{}{
		"email": "masfasf@gmail.com",
	}
	req, _ = http.NewRequest("POST", "/api/auth/resend-verify-email", utils.MapToJsonBytesBuffer(data))
	router.ServeHTTP(w, req)
	a.Equal(http.StatusOK, w.Code, "status code problem")

	//test with right email but twice in row
	w = httptest.NewRecorder()
	data = map[string]interface{}{
		"email": user.Email,
	}
	suite.database.DB.Model(&user).UpdateColumn("last_verify_email_date", time.Now())
	req, _ = http.NewRequest("POST", "/api/auth/resend-verify-email", utils.MapToJsonBytesBuffer(data))
	router.ServeHTTP(w, req)
	a.Equal(http.StatusBadRequest, w.Code, "status code problem")
	suite.database.DB.Find(&user)
}
