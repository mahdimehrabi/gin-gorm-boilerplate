package tests

import (
	"net/http"
	"net/http/httptest"
)

func (suite *TestSuiteEnv) TestPing() {
	router := suite.router.Gin
	a := suite.Assert()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ping", nil)
	router.ServeHTTP(w, req)
	a.Equal(w.Code, http.StatusOK, "Status code problem")
	a.Contains(w.Body.String(), "pong", "Contain data problem")
}
