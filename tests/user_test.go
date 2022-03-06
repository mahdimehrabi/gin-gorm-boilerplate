package tests

import (
	"boilerplate/models"
	"fmt"
	"net/http"
	"net/http/httptest"
)

func (suite TestSuiteEnv) TestUser() {
	fmt.Println("---------------afsfsafas------------------")
	router := suite.router.Gin
	db := suite.database.DB
	a := suite.Assert()
	user := models.User{Email: "test@gmail.com", FullName: "mahdi mehrabi"}
	db.Create(&user)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/users", nil)
	router.ServeHTTP(w, req)
	a.Equal(w.Code, http.StatusOK, "Status code problem")
	a.Contains(w.Body.String(), user.FullName, user.Email, "Contain data problem")
}
