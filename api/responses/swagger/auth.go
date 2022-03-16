package swagger

import (
	"boilerplate/models"
)

type RegisterLoginResponse struct {
	SuccessResonse
	Data models.LoginResponse `json:"data"`
}
