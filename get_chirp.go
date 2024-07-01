package main

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/enjaytarigan/chirpy-web-server/internal/database"
)

func (api *apiConfig) handlerGetChirpByID(w http.ResponseWriter, r *http.Request) {
	chirpID, err := strconv.Atoi(r.PathValue("chirpId"))

	if err != nil {
		respondWithError(w, http.StatusNotFound, "Chirp Not Found")
		return
	}

	chirp, err := api.db.GetChirpByID(chirpID)

	if err != nil {
		if errors.Is(err, database.ErrChirpNotFound) {
			respondWithError(w, http.StatusNotFound, err.Error())
		} else {
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	respondWithJSON(w, http.StatusOK, chirp)
}
