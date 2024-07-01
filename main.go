package main

import (
	"log"
	"net/http"

	"github.com/enjaytarigan/chirpy-web-server/internal/database"
)

type apiConfig struct {
	fileserverHits int
	db             *database.DB
}

func main() {
	port := "8080"

	db, err := database.NewDB("./database.json")

	if err != nil {
		log.Fatalf("failed connecting to db: %v", err)
	}

	apiCfg := &apiConfig{fileserverHits: 0, db: db}

	mux := http.NewServeMux()

	fsHandler := http.StripPrefix("/app", http.FileServer(http.Dir(".")))

	mux.Handle("/app/", apiCfg.middlewareMetricsInc(fsHandler))

	mux.HandleFunc("GET /api/healthz", handlerReadiness)

	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)

	mux.HandleFunc("GET /api/reset", apiCfg.handlerReset)

	mux.HandleFunc("POST /api/chirps", apiCfg.handlerCreateChirp)
	mux.HandleFunc("GET /api/chirps", apiCfg.handlerGetChirps)
	mux.HandleFunc("GET /api/chirps/{chirpId}", apiCfg.handlerGetChirpByID)

	mux.HandleFunc("POST /api/users", apiCfg.handlerCreateUser)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Fatal(server.ListenAndServe())
}
