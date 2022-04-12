package models

import (
	"boilerplate/apps/genericApp/models"
	"time"

	"gorm.io/datatypes"
)

type User struct {
	models.Base
	IsAdmin   bool   `json:"-"`
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
		"Email":     m.Email,
		"FirstName": m.FirstName,
		"LastName":  m.LastName,
		"IsAdmin":   m.IsAdmin,
	}
}

func (u User) ToResponse() UserResponse {
	return UserResponse{
		BaseResponse: models.BaseResponse{
			CreatedAt: u.CreatedAt.Unix(),
			UpdatedAt: u.UpdatedAt.Unix(),
			ID:        u.ID,
		},
		IsAdmin:             u.IsAdmin,
		Email:               u.Email,
		FirstName:           u.FirstName,
		LastName:            u.LastName,
		Password:            u.Password,
		Devices:             u.Devices,
		VerifiedEmail:       u.VerifiedEmail,
		LastVerifyEmailDate: u.LastVerifyEmailDate,
		VerifyEmailToken:    u.VerifyEmailToken,
		ForgotPasswordToken: u.ForgotPasswordToken,
		LastForgotEmailDate: u.LastForgotEmailDate,
		Picture:             u.Picture,
	}
}

func UsersToResponses(users []User) []UserResponse {
	userResponses := make([]UserResponse, len(users))
	for i, v := range users {
		userResponses[i] = v.ToResponse()
	}
	return userResponses
}

type UserResponse struct {
	models.BaseResponse
	IsAdmin             bool           `json:"isAdmin"`
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

type CreateUserRequestAdmin struct {
	IsAdmin        bool   `json:"isAdmin"`
	Email          string `json:"email" binding:"required,uniqueDB=users&email,email"`
	FirstName      string `json:"firstName" binding:"required"`
	LastName       string `json:"lastName" binding:"required"`
	Password       string `json:"password" binding:"required"`
	RepeatPassword string `json:"repeatPassword" binding:"required,eqfield=Password"`
}

type UpdateUserRequestAdmin struct {
	ID        uint   `json:"-"` //this is required for unique validation in update
	IsAdmin   bool   `json:"isAdmin"`
	Email     string `json:"email" binding:"required,uniqueDB=users&email,email"`
	FirstName string `json:"firstName" binding:"required"`
	LastName  string `json:"lastName" binding:"required"`
}

type PictureRequest struct {
	Picture string `json:"picture" binding:"required"`
}
