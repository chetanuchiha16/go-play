package api

import (
	"encoding/json"
	"net/http"
)

// WriteJSON marshals v as JSON and writes it with the given status code.
func WriteJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}

// WriteError writes an RFC 7807 Problem Details JSON error response.
func WriteError(w http.ResponseWriter, status int, title, detail string) {
	WriteJSON(w, status, ErrorResponse{
		Status: &status,
		Title:  &title,
		Detail: &detail,
	})
}
