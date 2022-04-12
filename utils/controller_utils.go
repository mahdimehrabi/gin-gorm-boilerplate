package utils

import (
	"boilerplate/core/responses"
	"strings"

	"github.com/gin-gonic/gin"
)

//validate and upload file
//uploadPath => path of file without extension like /media/mahdi/image
func UploadFile(uploadPath string, c *gin.Context, key string, types []string) (bool, string, error) {
	file, err := c.FormFile(key)
	if err != nil {
		fieldErrors := make(map[string]string, 0)
		fieldErrors[key] = "You must upload a image with type of jpeg or png"
		responses.ValidationErrorsJSON(c, err, "", fieldErrors)
		return false, "", nil
	}
	fileSlice := strings.Split(file.Filename, ".")
	extension := fileSlice[len(fileSlice)-1]
	if !StringInSlice(extension, types) {
		fieldErrors := make(map[string]string, 0)
		fieldErrors[key] = "You must upload a image with type of jpeg or png"
		responses.ValidationErrorsJSON(c, err, "", fieldErrors)
		return false, "", nil
	}
	uploadPath += "." + extension
	err = c.SaveUploadedFile(file, uploadPath)
	if err != nil {
		return false, "", err
	}
	return true, uploadPath, nil
}
