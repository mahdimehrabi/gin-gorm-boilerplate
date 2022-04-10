package utils

import (
	"image/jpeg"
	"os"

	"github.com/nfnt/resize"
)

func ResizeImage(srcPath string, dstPath string, width uint, height uint) error {
	file, err := os.Open(srcPath)
	if err != nil {
		return err
	}

	// decode jpeg into image.Image
	img, err := jpeg.Decode(file)
	if err != nil {
		return err
	}
	file.Close()

	m := resize.Resize(width, height, img, resize.Lanczos3)

	out, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer out.Close()

	// write new image to file
	jpeg.Encode(out, m, nil)
	return nil
}
