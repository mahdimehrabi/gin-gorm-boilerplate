package utils

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
)

// JsonStringify convert map to string without error handling
func JsonStringify(input gin.H) string {
	retStr, _ := json.Marshal(input)
	return string(retStr)
}
