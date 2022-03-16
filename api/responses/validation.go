package responses

import "fmt"

func MsgForTag(tag string, fieldName string) string {
	switch tag {
	case "required":
		return "This field is required"
	case "email":
		return "Please enter a valid email like example@example.com"
	case "uniqueDB":
		return fmt.Sprintf("Entered %s is already exist", fieldName)
	case "numeric":
		return "You must use numeric value for this field"
	}
	return "You cannot use this data for this field"
}
