package main

import (
	"github.com/Gustcat/quotation-book/internal/http-server/handlers/quote"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

const (
	baseUrl = "localhost:8080"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/quotes", quote.Create).Methods("POST")
	r.HandleFunc("/quotes", quote.List).Methods("GET")
	r.HandleFunc("/quotes/random", quote.GetRandom).Methods("GET")
	r.HandleFunc("/quotes/{id}", quote.Delete).Methods("DELETE")

	log.Printf("Server listen on %s", baseUrl)

	err := http.ListenAndServe(baseUrl, r)
	if err != nil {
		log.Fatal(err)
	}
}
