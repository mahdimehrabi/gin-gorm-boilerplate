package responses

import (
	"errors"
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

func ValidationErrorsJSON(c *gin.Context, err error, message string) {
	if message == "" {
		message = "Please review entered data"
	}
	defer func() {
		if recoveryMessage := recover(); recoveryMessage != nil {
			c.JSON(http.StatusBadRequest, gin.H{"data": gin.H{}, "msg": message, "ok": false})
		}
	}()
	var ve validator.ValidationErrors
	errs := make([]ValidationError, 0)
	if errors.As(err, &ve) {
		for _, fe := range ve {
			field := fe.Field()
			field = strings.ToLower(field[0:1]) + field[1:]
			msg := MsgForTag(fe.Tag())
			ve := ValidationError{field, msg}
			errs = append(errs, ve)
		}
	}
	c.JSON(http.StatusBadRequest, gin.H{"data": gin.H{"errors": errs}, "msg": message, "ok": false})
}

// JSONCount : json response function
func JSONCount(c *gin.Context, statusCode int, data interface{}, message string, count int64) {
	c.JSON(statusCode, gin.H{"data": data, "count": count, "msg": message, "ok": true})
}
