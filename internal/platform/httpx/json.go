// Package httpx holds HTTP helpers shared by every module: JSON encoding,
// error responses and middleware.
package httpx

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
)

const maxBodyBytes = 1 << 20 // 1 MiB; uploads get their own route and limit

type errorBody struct {
	Error string `json:"error"`
}

func WriteJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		slog.Error("httpx: encode response", "err", err)
	}
}

func WriteError(w http.ResponseWriter, status int, msg string) {
	WriteJSON(w, status, errorBody{Error: msg})
}

// WriteInternalError logs the real cause and returns a generic message, so
// internal details never leak to the client.
func WriteInternalError(w http.ResponseWriter, r *http.Request, err error) {
	slog.ErrorContext(r.Context(), "internal error",
		"method", r.Method, "path", r.URL.Path, "err", err)
	WriteError(w, http.StatusInternalServerError, "internal error")
}

// DecodeJSON reads a single JSON object into dst, rejecting unknown fields
// and bodies above maxBodyBytes.
func DecodeJSON(w http.ResponseWriter, r *http.Request, dst any) error {
	r.Body = http.MaxBytesReader(w, r.Body, maxBodyBytes)

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	if err := dec.Decode(dst); err != nil {
		return fmt.Errorf("invalid JSON body: %w", err)
	}
	if dec.More() {
		return errors.New("invalid JSON body: multiple objects")
	}
	return nil
}
