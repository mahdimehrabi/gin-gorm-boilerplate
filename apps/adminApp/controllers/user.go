package controllers

import (
	authServices "boilerplate/apps/authApp/services"
	"boilerplate/apps/userApp/models"
	"boilerplate/apps/userApp/repositories"
	"boilerplate/apps/userApp/services"
	"boilerplate/core/infrastructure"
	"boilerplate/core/responses"
	_ "boilerplate/core/swagger"
	"boilerplate/utils"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserController struct {
	logger         infrastructure.Logger
	env            infrastructure.Env
	encryption     infrastructure.Encryption
	authService    authServices.AuthService
	userService    services.UserService
	userRepository repositories.UserRepository
	db             infrastructure.Database
}

func NewUserController(logger infrastructure.Logger,
	env infrastructure.Env,
	encryption infrastructure.Encryption,
	authService authServices.AuthService,
	userService services.UserService,
	userRepository repositories.UserRepository,
	db infrastructure.Database,
) UserController {
	return UserController{
		logger:         logger,
		env:            env,
		encryption:     encryption,
		authService:    authService,
		userService:    userService,
		userRepository: userRepository,
		db:             db,
	}
}

// @Summary get users list
// @Schemes
// @Description list of paginated response , authentication required
// @Tags admin
// @Accept json
// @Produce json
// @Success 200 {object} swagger.UsersListResponse
// @failure 401 {object} swagger.UnauthenticatedResponse
// @failure 403 {object} swagger.AccessForbiddenResponse
// @Router /admin/users [get]
func (uc UserController) ListUser(c *gin.Context) {
	uc.usersListResponse(c, "")
}
func (uc UserController) usersListResponse(c *gin.Context, message string) {
	pagination := utils.BuildPagination(c)
	users, count, err := uc.userRepository.GetAllUsers(pagination)
	if err != nil {
		uc.logger.Zap.Error("Failed to get users", err.Error())
		responses.ErrorJSON(c, http.StatusInternalServerError, gin.H{}, "Sorry an error occoured 😢")
		return
	}

	responses.JSON(c, http.StatusOK, gin.H{
		"count": count,
		"list":  users,
	}, message)
}

// @Summary create users
// @Schemes
// @Description create user and admin , admin only
// @Tags admin
// @Accept json
// @Produce json
// @Param email formData string true "unique email"
// @Param password formData string true "password that have at least 8 length and contain an alphabet and number "
// @Param repeatPassword formData string true "repeatPassword that have at least 8 length and contain an alphabet and number "
// @Param firstName formData string true "firstName"
// @Param lastName formData string true "lastName"
// @Param isAdmin formData bool true "isAdmin"
// @Success 200 {object} swagger.UsersListResponse
// @failure 401 {object} swagger.UnauthenticatedResponse
// @failure 403 {object} swagger.AccessForbiddenResponse
// @Router /admin/users [post]
func (uc UserController) CreateUser(c *gin.Context) {
	var userData models.CreateUserRequestAdmin
	if err := c.ShouldBindJSON(&userData); err != nil {
		fieldErrors := make(map[string]string, 0)
		if !utils.IsGoodPassword(userData.Password) {
			fieldErrors["password"] = "Password must contain at least one alphabet and one number and its length must be 8 characters or more"

		}
		responses.ValidationErrorsJSON(c, err, "", fieldErrors)
		return
	}
	if !utils.IsGoodPassword(userData.Password) {
		fieldErrors := map[string]string{
			"password": "Password must contain at least one alphabet and one number and its length must be 8 characters or more",
		}
		responses.ManualValidationErrorsJSON(c, fieldErrors, "")
		return
	}
	encryptedPassword := uc.encryption.SaltAndSha256Encrypt(userData.Password)
	user := models.User{
		Password:      encryptedPassword,
		FirstName:     userData.FirstName,
		LastName:      userData.LastName,
		Email:         userData.Email,
		IsAdmin:       userData.IsAdmin,
		VerifiedEmail: true,
	}
	err := uc.userService.CreateUser(&user)
	if err != nil {
		uc.logger.Zap.Error("Failed to create user ", err.Error())
		responses.ErrorJSON(c, http.StatusInternalServerError, gin.H{}, "Sorry an error occoured!")
		return
	}
	uc.usersListResponse(c, "User created succesfuly!")
}

// @Summary update user
// @Schemes
// @Description update user and admin , admin only
// @Tags admin
// @Accept json
// @Produce json
// @Param id path int true "user id"
// @Param email formData string true "unique email"
// @Param firstName formData string true "firstName"
// @Param lastName formData string true "lastName"
// @Param isAdmin formData bool true "isAdmin"
// @Success 200 {object} swagger.SingleUserResponse
// @failure 401 {object} swagger.UnauthenticatedResponse
// @failure 404 {object} swagger.NotFoundResponse
// @failure 403 {object} swagger.AccessForbiddenResponse
// @Router /admin/users/{id} [put]
func (uc UserController) UpdateUser(c *gin.Context) {
	var userData models.UpdateUserRequestAdmin
	id, _ := strconv.Atoi(c.Param("id"))
	userData.ID = uint(id)
	if err := c.ShouldBindJSON(&userData); err != nil {
		responses.ValidationErrorsJSON(c, err, "", map[string]string{})
		return
	}
	user, err := uc.userRepository.FindByField("id", c.Param("id"))
	if errors.Is(err, gorm.ErrRecordNotFound) {
		responses.ErrorJSON(c, http.StatusNotFound, gin.H{}, "No user found")
		return
	}
	if err != nil {
		uc.logger.Zap.Error("Error to find user", err.Error())
		responses.ErrorJSON(c, http.StatusInternalServerError, gin.H{}, "An error occured")
		return
	}
	updateUser := models.User{
		FirstName: userData.FirstName,
		LastName:  userData.LastName,
		Email:     userData.Email,
		IsAdmin:   userData.IsAdmin,
	}
	err = uc.db.DB.Model(&user).UpdateColumns(updateUser).Error
	if err != nil {
		uc.logger.Zap.Error("Error to update user", err.Error())
		responses.ErrorJSON(c, http.StatusInternalServerError, gin.H{}, "An error occured")
		return
	}

	responses.JSON(c, http.StatusOK, user.ToResponse(), "User updated successfuly")
}

// @Summary delete user
// @Schemes
// @Description delete user or admin , admin only
// @Tags admin
// @Accept json
// @Produce json
// @Param id path int true "user id"
// @Success 200 {object} swagger.UsersListResponse
// @failure 401 {object} swagger.UnauthenticatedResponse
// @failure 404 {object} swagger.NotFoundResponse
// @failure 403 {object} swagger.AccessForbiddenResponse
// @Router /admin/users/{id} [delete]
func (uc UserController) DeleteUser(c *gin.Context) {
	intId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		responses.ErrorJSON(c, http.StatusNotFound, gin.H{}, "No user found")
		return
	}
	_, err = uc.userRepository.FindByField("id", c.Param("id"))
	if errors.Is(err, gorm.ErrRecordNotFound) {
		responses.ErrorJSON(c, http.StatusNotFound, gin.H{}, "No user found")
		return
	}
	if err != nil {
		uc.logger.Zap.Error("Error to delete user", err.Error())
		responses.ErrorJSON(c, http.StatusInternalServerError, gin.H{}, "An error occured")
		return
	}

	var ID uint = uint(intId)
	err = uc.userRepository.DeleteByID(ID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		responses.ErrorJSON(c, http.StatusNotFound, gin.H{}, "No user found")
		return
	}
	if err != nil {
		uc.logger.Zap.Error("Error to delete user", err.Error())
		responses.ErrorJSON(c, http.StatusInternalServerError, gin.H{}, "An error occured")
		return
	}

	uc.usersListResponse(c, "User deleted succesfuly !")
}
