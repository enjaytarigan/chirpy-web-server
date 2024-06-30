package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func handersValidateChirp(w http.ResponseWriter, r *http.Request) {
	var reqBody = struct {
		Body string `json:"body"`
	}{}

	err := json.NewDecoder(r.Body).Decode(&reqBody)

	if err != nil {
		sendJSON(http.StatusInternalServerError, map[string]string{
			"error": "Something went wrong",
		}, w)
		return
	}

	if len(reqBody.Body) > 140 {
		sendJSON(http.StatusBadRequest, map[string]string{
			"error": "Chirp is too long",
		}, w)
		return
	}

	responseBody := struct {
		Valid bool `json:"valid"`
	}{
		Valid: true,
	}
	sendJSON(http.StatusOK, responseBody, w)
}

func sendJSON(statusCode int, body any, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	err := json.NewEncoder(w).Encode(body)

	if err != nil {
		fmt.Printf("error sending response: %v", err)
	}
}
