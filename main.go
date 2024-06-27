package main

import (
	"log"
	"net/http"
)

type apiConfig struct {
	fileserverHits int
}

func main() {
	port := "8080"
	apiCfg := &apiConfig{fileserverHits: 0}

	mux := http.NewServeMux()

	fsHandler := http.StripPrefix("/app", http.FileServer(http.Dir(".")))

	mux.Handle("/app/*", apiCfg.middlewareMetricsInc(fsHandler))

	mux.HandleFunc("GET /healthz", handlerReadiness)

	mux.HandleFunc("GET /metrics", apiCfg.handlerMetrics)

	mux.HandleFunc("GET /reset", apiCfg.handlerReset)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Fatal(server.ListenAndServe())
}
