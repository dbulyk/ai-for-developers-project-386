package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

// APIError is the error payload returned to clients.
type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// errorResponse wraps APIError in the envelope expected by the API contract.
type errorResponse struct {
	Error APIError `json:"error"`
}

// writeError writes a JSON error response with the given HTTP status.
func writeError(w http.ResponseWriter, status int, code, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	resp := errorResponse{Error: APIError{Code: code, Message: message}}
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		slog.Default().Error("failed to encode error response", slog.Any("error", err))
	}
}

// writeJSON writes a JSON response with the given HTTP status.
func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(payload); err != nil {
		slog.Default().Error("failed to encode json response", slog.Any("error", err))
	}
}
