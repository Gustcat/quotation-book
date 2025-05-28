package quote

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/Gustcat/quotation-book/internal/http-server/response"
	"github.com/Gustcat/quotation-book/internal/storage"
)

const (
	AuthorMax int = 50
	QuoteMax  int = 200
)

type createResponse struct {
	ID int64 `json:"id"`
}

type Creator interface {
	Create(quote *storage.Quote) (int64, error)
}

func Create(w http.ResponseWriter, r *http.Request, creator Creator) {
	w.Header().Set("Content-Type", "application/json")

	quote := &storage.Quote{}

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(quote); err != nil {
		log.Printf("Error decoding JSON: %s", err)
		response.Error(fmt.Sprintf("Invalid JSON - %s", err), w, http.StatusBadRequest)
		return
	}

	if len(quote.Author) > AuthorMax || len(quote.Author) < 1 {
		log.Printf("invalid field 'author' - %s", quote.Author)
		response.Error(fmt.Sprintf("field 'author' must be and contain from 1 to %d characters", AuthorMax),
			w, http.StatusBadRequest)
		return
	}

	if len(quote.Quote) > QuoteMax || len(quote.Quote) < 1 {
		log.Printf("invalid field 'quote' - %s", quote.Quote)
		response.Error(fmt.Sprintf("field 'quote' must be and contain from 1 to %d characters", QuoteMax),
			w, http.StatusBadRequest)
		return
	}

	id, err := creator.Create(quote)

	if errors.Is(err, storage.ErrQuoteExists) {
		log.Printf("Quote '%s':'%s' already exists", quote.Author, quote.Quote)
		response.Error("Quote already exists", w, http.StatusConflict)
		return
	}

	if err != nil {
		log.Printf("Error creating quote: %s", err)
		response.Error("Failed to create quote", w, http.StatusInternalServerError)
		return
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(createResponse{ID: id}); err != nil {
		log.Printf("Error encoding JSON: %s", err)
		response.Error("Failed to encode quote", w, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	_, err = buf.WriteTo(w)
	if err != nil {
		log.Printf("Error writing response: %s", err)
	}
}
