package main

import (
	"context"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/enjaytarigan/chirpy-web-server/internal/database"
)

type TypeUserRefreshTokenKey string

var (
	UserRefreshTokenKey TypeUserRefreshTokenKey = "refresh_token"
)

func (api *apiConfig) WithRefreshToken(next http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")

		authToken, hasPrefix := strings.CutPrefix(authHeader, "Bearer ")

		if !hasPrefix {
			respondWithError(w, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
			return
		}

		user, err := api.db.GetUserByRefreshToken(authToken)

		if err != nil {
			respondWithError(w, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
			return
		}

		if user.RefreshToken.IsExpired() {
			respondWithError(w, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
			return
		}

		req := r.WithContext(context.WithValue(r.Context(), UserRefreshTokenKey, user))

		next.ServeHTTP(w, req)
	})
}

func (api *apiConfig) handlerRefreshToken(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(UserRefreshTokenKey).(database.User)

	token, err := api.jwt.GenerateToken(user.ID, 1*time.Hour)

	if err != nil {
		log.Println(err)
		respondWithError(w, http.StatusInternalServerError, "could not generate new access token")
	}

	respondWithJSON(w, http.StatusOK, map[string]any{
		"token": token,
	})
}

func (api *apiConfig) handlerRevokeRefreshToken(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(UserRefreshTokenKey).(database.User)

	err := api.db.SaveRefreshToken(user.ID, nil)

	if err != nil {
		log.Println(err)
		respondWithError(w, http.StatusInternalServerError, "could not revoke the refresh token")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
