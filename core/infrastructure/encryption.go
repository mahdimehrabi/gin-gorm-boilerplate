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

//do not change this function if you want your users can use their password(no change after first deployment)
func (e Encryption) SaltPassword(password string, salt string) string {
	if salt == "" {
		salt = e.env.Secret
	}
	return salt[:3] + password + salt[3:]
}

func (e Encryption) Sha256Encrypt(password string) string {
	hash := sha256.New()
	hash.Write([]byte(password))
	md := hash.Sum(nil)
	encodedStr := hex.EncodeToString(md)
	return encodedStr
}

func (e Encryption) SaltAndSha256Encrypt(password string) string {
	password = e.SaltPassword(password, "")
	return e.Sha256Encrypt(password)
}
