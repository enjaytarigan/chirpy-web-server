package main

import (
	"log"
	"net/http"
	"os"

	"github.com/enjaytarigan/chirpy-web-server/internal/database"
	"github.com/enjaytarigan/chirpy-web-server/internal/security"
	"github.com/joho/godotenv"
)

const (
	JWT_ISSUER = "chirpy"
)

type apiConfig struct {
	fileserverHits int
	db             *database.DB
	jwt            *security.JwtProvider
}

func main() {
	port := "8080"

	err := godotenv.Load()

	if err != nil {
		log.Fatalln("error loading .env file")
	}

	jwtProvider := security.NewJwtProvider(
		[]byte(os.Getenv("JWT_SECRET")),
		JWT_ISSUER,
	)

	db, err := database.NewDB("./database.json")

	if err != nil {
		log.Fatalf("failed connecting to db: %v", err)
	}

	apiCfg := &apiConfig{fileserverHits: 0, db: db, jwt: jwtProvider}

	mux := http.NewServeMux()

	fsHandler := http.StripPrefix("/app", http.FileServer(http.Dir(".")))

	mux.Handle("/app/", apiCfg.middlewareMetricsInc(fsHandler))

	mux.HandleFunc("GET /api/healthz", handlerReadiness)

	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)

	mux.HandleFunc("GET /api/reset", apiCfg.handlerReset)

	mux.Handle("POST /api/chirps", apiCfg.WithAuth(apiCfg.handlerCreateChirp))
	mux.HandleFunc("GET /api/chirps", apiCfg.handlerGetChirps)
	mux.HandleFunc("GET /api/chirps/{chirpId}", apiCfg.handlerGetChirpByID)

	mux.HandleFunc("POST /api/users", apiCfg.handlerCreateUser)
	mux.Handle("PUT /api/users", apiCfg.WithAuth(apiCfg.handlerUpdateUser))
	mux.HandleFunc("POST /api/login", apiCfg.handlerLogin)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Fatal(server.ListenAndServe())
}
