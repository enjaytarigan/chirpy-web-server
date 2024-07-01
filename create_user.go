package main

import (
	"net/http"

	"github.com/enjaytarigan/chirpy-web-server/internal/security"
)

func (api *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	body := struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}{}

	if err := decodeJSON(r, &body); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	if api.db.CheckUserExist(body.Email) {
		respondWithError(w, http.StatusConflict, "user has been registered")
		return
	}

	hashedPassword, err := security.GenerateHash([]byte(body.Password))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	user, err := api.db.CreateUser(body.Email, hashedPassword)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, user)
}
