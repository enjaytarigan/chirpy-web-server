package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/enjaytarigan/chirpy-web-server/internal/database"
	"github.com/enjaytarigan/chirpy-web-server/internal/security"
)

func (api *apiConfig) handlerUpdateUser(w http.ResponseWriter, r *http.Request) {
	body := struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}{}

	err := json.NewDecoder(r.Body).Decode(&body)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}

	hashedPassword, err := security.GenerateHash([]byte(body.Password))

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	ctx := r.Context()
	userID, _ := strconv.Atoi(ctx.Value(UserLoggedInKey).(string))

	updatedUser, err := api.db.UpdateUser(userID, database.UpdateUserIn{
		Email:    body.Email,
		Password: hashedPassword,
	})

	if err != nil {
		if errors.Is(err, database.ErrEmailAlreadyRegistered) {
			respondWithError(w, http.StatusConflict, database.ErrEmailAlreadyRegistered.Error())
			return
		}

		if errors.Is(err, database.ErrUserNotFound) {
			respondWithError(w, http.StatusNotFound, err.Error())
			return
		}

		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	resp := struct {
		ID    int    `json:"id"`
		Email string `json:"email"`
	}{
		ID:    updatedUser.ID,
		Email: updatedUser.Email,
	}

	respondWithJSON(w, http.StatusOK, resp)
}
