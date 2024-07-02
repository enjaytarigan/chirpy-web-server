package security

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrTokenInvalid = errors.New("invalid token")
	ErrTokenExpired = errors.New("token has expired")
)

type JwtProvider struct {
	Issuer string
	secret []byte
}

func NewJwtProvider(secret []byte, issuer string) *JwtProvider {
	return &JwtProvider{
		secret: secret,
		Issuer: issuer,
	}
}

func (j *JwtProvider) GenerateToken(userID int, expiresIn time.Duration) (string, error) {
	claims := jwt.RegisteredClaims{
		Issuer:    j.Issuer,
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expiresIn)),
		Subject:   fmt.Sprintf("%d", userID),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	ss, err := token.SignedString(j.secret)

	if err != nil {
		log.Println(err)
		return "", fmt.Errorf("failed to sign jwt: %w", err)
	}

	return ss, nil
}

func (j *JwtProvider) VerifyToken(token string) (string, error) {
	rc := jwt.RegisteredClaims{}
	t, err := jwt.ParseWithClaims(token, &rc, func(t *jwt.Token) (interface{}, error) {
		return j.secret, nil
	})

	if err != nil {
		return "", ErrTokenInvalid
	}

	if !t.Valid {
		return "", ErrTokenInvalid
	}

	userID, err := t.Claims.GetSubject()

	if err != nil {
		return "", fmt.Errorf("could not get subject: %w", err)
	}

	issuer, err := t.Claims.GetIssuer()

	if err != nil {
		return "", fmt.Errorf("could not get issuer: %w", err)
	}

	if issuer != j.Issuer {
		return "", ErrTokenInvalid
	}

	return userID, nil
}
