package security

import (
	"errors"
	"log"

	"golang.org/x/crypto/bcrypt"
)

func GenerateHash(password []byte) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword(password, bcrypt.MinCost)

	if err != nil {
		log.Printf("GenerateHash: %v", err)
		return "", err
	}

	return string(hashed), nil
}

func Compare(hashedPassword []byte, plainPassword []byte) bool {
	err := bcrypt.CompareHashAndPassword(hashedPassword, plainPassword)

	if err != nil && !errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		log.Printf("Compare: %v", err)
	}

	return err == nil
}
