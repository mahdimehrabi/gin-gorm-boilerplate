package models

type User struct {
	Base
	Email     string `json:"email" binding:"required" gorm:"unique"`
	FirstName string `json:"firstName" binding:"required"`
	LastName  string `json:"lastName" binding:"required"`
	Password  string `json:"password" binding:"required"`
	//change this field to a random string with length less then 50
	// if you want to make user logout of all devices
	RefreshTokenSecret string `json:"-"` //change this
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
	Email              string `json:"email"`
	FirstName          string `json:"firstName" binding:"required"`
	LastName           string `json:"lastName" binding:"required"`
	Password           string `json:"-"`
	RefreshTokenSecret string `json:"-"`
}
