package swagger

type EmptyData struct {
}

//for override
type SuccessResponse struct {
	Msg  string    `json:"msg" example:"Successful message"`
	Ok   bool      `json:"ok" example:"true"`
	Data EmptyData `json:"data"`
}

//for override
type FailedResponse struct {
	Msg  string    `json:"msg" example:"Error or warnnig message"`
	Ok   bool      `json:"ok" example:"false"`
	Data EmptyData `json:"data"`
}

//for override
type NotFoundResponse struct {
	Msg  string    `json:"msg" example:"404 not found!"`
	Ok   bool      `json:"ok" example:"false"`
	Data EmptyData `json:"data"`
}

type AccessForbiddenResponse struct {
	Msg  string    `json:"msg" example:"Sorry you don't have access to visit this page!"`
	Ok   bool      `json:"ok" example:"false"`
	Data EmptyData `json:"data"`
}

type validationErrors struct {
	Errors map[string]string `json:"errors" example:"field1:This field is required,field2:This field must be numeric"`
}

//for override
type FailedValidationResponse struct {
	Msg  string           `json:"msg" example:"Please review your entered data"`
	Ok   bool             `json:"ok" example:"false"`
	Data validationErrors `json:"data"`
}

type PingResponse struct {
	Msg  string            `json:"msg" example:"pong"`
	Ok   bool              `json:"ok" example:"true"`
	Data map[string]string `json:"data" example:"pingpong:🏓🏓🏓🏓🏓🏓"`
}
