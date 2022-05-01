package tests

import (
	"boilerplate/core/models"
	"boilerplate/utils"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
)

func (suite TestSuiteEnv) TestAdminCreateUser() {
	router := suite.router.Gin
	db := suite.database.DB
	a := suite.Assert()
	admin := CreateAdmin("m1234567", db, suite.encryption)
	var beforeUserCount int64
	db.Model(models.User{}).Count(&beforeUserCount)

	/*test with completed credentials*/
	data := map[string]interface{}{
		"email":          "mahdi@gmail.com",
		"password":       "m12345678",
		"repeatPassword": "m12345678",
		"firstName":      "mahdi",
		"lastName":       "mehrabi",
		"isAdmin":        false,
	}
	w := httptest.NewRecorder()
	req, _, _ := NewAuthenticatedRequest(
		suite.authService,
		suite.database,
		admin,
		"POST",
		"/api/admin/users",
		utils.MapToJsonBytesBuffer(data))

	router.ServeHTTP(w, req)
	a.Equal(http.StatusOK, w.Code, "Status code problem")
	res := w.Body.String()
	fmt.Println(res)
	var afterUserCount int64
	db.Model(models.User{}).Count(&afterUserCount)
	var user models.User
	db.Model(models.User{}).Last(&user)
	a.True(afterUserCount == beforeUserCount+1, "User count problem")

	/*test with duplicate email*/
	w = httptest.NewRecorder()
	req, _, _ = NewAuthenticatedRequest(
		suite.authService,
		suite.database,
		admin,
		"POST",
		"/api/admin/users",
		utils.MapToJsonBytesBuffer(data))
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
		"isAdmin":   false,
	}
	w = httptest.NewRecorder()
	req, _, _ = NewAuthenticatedRequest(
		suite.authService,
		suite.database,
		admin,
		"POST",
		"/api/admin/users",
		utils.MapToJsonBytesBuffer(data))
	router.ServeHTTP(w, req)
	a.Equal(http.StatusUnprocessableEntity, w.Code, "Status code problem")
	db.Model(models.User{}).Count(&afterUserCount)
	a.True(afterUserCount == beforeUserCount+1, "User count problem")
}

func (suite TestSuiteEnv) TestGetUsers() {
	router := suite.router.Gin
	db := suite.database.DB
	a := suite.Assert()
	admin := CreateAdmin("m12345678", db, suite.encryption)
	//test correct
	data := map[string]interface{}{}
	w := httptest.NewRecorder()

	req, _, _ := NewAuthenticatedRequest(
		suite.authService,
		suite.database,
		admin,
		"GET",
		"/api/admin/users",
		utils.MapToJsonBytesBuffer(data))

	router.ServeHTTP(w, req)
	a.Equal(http.StatusOK, w.Code, "Status code problem")
	a.Contains(w.Body.String(), admin.Email)

	//test with not admin user
	data = map[string]interface{}{}
	w = httptest.NewRecorder()
	user := CreateUser("m12345678", db, suite.encryption)
	req, _, _ = NewAuthenticatedRequest(
		suite.authService,
		suite.database,
		user,
		"GET",
		"/api/admin/users",
		utils.MapToJsonBytesBuffer(data))
	router.ServeHTTP(w, req)
	a.Equal(http.StatusForbidden, w.Code, "Status code problem")

}

func (suite TestSuiteEnv) TestAdminUpdateUser() {
	router := suite.router.Gin
	db := suite.database.DB
	a := suite.Assert()
	admin := CreateAdmin("m1234567", db, suite.encryption)
	user := CreateUser("m1234567", db, suite.encryption)

	/*test with completed credentials*/
	data := map[string]interface{}{
		"email":     "mahdi@gmail.com",
		"firstName": "mahdi",
		"lastName":  "mehrabi",
		"isAdmin":   false,
	}
	w := httptest.NewRecorder()
	req, _, _ := NewAuthenticatedRequest(
		suite.authService,
		suite.database,
		admin,
		"PUT",
		"/api/admin/users/"+strconv.Itoa(int(user.ID)),
		utils.MapToJsonBytesBuffer(data))

	router.ServeHTTP(w, req)
	a.Equal(http.StatusOK, w.Code, "Status code problem")
	db.Find(&user)
	a.Equal(user.Email, data["email"], "User email equality problem")

	//test update with his current
	w = httptest.NewRecorder()
	req, _, _ = NewAuthenticatedRequest(
		suite.authService,
		suite.database,
		admin,
		"PUT",
		"/api/admin/users/"+strconv.Itoa(int(user.ID)),
		utils.MapToJsonBytesBuffer(data))

	router.ServeHTTP(w, req)
	a.Equal(http.StatusOK, w.Code, "Status code problem")

	/*test with duplicate email*/
	user2 := CreateUser("m1234567", db, suite.encryption)

	w = httptest.NewRecorder()
	req, _, _ = NewAuthenticatedRequest(
		suite.authService,
		suite.database,
		admin,
		"PUT",
		"/api/admin/users/"+strconv.Itoa(int(user2.ID)),
		utils.MapToJsonBytesBuffer(data))

	router.ServeHTTP(w, req)
	a.Equal(http.StatusUnprocessableEntity, w.Code, "Status code problem")
	db.Find(&user2)
	a.NotEqual(user2.Email, data["email"], "User email equality problem")
}

func (suite TestSuiteEnv) TestAdminDeleteUser() {
	router := suite.router.Gin
	db := suite.database.DB
	a := suite.Assert()
	admin := CreateAdmin("m1234567", db, suite.encryption)
	user := CreateAdmin("m1234567", db, suite.encryption)
	user.Picture = "/media/picture.jpg"
	db.Save(&user)
	os.Create(suite.env.BasePath + user.Picture)
	CreateAdmin("m1234567", db, suite.encryption)
	CreateAdmin("m1234567", db, suite.encryption)
	var beforeUserCount int64
	db.Model(models.User{}).Count(&beforeUserCount)

	/*test with completed credentials*/
	data := map[string]interface{}{}
	w := httptest.NewRecorder()
	req, _, _ := NewAuthenticatedRequest(
		suite.authService,
		suite.database,
		admin,
		"DELETE",
		"/api/admin/users/"+strconv.Itoa(int(user.ID)),
		utils.MapToJsonBytesBuffer(data))

	router.ServeHTTP(w, req)
	a.Equal(http.StatusOK, w.Code, "Status code problem")
	var afterUserCount int64
	db.Model(models.User{}).Count(&afterUserCount)
	a.True(afterUserCount == beforeUserCount-1, "User count problem")
	//check profile picture deleted
	_, err := os.Stat(suite.env.BasePath + user.Picture)
	a.True(os.IsNotExist(err))

	/*test with wrong id */
	db.Model(models.User{}).Count(&beforeUserCount)
	w = httptest.NewRecorder()
	req, _, _ = NewAuthenticatedRequest(
		suite.authService,
		suite.database,
		admin,
		"DELETE",
		"/api/admin/users/"+"25115512",
		utils.MapToJsonBytesBuffer(data))
	router.ServeHTTP(w, req)
	a.Equal(http.StatusNotFound, w.Code, "Status code problem")
	db.Model(models.User{}).Count(&afterUserCount)
	a.True(afterUserCount == beforeUserCount, "User count problem")
}
