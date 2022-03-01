package responses

type ValidationError struct {
	Field string `json:"field"`
	Msg   string `json:"msg"`
}

func MsgForTag(tag string) string {
	switch tag {
	case "required":
		return "This field is required"
	case "email":
		return "Please email a valid email like example@example.com"
	}
	return "You cannot use this data for this field"
}
