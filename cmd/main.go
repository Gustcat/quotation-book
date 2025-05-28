package main

import (
	"github.com/Gustcat/quotation-book/internal/storage"
	"log"
	"net/http"

	"github.com/Gustcat/quotation-book/internal/http-server/handlers/quote"
	"github.com/gorilla/mux"
)

const (
	baseUrl = "localhost:8080"
)

func main() {
	r := mux.NewRouter()

	qb := storage.NewQBook()

	creator := func(w http.ResponseWriter, r *http.Request) {
		quote.Create(w, r, qb)
	}

	lister := func(w http.ResponseWriter, r *http.Request) {
		quote.List(w, r, qb)
	}

	deleter := func(w http.ResponseWriter, r *http.Request) {
		quote.Delete(w, r, qb)
	}

	randomizer := func(w http.ResponseWriter, r *http.Request) {
		quote.GetRandom(w, r, qb)
	}

	r.HandleFunc("/quotes", creator).Methods("POST")
	r.HandleFunc("/quotes", lister).Methods("GET")
	r.HandleFunc("/quotes/random", randomizer).Methods("GET")
	r.HandleFunc("/quotes/{id}", deleter).Methods("DELETE")

	log.Printf("Server listen on %s", baseUrl)

	err := http.ListenAndServe(baseUrl, r)
	if err != nil {
		log.Fatal(err)
	}
}
