package swagger

type PingResponse struct {
	Msg  string            `json:"msg" example:"pong"`
	Ok   bool              `json:"ok" example:"true"`
	Data map[string]string `json:"data" example:"pingpong:🏓🏓🏓🏓🏓🏓"`
}
