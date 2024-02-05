package util

import (
	"crypto/rand"
	"encoding/base64"

	"golang.org/x/crypto/scrypt"
)

func GenerateSalt() string {
	salt := make([]byte, 32)
	rand.Read(salt)
	return base64.StdEncoding.EncodeToString(salt)
}

func HashPassword(password, salt string) string {
	key, _ := scrypt.Key([]byte(password), []byte(salt), 32768, 8, 1, 32)
	return base64.StdEncoding.EncodeToString(key)
}
