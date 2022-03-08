package controllers

import (
	"boilerplate/api/repositories"
	"boilerplate/api/responses"
	"boilerplate/api/services"
	"boilerplate/constants"
	"boilerplate/infrastructure"
	"boilerplate/models"
	"boilerplate/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// UserController -> struct
type UserController struct {
	logger         infrastructure.Logger
	userService    services.UserService
	userRepository repositories.UserRepository
	env            infrastructure.Env
}

// NewUserController -> constructor
func NewUserController(
	logger infrastructure.Logger,
	userService services.UserService,
	env infrastructure.Env,
	userRepository repositories.UserRepository,
) UserController {
	return UserController{
		logger:         logger,
		userService:    userService,
		env:            env,
		userRepository: userRepository,
	}
}

// CreateUser -> Create User
func (uc UserController) CreateUser(c *gin.Context) {
	user := models.User{}
	trx := c.MustGet(constants.DBTransaction).(*gorm.DB)

	if err := c.ShouldBindJSON(&user); err != nil {
		responses.ValidationErrorsJSON(c, err, "")
		return
	}

	if err := uc.userService.WithTrx(trx).CreateUser(&user); err != nil {
		uc.logger.Zap.Error("Error [CreateUser] [db CreateUser]: ", err.Error())
		responses.ErrorJSON(c, http.StatusInternalServerError, gin.H{}, "Failed To Create User")
		return
	}

	responses.JSON(c, http.StatusOK, user, "User created sucessfully")
}

// GetAllUser -> Get All User
func (uc UserController) GetAllUsers(c *gin.Context) {
	pagination := utils.BuildPagination(c)
	users, count, err := uc.userService.GetAllUsers(pagination)

	if err != nil {
		uc.logger.Zap.Error("Error finding user records", err.Error())
		responses.ErrorJSON(c, http.StatusInternalServerError, gin.H{}, "Failed To Find User")
		return
	}

	responses.JSONCount(c, http.StatusOK, users, "", count)
}

func (uc UserController) DeleteUser(c *gin.Context) {
	id := c.Param("id")
	if exist, _ := uc.userRepository.IsExist("id", id); exist == false {
		responses.ErrorJSON(c, http.StatusNotFound, gin.H{}, "User not found !")
		return
	}
	intId, _ := strconv.Atoi(id)
	if err := uc.userService.DeleteUserByID(uint(intId)); err != nil {
		uc.logger.Zap.Error("Error in DeleteUserByID : ", err.Error())
		responses.ErrorJSON(c, http.StatusInternalServerError, gin.H{}, "Failed To Create User")
		return
	}

	responses.JSON(c, http.StatusOK, gin.H{}, "User deleted sucessfully")
}
