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
