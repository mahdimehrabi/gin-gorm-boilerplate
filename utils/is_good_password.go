package utils

import "regexp"

func IsGoodPassword(password string) bool {
	containAlphabet, _ := regexp.MatchString("[a-zA-Z]", password)
	containNumber, _ := regexp.MatchString("\\d", password)
	return containAlphabet && containNumber && len(password) >= 8
}
