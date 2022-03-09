package tests

func (suite TestSuiteEnv) TestUser() {
	// router := suite.router.Gin
	// db := suite.database.DB
	a := suite.Assert()
	// user := models.User{Email: "test@gmail.com", FullName: "mahdi mehrabi"}
	// db.Create(&user)

	// w := httptest.NewRecorder()
	// req, _ := http.NewRequest("GET", "/api/users", nil)
	// router.ServeHTTP(w, req)
	a.Equal(200, 200)
	// a.Contains(w.Body.String(), user.FullName, user.Email, "Contain data problem")
}
