package quote

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/Gustcat/quotation-book/internal/http-server/response"
	"github.com/Gustcat/quotation-book/internal/storage"
)

type Randomizer interface {
	GetRandom() (*storage.QuoteWithID, error)
}

func GetRandom(w http.ResponseWriter, r *http.Request, randomizer Randomizer) {
	w.Header().Set("Content-Type", "application/json")

	quote, err := randomizer.GetRandom()

	if errors.Is(err, storage.ErrQuoteNotFound) {
		log.Printf("Quote not found: %s", err)
		response.Error("Quote not found", w, http.StatusNotFound)
		return
	}

	if err != nil {
		log.Printf("Error getting random quote: %s", err)
		response.Error("Error getting random quote", w, http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(quote); err != nil {
		log.Printf("Error encoding JSON: %s", err)
		response.Error("Failed to encode quote", w, http.StatusInternalServerError)
		return
	}
}
