package main

import (
	"log"
	"net/http"

	"github.com/Gustcat/quotation-book/internal/storage"

	"github.com/Gustcat/quotation-book/internal/http-server/handlers/quote"
	"github.com/gorilla/mux"
)

const (
	baseUrl = "localhost:8080"
)

func main() {
	r := mux.NewRouter()

	qb := storage.NewQBook()

	r.HandleFunc("/quotes", quote.Create(qb)).Methods("POST")
	r.HandleFunc("/quotes", quote.List(qb)).Methods("GET")
	r.HandleFunc("/quotes/random", quote.GetRandom(qb)).Methods("GET")
	r.HandleFunc("/quotes/{id}", quote.Delete(qb)).Methods("DELETE")

	log.Printf("Server listen on %s", baseUrl)

	err := http.ListenAndServe(baseUrl, r)
	if err != nil {
		log.Fatal(err)
	}
}
