package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

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

func decodeJSON(r *http.Request, v any) error {
	err := json.NewDecoder(r.Body).Decode(v)

	if err != nil {
		return fmt.Errorf("failed unmarshalling request body: %v", err)
	}

	return nil
}
