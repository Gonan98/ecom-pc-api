package util

import (
	"encoding/json"
	"net/http"

	"github.com/gonan98/ecom-pc-api/internal/errors"
)

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

func WriteError(w http.ResponseWriter, err errors.APIError) error {
	return WriteJSON(w, err.Code, err)
}
