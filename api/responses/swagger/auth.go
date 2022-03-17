package swagger

import (
	"boilerplate/models"
)

type RegisterLoginResponse struct {
	SuccessResonse
	Data models.LoginResponse `json:"data"`
}

type FailedLoginResponse struct {
	FailedResponse
	Msg string `json:"msg" example:"No user found with entered credentials"`
}

type SuccessVerifyAccessTokenResponse struct {
	SuccessResonse
	Data models.AccessTokenReqRes `json:"data"`
}
