package infrastructure

import (
	"crypto/sha256"
	"encoding/hex"
)

//Encryption -> Encryption Struct
type Encryption struct {
	logger Logger
	env    Env
}

//NewEncryption -> return new Encryption struct
func NewEncryption(
	logger Logger,
	env Env,
) Encryption {
	return Encryption{
		logger: logger,
		env:    env,
	}
}

func (e Encryption) SaltPassword(password string) string {
	return password + e.env.Secret
}

func (e Encryption) Sha256Encrypt(password string) string {
	hash := sha256.New()
	hash.Write([]byte(password))
	md := hash.Sum(nil)
	encodedStr := hex.EncodeToString(md)
	return encodedStr
}
