package swagger

import "boilerplate/core/models"

type LoginResponse struct {
	SuccessResponse
	Data models.LoginResponse `json:"data"`
}

type FailedLoginResponse struct {
	FailedResponse
	Msg string `json:"msg" example:"No user found with entered credentials"`
}

type SuccessVerifyAccessTokenResponse struct {
	SuccessResponse
	Data models.AccessTokenRes `json:"data"`
}

type UnauthenticatedResponse struct {
	FailedResponse
	Msg string `json:"msg" example:"You must login first!"`
}
