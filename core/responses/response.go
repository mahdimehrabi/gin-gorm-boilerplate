package responses

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// JSON : json response function
func JSON(c *gin.Context, statusCode int, data interface{}, message string) {
	c.JSON(statusCode, gin.H{"data": data, "msg": message, "ok": true})
}

// ErrorJSON : json error response function
func ErrorJSON(c *gin.Context, statusCode int, data interface{}, message string) {
	c.JSON(statusCode, gin.H{"data": data, "msg": message, "ok": false})
}

//automatic validation errors just send validator.ValidationErrors and this func automatic generate response
func ValidationErrorsJSON(c *gin.Context, err error, message string, extraFieldErrors map[string]string) {
	if message == "" {
		message = "Please review entered data"
	}
	defer func() {
		if recoveryMessage := recover(); recoveryMessage != nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"data": gin.H{}, "msg": message, "ok": false})
		}
	}()
	//auto generate errors from validation erros
	var ve validator.ValidationErrors
	errs := make(map[string]string)
	typeErr := new(json.UnmarshalTypeError)
	if errors.As(err, &ve) {
		for _, fe := range ve {
			field := fe.Field()
			lowerCaseField := strings.ToLower(field[0:1]) + field[1:]
			msg := MsgForTag(fe.Tag(), field, fe.Param())
			errs[lowerCaseField] = msg
		}
	} else if errors.As(err, &typeErr) {
		errs[typeErr.Field] = fmt.Sprintf("Incorrect type %s for field %s", typeErr.Value, typeErr.Field)
	}

	//merge errors with extraFieldErrors
	for k, v := range extraFieldErrors {
		errs[k] = v
	}
	c.JSON(http.StatusUnprocessableEntity, gin.H{"data": gin.H{"errors": errs}, "msg": message, "ok": false})
}

//manual validation errors, please send errors like map[fieldName]errorMsg in second parameter
func ManualValidationErrorsJSON(c *gin.Context, fieldErrors map[string]string, message string) {
	if message == "" {
		message = "Please review entered data"
	}
	c.JSON(http.StatusUnprocessableEntity, gin.H{"data": gin.H{"errors": fieldErrors}, "msg": message, "ok": false})
}

// JSONCount : json response function
func JSONCount(c *gin.Context, statusCode int, data interface{}, message string, count int64) {
	c.JSON(statusCode, gin.H{"data": data, "count": count, "msg": message, "ok": true})
}
