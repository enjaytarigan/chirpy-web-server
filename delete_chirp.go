package main

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/enjaytarigan/chirpy-web-server/internal/database"
)

func (api *apiConfig) handlerDeleteChirp(w http.ResponseWriter, r *http.Request) {
	chirpID, err := strconv.Atoi(r.PathValue("chirpId"))

	if err != nil {
		respondWithError(w, http.StatusNotFound, database.ErrChirpNotFound.Error())
		return
	}

	chirp, err := api.db.GetChirpByID(chirpID)

	if err != nil {
		if errors.Is(err, database.ErrChirpNotFound) {
			respondWithError(w, http.StatusNotFound, database.ErrChirpNotFound.Error())
			return
		} else {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
	}

	userID, _ := strconv.Atoi(r.Context().Value(UserLoggedInKey).(string))

	if userID != chirp.AuthorID {
		respondWithError(w, http.StatusForbidden, http.StatusText(http.StatusForbidden))
		return
	}

	err = api.db.DeleteChirpByID(chirpID)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
