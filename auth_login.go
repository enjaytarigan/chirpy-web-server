package main

import (
	"errors"
	"net/http"

	"github.com/enjaytarigan/chirpy-web-server/internal/database"
	"github.com/enjaytarigan/chirpy-web-server/internal/security"
)

func (api *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	payload := struct {
		Email    string `json:"email"`
		Password string `json:"password"`
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

	responeBody := struct {
		ID    int    `json:"id"`
		Email string `json:"email"`
	}{
		ID:    user.ID,
		Email: user.Email,
	}

	respondWithJSON(w, http.StatusOK, responeBody)
}
