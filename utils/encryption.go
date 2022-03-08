package utils

import (
	"crypto/sha256"
	"encoding/hex"
)

func Sha256Encrypt(str string) string {
	hash := sha256.New()
	hash.Write([]byte(str))
	md := hash.Sum(nil)
	encodedStr := hex.EncodeToString(md)
	return encodedStr
}
