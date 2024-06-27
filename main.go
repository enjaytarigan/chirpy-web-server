package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	filePathDir := "."
	port := "8080"
	mux := http.NewServeMux()

	mux.Handle("/app/*", http.StripPrefix("/app", http.FileServer(http.Dir(filePathDir))))

	mux.HandleFunc("/healthz", handlerReadiness)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Fatal(server.ListenAndServe())

	fmt.Printf("Starting server on localhost:8080")
}

func handlerReadiness(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.Write([]byte("OK"))
}
