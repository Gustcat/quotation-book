package quote

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Gustcat/quotation-book/internal/http-server/response"
	"github.com/Gustcat/quotation-book/internal/storage"
)

func List(w http.ResponseWriter, r *http.Request) {
	var author *string
	query := r.URL.Query()
	if authorQuery := query.Get("author"); authorQuery != "" {
		author = &authorQuery
	}

	quotes := storage.List(author)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(quotes); err != nil {
		log.Printf("Error encoding JSON: %s", err)
		response.Error("Failed to encode quotes", w, http.StatusInternalServerError)
		return
	}
}
