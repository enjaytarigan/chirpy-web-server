package main

import "net/http"

func (api *apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {
	chirps, err := api.db.GetChirps()

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, chirps)
}
