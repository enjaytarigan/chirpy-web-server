package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	var dir http.Dir

	mux.Handle("/", http.FileServer(dir))

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	log.Fatal(server.ListenAndServe())

	fmt.Printf("Starting server on localhost:8080")
}
