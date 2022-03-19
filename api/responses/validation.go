package responses

import (
	"fmt"
	"strings"
)

func MsgForTag(tag string, fieldName string, param string) string {
	switch tag {
	case "required":
		return "This field is required"
	case "email":
		return "Please enter a valid email like example@example.com"
	case "uniqueDB":
		return fmt.Sprintf("Entered %s is already exist", fieldName)
	case "numeric":
		return "You must use numeric value for this field"
	case "eqfield":
		param = strings.ToLower(param)
		return fmt.Sprintf("Field %s must be equal to field %s", fieldName, param)
	}
	return "You cannot use this data for this field"
}
