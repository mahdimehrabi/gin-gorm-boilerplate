package models

type PictureRequest struct {
	Picture string `json:"picture" binding:"required"`
}
