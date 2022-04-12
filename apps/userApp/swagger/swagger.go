package swagger

import "boilerplate/apps/auth"

type DevicesResponse struct {
	Msg  string                  `json:"msg" example:""`
	Ok   bool                    `json:"ok" example:"true"`
	Data auth.DeviceListResponse `json:"data" `
}
