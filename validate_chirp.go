package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func cleanChirp(chirp string) string {
	var cleanedChirp strings.Builder

	words := strings.Split(chirp, " ")

	for i, word := range words {
		switch strings.ToLower(word) {
		case "kerfuffle", "sharbert", "fornax": // Profone words
			cleanedChirp.WriteString("****")
		default:
			cleanedChirp.WriteString(word)
		}

		if isLastIndex := i == len(words)-1; !isLastIndex {
			cleanedChirp.WriteByte(' ')
		}
	}

	return cleanedChirp.String()
}

func handersValidateChirp(w http.ResponseWriter, r *http.Request) {
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

	responseBody := struct {
		CleanedBody string `json:"cleaned_body"`
	}{
		CleanedBody: cleanChirp(reqBody.Body),
	}

	respondWithJSON(w, http.StatusOK, responseBody)
}

func respondWithError(w http.ResponseWriter, statusCode int, msg string) {
	errResp := struct {
		Error string `json:"error"`
	}{
		Error: msg,
	}

	respondWithJSON(w, statusCode, errResp)
}

func respondWithJSON(w http.ResponseWriter, statusCode int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	err := json.NewEncoder(w).Encode(payload)

	if err != nil {
		fmt.Printf("error sending response: %v", err)
	}
}
