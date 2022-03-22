package tests

import (
	"boilerplate/models"
	"net/http"
	"net/http/httptest"
)

func (suite TestSuiteEnv) TestLogin() {
	router := suite.router.Gin
	db := suite.database.DB
	a := suite.Assert()
	CreateUser("m12345678", db, suite.encryption)

	//test correct credentials
	data := map[string]interface{}{
		"email":    "mahdi@gmail.com",
		"password": "m12345678",
	}
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/auth/login", MapToJsonBytesBuffer(data))
	router.ServeHTTP(w, req)
	a.Equal(http.StatusOK, w.Code, "Status code problem")
	response := ExtractResponseAsMap(w)
	dt, _ := response["data"].(map[string]interface{})
	accessToken, _ := dt["accessToken"].(string)
	refreshToken, _ := dt["refreshToken"].(string)
	a.True(len(accessToken) > 7, "Access token invalid")
	a.True(len(refreshToken) > 7, "Refresh token invalid")

	//test access token
	data = map[string]interface{}{
		"accessToken": accessToken,
	}
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/auth/access-token-verify", MapToJsonBytesBuffer(data))
	router.ServeHTTP(w, req)
	a.Equal(http.StatusOK, w.Code, "Status code problem")

	//test refresh token
	data = map[string]interface{}{
		"refreshToken": refreshToken,
	}
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/auth/renew-access-token", MapToJsonBytesBuffer(data))
	router.ServeHTTP(w, req)
	a.Equal(http.StatusOK, w.Code, "Status code problem")

	//test wrong email
	data = map[string]interface{}{
		"email":    "mahdi1@gmail.com",
		"password": "m12345678",
	}
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/auth/login", MapToJsonBytesBuffer(data))
	router.ServeHTTP(w, req)
	a.Equal(http.StatusUnauthorized, w.Code, "Status code problem")

	//test wrong password
	data = map[string]interface{}{
		"email":    "mahdi1@gmail.com",
		"password": "m123456781",
	}
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/auth/login", MapToJsonBytesBuffer(data))
	router.ServeHTTP(w, req)
	a.Equal(http.StatusUnauthorized, w.Code, "Status code problem")

	//test without email
	data = map[string]interface{}{
		"password": "m123456781",
	}
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/auth/login", MapToJsonBytesBuffer(data))
	router.ServeHTTP(w, req)
	a.Equal(http.StatusUnprocessableEntity, w.Code, "Status code problem")

	//test without password
	data = map[string]interface{}{
		"email": "mahdi1@gmail.com",
	}
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/auth/login", MapToJsonBytesBuffer(data))
	router.ServeHTTP(w, req)
	a.Equal(http.StatusUnprocessableEntity, w.Code, "Status code problem")
}

func (suite TestSuiteEnv) TestRegister() {
	router := suite.router.Gin
	db := suite.database.DB
	a := suite.Assert()
	var beforeUserCount int64
	db.Model(models.User{}).Count(&beforeUserCount)

	//test with completed credentials
	data := map[string]interface{}{
		"email":          "mahdi@gmail.com",
		"password":       "m12345678",
		"repeatPassword": "m12345678",
		"firstName":      "mahdi",
		"lastName":       "mehrabi",
	}
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/auth/register", MapToJsonBytesBuffer(data))
	router.ServeHTTP(w, req)
	a.Equal(http.StatusOK, w.Code, "Status code problem")
	response := ExtractResponseAsMap(w)
	dt, _ := response["data"].(map[string]interface{})
	accessToken, _ := dt["accessToken"].(string)
	refreshToken, _ := dt["refreshToken"].(string)
	a.True(len(accessToken) > 7, "Access token invalid")
	a.True(len(refreshToken) > 7, "Refresh token invalid")

	//test access token
	data = map[string]interface{}{
		"accessToken": accessToken,
	}
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/auth/access-token-verify", MapToJsonBytesBuffer(data))
	router.ServeHTTP(w, req)
	a.Equal(http.StatusOK, w.Code, "Status code problem")

	//test refresh token
	data = map[string]interface{}{
		"refreshToken": refreshToken,
	}
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/auth/renew-access-token", MapToJsonBytesBuffer(data))
	router.ServeHTTP(w, req)
	a.Equal(http.StatusOK, w.Code, "Status code problem")

	var afterUserCount int64
	db.Model(models.User{}).Count(&afterUserCount)
	a.True(afterUserCount == beforeUserCount+1, "User count problem")

	//test with duplicate email
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/auth/register", MapToJsonBytesBuffer(data))
	router.ServeHTTP(w, req)
	a.Equal(http.StatusUnprocessableEntity, w.Code, "Status code problem")
	db.Model(models.User{}).Count(&afterUserCount)
	a.True(afterUserCount == beforeUserCount+1, "User count problem")

	//test with weak password
	data = map[string]interface{}{
		"email":     "mahdi@gmail.com",
		"password":  "12345678",
		"firstName": "mahdi",
		"lastName":  "mehrabi",
	}
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/auth/register", MapToJsonBytesBuffer(data))
	router.ServeHTTP(w, req)
	a.Equal(http.StatusUnprocessableEntity, w.Code, "Status code problem")
	db.Model(models.User{}).Count(&afterUserCount)
	a.True(afterUserCount == beforeUserCount+1, "User count problem")
}
