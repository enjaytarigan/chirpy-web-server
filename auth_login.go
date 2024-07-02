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

	payloadExpiresIn, err := time.ParseDuration(fmt.Sprintf("%ds", payload.ExpiresInSeconds))

	tokenExpiresIn := 24 * time.Hour // Default

	if err == nil && payload.ExpiresInSeconds != 0 && payloadExpiresIn.Hours() < 24.0 {
		tokenExpiresIn = payloadExpiresIn
	}

	token, err := api.jwt.GenerateToken(user.ID, tokenExpiresIn)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	responeBody := struct {
		ID    int    `json:"id"`
		Email string `json:"email"`
		Token string `json:"token"`
	}{
		ID:    user.ID,
		Email: user.Email,
		Token: token,
	}

	respondWithJSON(w, http.StatusOK, responeBody)
}
