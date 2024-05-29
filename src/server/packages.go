package server

import (
	"crypto/sha512"
	"encoding/hex"
)

func hashPassword(password string, login string) string {
	hash := sha512.Sum512([]byte(login + password))
	hashedPassword := hex.EncodeToString(hash[:])

	return hashedPassword
}
