package models

type Register struct {
	Password string `json:"password" binding:"required"`
	Email    string `json:"email" binding:"required,uniqueDB=users&email,email"`
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

type AccessTokenReqRes struct {
	AccessToken string `json:"accessToken" binding:"required"`
}
