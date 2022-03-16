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

type PingResponse struct {
	Msg  string            `json:"msg" example:"pong"`
	Ok   bool              `json:"ok" example:"true"`
	Data map[string]string `json:"data" example:"pingpong:🏓🏓🏓🏓🏓🏓"`
}
