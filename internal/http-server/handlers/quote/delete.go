package quote

import (
	"errors"
	"github.com/Gustcat/quotation-book/internal/http-server/response"
	"github.com/Gustcat/quotation-book/internal/storage"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

func Delete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	idStr := mux.Vars(r)["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		log.Printf("Invalid id parameter: %s", err)
		response.Error("Invalid id parameter", w, http.StatusBadRequest)
		return
	}

	err = storage.Delete(id)

	if errors.Is(err, storage.ErrQuoteNotFound) {
		log.Printf("Quote not found: %s", err)
		response.Error("Quote not found", w, http.StatusNotFound)
		return
	}

	if err != nil {
		log.Printf("Error deleting quote: %s", err)
		response.Error("Error deleting quote", w, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
