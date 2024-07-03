package main

import (
	"net/http"

	"github.com/enjaytarigan/chirpy-web-server/internal/database"
)

func (api *apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	query := database.QueryGetChirps{
		AuthorID: q.Get("author_id"),
		Sort:     q.Get("sort"),
	}

	chirps, err := api.db.GetChirps(query)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, chirps)
}
