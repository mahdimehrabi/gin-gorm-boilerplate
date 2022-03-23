package models

type User struct {
	Base
	Email     string `json:"email" binding:"required" gorm:"unique"`
	FirstName string `json:"firstName" binding:"required"`
	LastName  string `json:"lastName" binding:"required"`
	Password  string `json:"password" binding:"required"`
	// Make this field true when user change password or request for logout in other devices or ...
	// Make sure you make this field to false on login
	// Make sure user cannot renew access token if this field is true
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
	Email     string `json:"email"`
	FirstName string `json:"firstName" binding:"required"`
	LastName  string `json:"lastName" binding:"required"`
	Password  string `json:"-"`
}
