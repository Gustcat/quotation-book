package response

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

func Error(msg string, w http.ResponseWriter, status int) {
	var buf bytes.Buffer

	if err := json.NewEncoder(&buf).Encode(ErrorResponse{Error: msg}); err != nil {
		log.Printf("%v", err)
		http.Error(w, "invalid error json", http.StatusInternalServerError)
	} else {
		w.WriteHeader(status)
		if _, err := buf.WriteTo(w); err != nil {
			log.Printf("%v", err)
		}
	}
}
