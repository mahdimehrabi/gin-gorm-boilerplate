package models

type User struct {
	Base
	Email    string `json:"email" binding:"required" gorm:"unique"`
	FullName string `json:"fullName" binding:"required"`
	Password string `json:"password" binding:"required"`
	// Make this field true when user change password or request for logout in other devices or ...
	// Make sure you make this field to false on login
	// Make sure user cannot renew access token if this field is true
	MustLogout bool `json:"-"`
}

// TableName gives table name of model
func (m User) TableName() string {
	return "users"
}

// ToMap convert User to map
func (m User) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"email":     m.Email,
		"full_name": m.FullName,
	}
}

type UserResponse struct {
	Base
	Email      string `json:"email"`
	FullName   string `json:"fullName"`
	Password   string `json:"-"`
	MustLogout bool   `json:"-"`
}

type CreateUser struct {
	Password string `json:"password" binding:"required"`
	Email    string `json:"email" binding:"required,uniqueDB"`
	FullName string `json:"fullName" binding:"required"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	AccessToken  string       `json:"accessToken"`
	RefreshToken string       `json:"refreshToken"`
	User         UserResponse `json:"user"`
}
