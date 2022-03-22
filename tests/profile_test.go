package tests

import (
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

	req, _ := NewAuthenticatedRequest(suite.authService, user, "POST", "/api/profile/change-password", MapToJsonBytesBuffer(data))

	router.ServeHTTP(w, req)
	a.Equal(http.StatusOK, w.Code, "Status code problem")
	suite.database.DB.Find(&user)
	encryptedPassword := suite.encryption.SaltAndSha256Encrypt("m987654321")
	a.Equal(encryptedPassword, user.Password, "encrypt password problem")
	a.True(user.MustLogin, "must logout problem")

	//test with weak password
	user = CreateUser("m12345678", db, suite.encryption)
	data = map[string]interface{}{
		"password":       "12345678",
		"repeatPassword": "12345678",
	}
	w = httptest.NewRecorder()
	req, _ = NewAuthenticatedRequest(suite.authService, user, "POST", "/api/profile/change-password", MapToJsonBytesBuffer(data))
	router.ServeHTTP(w, req)
	a.Equal(http.StatusUnprocessableEntity, w.Code, "Status code problem")
}
