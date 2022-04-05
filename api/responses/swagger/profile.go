package swagger

import "boilerplate/models"

type DevicesResponse struct {
	Msg  string                    `json:"msg" example:""`
	Ok   bool                      `json:"ok" example:"true"`
	Data models.DeviceListResponse `json:"data" `
}
