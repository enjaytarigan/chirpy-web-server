package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("receiving request")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Not Found"))
	})

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	log.Fatal(server.ListenAndServe())

	fmt.Printf("Starting server on localhost:8080")
}
