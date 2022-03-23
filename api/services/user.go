package services

import (
	"boilerplate/api/repositories"
	"boilerplate/infrastructure"
	"boilerplate/models"
	"boilerplate/utils"

	"gorm.io/gorm"
)

// UserService -> struct
type UserService struct {
	repositories repositories.UserRepository
	db           infrastructure.Database
}

// NewUserService -> creates a new Userservice
func NewUserService(repositories repositories.UserRepository, db infrastructure.Database) UserService {
	return UserService{
		repositories: repositories,
		db:           db,
	}
}

// WithTrx -> enables repositories with transaction
func (us UserService) WithTrx(trxHandle *gorm.DB) UserService {
	us.repositories = us.repositories.WithTrx(trxHandle)
	return us
}

// CreateUser -> call to create the User
func (us UserService) CreateUser(user *models.User) error {
	err := us.repositories.Create(user)
	return err
}

// CreateUser -> call to create the User
func (us UserService) UpdatePassword(user *models.User, password string) error {
	return us.db.DB.Model(user).UpdateColumns(models.User{Password: password}).Error
}

// GetAllUser -> call to get all the User
func (us UserService) GetAllUsers(pagination utils.Pagination) ([]models.User, int64, error) {
	return us.repositories.GetAllUsers(pagination)
}

func (us UserService) DeleteUserByID(id uint) error {
	return us.repositories.DeleteByID(id)
}
