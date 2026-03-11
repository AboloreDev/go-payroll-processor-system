package helpers

import (
	"encoding/json"
	"net/http"
)

// writeJSON is a helper used by every handler
// It sets the Content-Type header and encodes
// any value as JSON in the response
func WriteJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}