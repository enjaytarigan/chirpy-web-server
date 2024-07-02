package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/enjaytarigan/chirpy-web-server/internal/database"
	"github.com/enjaytarigan/chirpy-web-server/internal/security"
)

const (
	UserUpgradeEventName = "user.upgraded"
)

func (api *apiConfig) handlerUserUpgradeEvent(w http.ResponseWriter, r *http.Request) {
	apiKey, found := strings.CutPrefix(r.Header.Get("Authorization"), "ApiKey ")

	if !found || !security.IsValidPolkaApiKey(apiKey) {
		respondWithError(w, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
		return
	}

	reqBody := struct {
		Event string `json:"event"`
		Data  struct {
			UserID int `json:"user_id"`
		} `json:"data"`
	}{}

	err := json.NewDecoder(r.Body).Decode(&reqBody)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "could not decode the request body")
		return
	}

	if reqBody.Event != UserUpgradeEventName {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	user, err := api.db.GetUserByID(reqBody.Data.UserID)
	if err != nil {
		if errors.Is(err, database.ErrUserNotFound) {
			respondWithError(w, http.StatusNotFound, err.Error())
			return
		}

		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	err = api.db.UpdateIsChirpyRedStatus(user.ID, true)

	if err != nil {
		if errors.Is(err, database.ErrUserNotFound) {
			respondWithError(w, http.StatusNotFound, err.Error())
			return
		}

		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
