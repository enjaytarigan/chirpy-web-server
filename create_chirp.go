package main

import (
	"encoding/json"
	"net/http"
)

func (api *apiConfig) handlerCreateChirp(w http.ResponseWriter, r *http.Request) {
	var reqBody = struct {
		Body string `json:"body"`
	}{}

	err := json.NewDecoder(r.Body).Decode(&reqBody)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	if len(reqBody.Body) > 140 {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}

	chirp, err := api.db.CreateChirp(cleanChirp(reqBody.Body))

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, chirp)
}
