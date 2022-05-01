package swagger

import "boilerplate/core/models"

type PaginateUsersData struct {
	Count int                   `json:"count" example:"10"`
	List  []models.UserResponse `json:"list"`
}

type UsersListResponse struct {
	SuccessResponse
	Data PaginateUsersData
}

type SingleUserResponse struct {
	SuccessResponse
	Data models.UserResponse
}
