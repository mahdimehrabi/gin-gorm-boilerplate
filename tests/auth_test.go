package tests

import (
	"net/http"
	"net/http/httptest"
)

func (suite TestSuiteEnv) TestLogin() {
	router := suite.router.Gin
	db := suite.database.DB
	a := suite.Assert()
	CreateUser(db)
	data := map[string]interface{}{
		"email":    "mahdi@gmail.com",
		"password": "m12345678",
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/auth/login", MapToJsonBytesBuffer(data))
	router.ServeHTTP(w, req)
	a.Equal(http.StatusOK, w.Code, "Status code problem")
}
