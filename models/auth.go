package models

type Register struct {
	Password       string `json:"password" binding:"required"`
	RepeatPassword string `json:"repeatPassword" binding:"required,eqfield=Password"`
	Email          string `json:"email" binding:"required,uniqueDB=users&email,email"`
	FirstName      string `json:"firstName" binding:"required"`
	LastName       string `json:"LastName" binding:"required"`
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

type RefreshTokenRequest struct {
	RefreshToken string `json:"refreshToken" binding:"required"`
}
