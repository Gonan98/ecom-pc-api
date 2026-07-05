package util

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gonan98/ecom-pc-api/internal/errors"
)

var Validate = validator.New(validator.WithRequiredStructEnabled())

type APIHandler func(http.ResponseWriter, *http.Request) error

func HttpHandler(h APIHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := h(w, r); err != nil {
			if apiError, ok := err.(errors.APIError); ok {
				writeError(w, apiError)
			} else {
				writeError(w, errors.APIError{Code: http.StatusInternalServerError, Message: "Internal server error"})
			}
			slog.Error("HTTP API", "err", err.Error(), "path", r.URL.Path)
		}
	}
}

func ReadJSON(r *http.Request, payload any) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	return decoder.Decode(payload)
}

func WriteJSON(w http.ResponseWriter, status int, data any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, err errors.APIError) error {
	return WriteJSON(w, err.Code, err)
}
