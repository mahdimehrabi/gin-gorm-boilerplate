package models

type Register struct {
	Password       string `json:"password" binding:"required"`
	RepeatPassword string `json:"repeatPassword" binding:"required,eqfield=Password"`
	Email          string `json:"email" binding:"required,uniqueDB=users&email,email"`
	FirstName      string `json:"firstName" binding:"required"`
	LastName       string `json:"LastName" binding:"required"`
}

type ChangePassword struct {
	Password       string `json:"password" binding:"required"`
	RepeatPassword string `json:"repeatPassword" binding:"required,eqfield=Password"`
}

type LoginRequest struct {
	Email      string `json:"email" binding:"required"`
	Password   string `json:"password" binding:"required"`
	DeviceName string `json:"deviceName" binding:"required"`
}

type LoginResponse struct {
	AccessToken     string       `json:"accessToken"`
	RefreshToken    string       `json:"refreshToken"`
	ExpRefreshToken string       `json:"expRefreshToken"`
	ExpAccessToken  string       `json:"expAccessToken"`
	User            UserResponse `json:"user"`
}

type AccessTokenReq struct {
	AccessToken string `json:"accessToken" binding:"required"`
}

type AccessTokenRes struct {
	AccessToken    string `json:"accessToken" binding:"required"`
	ExpAccessToken string `json:"expAccessToken" binding:"required"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refreshToken" binding:"required"`
}

type TokenRequest struct {
	Token string `json:"token" binding:"required"`
}

type EmailRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type Device struct {
	Ip         string `json:"ip"`
	City       string `json:"city"`
	Date       string `json:"date"`
	DeviceName string `json:"deviceName"`
}
