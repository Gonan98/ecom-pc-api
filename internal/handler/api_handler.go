package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gonan98/ecom-pc-api/internal/types"
)

var validate = validator.New(validator.WithRequiredStructEnabled())

type APIHandler func(http.ResponseWriter, *http.Request) error

func httpHandler(h APIHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := h(w, r); err != nil {
			if apiError, ok := err.(types.APIError); ok {
				writeError(w, apiError)
			} else {
				writeError(w, types.APIError{Code: http.StatusInternalServerError, Message: "Internal server error"})
			}
			slog.Error("HTTP API", "err", err.Error(), "path", r.URL.Path)
		}
	}
}

func readJSON(r *http.Request, payload any) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	return decoder.Decode(payload)
}

func writeJSON(w http.ResponseWriter, status int, data any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, apiErr types.APIError) error {
	return writeJSON(w, apiErr.Code, apiErr)
}

func writeResponse(w http.ResponseWriter, resp types.APIResponse) error {
	return writeJSON(w, resp.Code, resp)
}
