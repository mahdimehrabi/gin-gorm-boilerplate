package swagger

import "boilerplate/apps/authApp/models"

type DevicesResponse struct {
	Msg  string                    `json:"msg" example:""`
	Ok   bool                      `json:"ok" example:"true"`
	Data models.DeviceListResponse `json:"data" `
}
