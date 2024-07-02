package security

import (
	"crypto/rand"
	"encoding/hex"
)

func GenerateRandToken(length int) (string, error) {
	token := make([]byte, length)

	_, err := rand.Read(token)

	if err != nil {
		return "", err
	}

	return hex.EncodeToString(token), nil
}
