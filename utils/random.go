package utils

import (
	"crypto/rand"
	"io"
	mr "math/rand"
	"strconv"
	"time"
)

var digitsAndNumbers = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

// GenerateRandomCode generate random string
func GenerateRandomCode(codeLength int) string {
	code := make([]rune, codeLength)
	r := mr.NewSource(time.Now().UnixMicro())
	rand := mr.New(r)

	for i := range code {
		code[i] = digitsAndNumbers[rand.Intn(len(digitsAndNumbers))]
	}
	return string(code)
}

//GenerateRandomFileName genrates the fileName with unique time
func GenerateRandomFileName() string {
	time := time.Now().UnixNano()
	return strconv.FormatInt(time, 10)
}

// GenerateRandomDigitSequence generates random digit sequence
func GenerateRandomDigitSequence(max int) string {
	b := make([]byte, max)
	n, err := io.ReadAtLeast(rand.Reader, b, max)
	if n != max {
		panic(err)
	}
	for i := 0; i < len(b); i++ {
		b[i] = table[int(b[i])%len(table)]
	}
	return string(b)
}

func GenerateRandomEmail(length int) string {
	return GenerateRandomCode(length) + "@gmail.com"
}

var table = [...]byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}
