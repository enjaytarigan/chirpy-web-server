package main

import "net/http"

func (apiCfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	apiCfg.fileserverHits = 0
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
