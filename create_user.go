package main

import (
	"net/http"
)

func (api *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	body := struct {
		Email string `json:"email"`
	}{}

	if err := decodeJSON(r, &body); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	user, err := api.db.CreateUser(body.Email)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, user)
}
