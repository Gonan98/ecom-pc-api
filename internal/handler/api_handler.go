package handler

import (
	"log/slog"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gonan98/ecom-pc-api/internal/errors"
	"github.com/gonan98/ecom-pc-api/internal/util"
)

var validate = validator.New(validator.WithRequiredStructEnabled())

type APIHandler func(http.ResponseWriter, *http.Request) error

func HttpHandler(h APIHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := h(w, r); err != nil {
			if apiError, ok := err.(errors.APIError); ok {
				util.WriteError(w, apiError)
			} else {
				util.WriteError(w, errors.APIError{Code: http.StatusInternalServerError, Message: "Internal server error"})
			}
			slog.Error("HTTP API", "err", err.Error(), "path", r.URL.Path)
		}
	}
}
