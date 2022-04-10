package models

import (
	"time"

	"gorm.io/datatypes"
)

type User struct {
	Base
	Email     string `json:"email" binding:"required" gorm:"unique"`
	FirstName string `json:"firstName" binding:"required"`
	LastName  string `json:"lastName" binding:"required"`
	Password  string `json:"password" binding:"required"`

	//if you want to make user logout of specefic devices delete only specefic Device index
	// if you want to make user logout of all devices make this column empty
	Devices datatypes.JSON `json:"-"`

	VerifiedEmail bool `json:"verifiedEmail"`
	//last time send verification email date(use this field for implement limit for user resending verify email )
	LastVerifyEmailDate time.Time `json:"-"`
	VerifyEmailToken    string    `json:"-"`
	ForgotPasswordToken string    `json:"-"`
	LastForgotEmailDate time.Time `json:"-"`

	//resized profile picture store in this as json
	Picture string `json:"picture"`
}

// TableName gives table name of model
func (m User) TableName() string {
	return "users"
}

// ToMap convert User to map
func (m User) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"email":     m.Email,
		"firstName": m.FirstName,
		"LastName":  m.LastName,
	}
}

type UserResponse struct {
	Base
	Email               string         `json:"email"`
	FirstName           string         `json:"firstName" binding:"required"`
	LastName            string         `json:"lastName" binding:"required"`
	Password            string         `json:"-"`
	Devices             datatypes.JSON `json:"-"`
	VerifiedEmail       bool           `json:"verifiedEmail"`
	LastVerifyEmailDate time.Time      `json:"-"`
	VerifyEmailToken    string         `json:"-"`
	ForgotPasswordToken string         `json:"-"`
	LastForgotEmailDate time.Time      `json:"-"`
	Picture             string         `json:"picture"`
}
