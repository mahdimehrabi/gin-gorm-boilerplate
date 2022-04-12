package swagger

import (
	"boilerplate/apps/auth"
	"boilerplate/apps/generic"
)

type SuccessLoginResponse struct {
	generic.SuccessResponse
	Data auth.LoginResponse `json:"data"`
}

type FailedLoginResponse struct {
	generic.FailedResponse
	Msg string `json:"msg" example:"No user found with entered credentials"`
}

type SuccessVerifyAccessTokenResponse struct {
	generic.SuccessResponse
	Data auth.AccessTokenRes `json:"data"`
}

type UnauthenticatedResponse struct {
	generic.FailedResponse
	Msg string `json:"msg" example:"You must login first!"`
}
