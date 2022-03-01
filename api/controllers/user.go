package controllers

import (
	"boilerplate/api/responses"
	"boilerplate/api/services"
	"boilerplate/constants"
	"boilerplate/infrastructure"
	"boilerplate/models"
	"boilerplate/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// UserController -> struct
type UserController struct {
	logger      infrastructure.Logger
	userService services.UserService
	env         infrastructure.Env
}

// NewUserController -> constructor
func NewUserController(
	logger infrastructure.Logger,
	userService services.UserService,
	env infrastructure.Env,
) UserController {
	return UserController{
		logger:      logger,
		userService: userService,
		env:         env,
	}
}

// CreateUser -> Create User
func (cc UserController) CreateUser(c *gin.Context) {
	user := models.User{}
	trx := c.MustGet(constants.DBTransaction).(*gorm.DB)

	if err := c.ShouldBindJSON(&user); err != nil {
		responses.ValidationErrorsJSON(c, err, "")
		return
	}

	if err := cc.userService.WithTrx(trx).CreateUser(user); err != nil {
		cc.logger.Zap.Error("Error [CreateUser] [db CreateUser]: ", err.Error())
		responses.ErrorJSON(c, http.StatusInternalServerError, gin.H{}, "Failed To Create User")
		return
	}

	responses.JSON(c, http.StatusOK, nil, "User Created Sucessfully")
}

// GetAllUser -> Get All User
func (cc UserController) GetAllUsers(c *gin.Context) {
	pagination := utils.BuildPagination(c)
	users, count, err := cc.userService.GetAllUsers(pagination)

	if err != nil {
		cc.logger.Zap.Error("Error finding user records", err.Error())
		responses.ErrorJSON(c, http.StatusInternalServerError, gin.H{}, "Failed To Find User")
		return
	}

	responses.JSONCount(c, http.StatusOK, users, "", count)
}
