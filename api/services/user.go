package services

import (
	"boilerplate/api/repositories"
	"boilerplate/models"
	"boilerplate/utils"

	"gorm.io/gorm"
)

// UserService -> struct
type UserService struct {
	repositories repositories.UserRepository
}

// NewUserService -> creates a new Userservice
func NewUserService(repositories repositories.UserRepository) UserService {
	return UserService{
		repositories: repositories,
	}
}

// WithTrx -> enables repositories with transaction
func (c UserService) WithTrx(trxHandle *gorm.DB) UserService {
	c.repositories = c.repositories.WithTrx(trxHandle)
	return c
}

// CreateUser -> call to create the User
func (c UserService) CreateUser(user *models.User) error {
	err := c.repositories.Create(user)
	return err
}

// GetAllUser -> call to get all the User
func (c UserService) GetAllUsers(pagination utils.Pagination) ([]models.User, int64, error) {
	return c.repositories.GetAllUsers(pagination)
}
