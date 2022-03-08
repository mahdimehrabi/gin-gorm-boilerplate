package repositories

import (
	"boilerplate/infrastructure"
	"boilerplate/models"
	"boilerplate/utils"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

// UserRepository -> database structure
type UserRepository struct {
	db     infrastructure.Database
	logger infrastructure.Logger
}

// NewUserRepository -> creates a new User repository
func NewUserRepository(db infrastructure.Database, logger infrastructure.Logger) UserRepository {
	return UserRepository{
		db:     db,
		logger: logger,
	}
}

// WithTrx enables repository with transaction
func (c UserRepository) WithTrx(trxHandle *gorm.DB) UserRepository {
	if trxHandle == nil {
		c.logger.Zap.Error("Transaction Database not found in gin context. ")
		return c
	}
	c.db.DB = trxHandle
	return c
}

// Save -> User
func (c UserRepository) Create(User *models.User) error {
	return c.db.DB.Create(User).Error
}

func (c UserRepository) FindByField(field string, value string) (user models.User, err error) {
	err = c.db.DB.Where(fmt.Sprintf("%s= ?", field), value).First(&user).Error
	return
}

func (c UserRepository) DeleteByID(id uint) error {
	return c.db.DB.Where("id=?", id).Delete(&models.User{}).Error
}

func (c UserRepository) IsExist(field string, value string) (bool, error) {
	_, err := c.FindByField(field, value)
	if err == nil {
		return true, nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}
	return false, err
}

// GetAllUser -> Get All users
func (c UserRepository) GetAllUsers(pagination utils.Pagination) ([]models.User, int64, error) {
	var users []models.User
	var totalRows int64 = 0
	queryBuilder := c.db.DB.Limit(pagination.PageSize).Offset(pagination.Offset).Order("created_at desc")
	queryBuilder = queryBuilder.Model(&models.User{})

	if pagination.Keyword != "" {
		searchQuery := "%" + pagination.Keyword + "%"
		queryBuilder.Where(c.db.DB.Where("`users`.`name` LIKE ?", searchQuery))
	}

	err := queryBuilder.
		Find(&users).
		Offset(-1).
		Limit(-1).
		Count(&totalRows).Error
	return users, totalRows, err
}
