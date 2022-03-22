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

type UnauthenticatedResponse struct {
	FailedResponse
	Msg string `json:"msg" example:"You must login first!"`
}
