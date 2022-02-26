package utils

import (
	"image"
	"image/jpeg"
	"image/png"
	"mime/multipart"

	"github.com/nfnt/resize"
)

// CreateThumbnail image
func CreateThumbnail(file multipart.File, fileType string, width, height uint) (image.Image, error) {
	var img image.Image
	var err error

	// Seek back to beginning of file for CreateThumbnail
	if _, err := file.Seek(0, 0); err != nil {
		return nil, err
	}
	if fileType == "image/jpeg" || fileType == "image/jpg" {
		img, err = jpeg.Decode(file)
	} else if fileType == "image/png" {
		img, err = png.Decode(file)
	}
	if err != nil {
		return nil, err
	}
	thumbnail := resize.Resize(width, height, img, resize.Lanczos3)
	return thumbnail, nil
}
