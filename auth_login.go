package main

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/enjaytarigan/chirpy-web-server/internal/database"
	"github.com/enjaytarigan/chirpy-web-server/internal/security"
)

func (api *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	payload := struct {
		Email            string `json:"email"`
		Password         string `json:"password"`
		ExpiresInSeconds int    `json:"expires_in_seconds"`
	}{}

	err := decodeJSON(r, &payload)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	user, err := api.db.GetUserByEmail(payload.Email)

	if err != nil {
		if errors.Is(err, database.ErrUserNotFound) {
			respondWithError(w, http.StatusNotFound, err.Error())
		} else {
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	isMatch := security.Compare([]byte(user.Password), []byte(payload.Password))

	if !isMatch {
		respondWithError(w, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
		return
	}

	token, err := api.jwt.GenerateToken(user.ID, 1*time.Hour)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	refreshToken, err := security.GenerateRandToken(32)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	err = api.db.SaveRefreshToken(user.ID, &database.UserRefreshToken{
		Token:           refreshToken,
		ExpiriationTime: time.Now().Add(60 * (24 * time.Hour)),
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("could not save refresh token: %v", err))
		return
	}

	responeBody := struct {
		ID           int    `json:"id"`
		Email        string `json:"email"`
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
		IsChirpyRed  bool   `json:"is_chirpy_red"`
	}{
		ID:           user.ID,
		Email:        user.Email,
		Token:        token,
		RefreshToken: refreshToken,
		IsChirpyRed:  user.IsChirpyRed,
	}

	respondWithJSON(w, http.StatusOK, responeBody)
}
