package swagger

//for override
type SuccessResonse struct {
	Msg  string            `json:"msg" example:"Operation was successful"`
	Ok   bool              `json:"ok" example:"true"`
	Data map[string]string //please rewrite this
}

//for override
type FailedResonse struct {
	Msg  string            `json:"msg" example:"Operation was not successful"`
	Ok   bool              `json:"ok" example:"false"`
	Data map[string]string //please rewrite this
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
