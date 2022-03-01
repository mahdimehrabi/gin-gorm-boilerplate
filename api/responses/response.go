package responses

import (
	"github.com/gin-gonic/gin"
)

// JSON : json response function
func JSON(c *gin.Context, statusCode int, data interface{}, message string) {
	c.JSON(statusCode, gin.H{"data": data, "message": message, "ok": true})
}

// ErrorJSON : json error response function
func ErrorJSON(c *gin.Context, statusCode int, data interface{}, message string) {
	c.JSON(statusCode, gin.H{"data": data, "message": message, "ok": false})
}

func ValidationErrorsJSON(c *gin.Context, statusCode int, data interface{}, message string) {
	if message == "" {
		message = "Please fill all required field"
	}
	c.JSON(statusCode, gin.H{"data": data, "message": message, "ok": false})
}

// JSONCount : json response function
func JSONCount(c *gin.Context, statusCode int, data interface{}, message string, count int64) {
	c.JSON(statusCode, gin.H{"data": data, "count": count, "message": message, "ok": true})
}
