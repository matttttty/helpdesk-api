package handler

import (
	"encoding/json"
	"log"
	"net/http"
)

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if data != nil {
		if err := json.NewEncoder(w).Encode(data); err != nil {
			log.Printf("failed to encode response: %v", err)
		}
	}
}

func writeError(w http.ResponseWriter, status int, msg string, internalErr error) {
	if internalErr != nil {
		log.Printf("error %d: %s: %v", status, msg, internalErr)
	}

	if status >= 500 {
		msg = "internal server error"
	}

	writeJSON(w, status, map[string]string{"error": msg})
}
