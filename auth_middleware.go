package main

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/enjaytarigan/chirpy-web-server/internal/security"
)

type AuthAuthenticatedUserKey string

const (
	UserLoggedInKey AuthAuthenticatedUserKey = "user"
)

func (api *apiConfig) WithAuth(next http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")

		authToken, hasPrefix := strings.CutPrefix(authHeader, "Bearer ")

		if !hasPrefix {
			respondWithError(w, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
			return
		}

		userID, err := api.jwt.VerifyToken(authToken)

		if errors.Is(err, security.ErrTokenExpired) || errors.Is(err, security.ErrTokenInvalid) {
			respondWithError(w, http.StatusUnauthorized, err.Error())
			return
		}

		ctx := r.Context()
		req := r.WithContext(context.WithValue(ctx, UserLoggedInKey, userID))

		next.ServeHTTP(w, req)
	})
}
