package auth

import (
	"crypto/rand"
	"encoding/hex"
)

func MakeRefreshToken() (string, error) {
	// Generate a random 32-byte token
	token := make([]byte, 32)
	_, err := rand.Read(token)
	if err != nil {
		return "", err
	}

	// Return the token as a hexadecimal string
	return hex.EncodeToString(token), nil
}
