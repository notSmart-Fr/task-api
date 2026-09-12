package response

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

// JSON sends a JSON payload with the specified HTTP status code.
func JSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)

	if data != nil {
		encoder := json.NewEncoder(w)
		encoder.SetEscapeHTML(false) // Prevents converting & to \u0026
		_ = encoder.Encode(data)
	}
}

// Error sends a standardized JSON error message.
func Error(w http.ResponseWriter, status int, message string) {
	JSON(w, status, ErrorResponse{Error: message})
}
