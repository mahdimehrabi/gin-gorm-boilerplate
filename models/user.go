package models

type User struct {
	Base
	Email    string `json:"email"`
	FullName string `json:"full_name"`
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
