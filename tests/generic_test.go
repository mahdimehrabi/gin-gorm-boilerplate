package tests

import (
	"fmt"
	"net/http"
	"net/http/httptest"
)

func (suite *TestSuiteEnv) TestPing() {
	router := suite.router.Gin

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ping", nil)
	fmt.Println(router)
	router.ServeHTTP(w, req)
}
