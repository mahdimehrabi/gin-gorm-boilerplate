package tests

import (
	"fmt"
	"net/http"
	"net/http/httptest"
)

func (suite TestSuiteEnv) TestLoginCorrectCredentials() {
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
	response := ExtractResponseAsMap(w)
	str := w.Body.String()
	fmt.Println(response["data"]["accessToken"])
	fmt.Println(response)
}
