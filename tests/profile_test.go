package tests

import (
	"boilerplate/utils"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"time"

	"gorm.io/datatypes"
)

func (suite TestSuiteEnv) TestChangePassword() {
	router := suite.router.Gin
	db := suite.database.DB
	a := suite.Assert()
	user := CreateUser("m12345678", db, suite.encryption)

	data := map[string]interface{}{
		"currentPassword": "m12345678",
		"password":        "m987654321",
		"repeatPassword":  "m987654321",
	}
	w := httptest.NewRecorder()

	req, _, _ := NewAuthenticatedRequest(suite.authService, suite.database, user, "POST", "/api/profile/change-password", utils.MapToJsonBytesBuffer(data))

	router.ServeHTTP(w, req)
	a.Equal(http.StatusOK, w.Code, "Status code problem")
	suite.database.DB.Find(&user)
	encryptedPassword := suite.encryption.SaltAndSha256Encrypt("m987654321")
	a.Equal(encryptedPassword, user.Password, "encrypt password problem")
	devicesBytes := []byte(user.Devices.String())
	devices, _ := utils.BytesJsonToMap(devicesBytes)
	a.Equal(len(devices), 0, "devices is not empty")

	//test with weak password
	user = CreateUser("m123456789", db, suite.encryption)
	data = map[string]interface{}{
		"currentPassword": "m12345678",
		"password":        "12345678",
		"repeatPassword":  "12345678",
	}
	w = httptest.NewRecorder()
	req, _, _ = NewAuthenticatedRequest(suite.authService, suite.database, user, "POST", "/api/profile/change-password", utils.MapToJsonBytesBuffer(data))
	router.ServeHTTP(w, req)
	a.Equal(http.StatusUnprocessableEntity, w.Code, "Status code problem")

	//test with wrong current password
	user = CreateUser("m123456789", db, suite.encryption)
	data = map[string]interface{}{
		"currentPassword": "521421251",
		"password":        "12345678",
		"repeatPassword":  "12345678",
	}
	w = httptest.NewRecorder()
	req, _, _ = NewAuthenticatedRequest(suite.authService, suite.database, user, "POST", "/api/profile/change-password", utils.MapToJsonBytesBuffer(data))
	router.ServeHTTP(w, req)
	a.Equal(http.StatusUnprocessableEntity, w.Code, "Status code problem")
}

func (suite TestSuiteEnv) TestGetLoggedInDevices() {
	router := suite.router.Gin
	db := suite.database.DB
	a := suite.Assert()
	user := CreateUser("m12345678", db, suite.encryption)

	data := map[string]interface{}{
		"password":       "m987654321",
		"repeatPassword": "m987654321",
	}
	w := httptest.NewRecorder()

	req, _, _ := NewAuthenticatedRequest(suite.authService, suite.database, user, "GET", "/api/profile/devices", utils.MapToJsonBytesBuffer(data))

	router.ServeHTTP(w, req)
	a.Equal(http.StatusOK, w.Code, "Status code problem")
	suite.database.DB.Find(&user)

	response := ExtractResponseAsMap(w)
	devices, _ := response["data"].(map[string]interface{})["devices"].([]interface{})
	a.Equal("alaki", devices[0].(map[string]interface{})["city"].(string))
}

func (suite TestSuiteEnv) TestTerminateDevice() {
	router := suite.router.Gin
	db := suite.database.DB
	a := suite.Assert()
	user := CreateUser("m12345678", db, suite.encryption)

	w := httptest.NewRecorder()
	deviceToken := utils.GenerateRandomCode(18)
	data := map[string]interface{}{
		"token": deviceToken,
	}
	req, _, _ := NewAuthenticatedRequestCustomDeviceToken(
		suite.authService,
		suite.database,
		user,
		"POST",
		"/api/profile/terminate-device",
		utils.MapToJsonBytesBuffer(data),
		deviceToken,
	)
	suite.database.DB.Find(&user)
	devicesBytes := []byte(user.Devices.String())
	devices, _ := utils.BytesJsonToMap(devicesBytes)
	a.NotNil(devices[deviceToken])

	router.ServeHTTP(w, req)
	a.Equal(http.StatusOK, w.Code, "Status code problem")
	suite.database.DB.Find(&user)

	devicesBytes = []byte(user.Devices.String())
	devices, _ = utils.BytesJsonToMap(devicesBytes)
	a.Nil(devices[deviceToken])

	//test wrong token
	w = httptest.NewRecorder()
	data = map[string]interface{}{
		"token": "wrong-device-token",
	}
	req, deviceToken, _ = NewAuthenticatedRequest(
		suite.authService,
		suite.database,
		user,
		"POST",
		"/api/profile/terminate-device",
		utils.MapToJsonBytesBuffer(data),
	)
	router.ServeHTTP(w, req)
	a.Equal(http.StatusNotFound, w.Code, "Status code problem")

}

func (suite TestSuiteEnv) TerminateDevicesExceptMe() {
	router := suite.router.Gin
	db := suite.database.DB
	a := suite.Assert()
	user := CreateUser("m12345678", db, suite.encryption)

	devices := make(map[string]interface{})
	device1Token := utils.GenerateRandomCode(17) + "1"
	device2Token := utils.GenerateRandomCode(17) + "2"
	devices[device1Token] = map[string]string{
		"ip":         "1.1.1.1",
		"city":       "alaki",
		"date":       strconv.Itoa(int(time.Now().Unix())),
		"deviceName": "windows10-chrome",
	}
	devices[device2Token] = map[string]string{
		"ip":         "1.1.1.1",
		"city":       "alaki",
		"date":       strconv.Itoa(int(time.Now().Unix())),
		"deviceName": "android-chrome",
	}
	user.Devices = datatypes.JSON(utils.MapToJsonBytesBuffer(devices).String())
	suite.database.DB.Save(&user)

	data := map[string]interface{}{}

	req, _, _ := NewAuthenticatedRequest(
		suite.authService,
		suite.database,
		user,
		"POST",
		"/api/profile/terminate-devices-except-me",
		utils.MapToJsonBytesBuffer(data),
	)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	a.Equal(http.StatusOK, w.Code)

	devicesBytes := []byte(user.Devices.String())
	devices, _ = utils.BytesJsonToMap(devicesBytes)
	a.NotNil(devices[device1Token])
	a.Nil(devices[device1Token])
}

func (suite TestSuiteEnv) TestProfilePicture() {
	router := suite.router.Gin
	db := suite.database.DB
	a := suite.Assert()
	user := CreateUser("m12345678", db, suite.encryption)

	f, _ := os.Open("./media/duck.png")
	defer f.Close()

	data, writer, _ := CreateFileRequestBody("picture", f)
	w := httptest.NewRecorder()

	req, _, _ := NewAuthenticatedRequest(
		suite.authService,
		suite.database,
		user,
		"POST",
		"/api/profile/upload-profile-picture",
		data)
	req.Header.Add("Content-Type", writer.FormDataContentType())
	router.ServeHTTP(w, req)
	a.Equal(http.StatusOK, w.Code, "Status code problem")
	suite.database.DB.Find(&user)
	a.Regexp("/media/users/"+strconv.Itoa(int(user.ID))+"*", user.Picture)
}
