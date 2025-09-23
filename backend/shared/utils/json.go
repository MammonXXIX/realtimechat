package utils

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/go-playground/validator/v10"
)

type APIResponse struct {
	Data  any       `json:"data,omitempty"`
	Error *APIError `json:"error,omitempty"`
}

type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

var Validate *validator.Validate

func init() {
	Validate = validator.New(validator.WithRequiredStructEnabled())
}

func ReadJSON(w http.ResponseWriter, r *http.Request, data any) error {
	r.Body = http.MaxBytesReader(w, r.Body, 1_048_576)
	defer r.Body.Close()

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(data); err != nil {
		return err
	}

	if decoder.More() {
		return errors.New("BODY MUST ONLY CONTAIN A SINGLE JSON VALUE")
	}

	return nil
}

func WriteJSON(w http.ResponseWriter, status int, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(data)
}

func DecodeJSON(r io.Reader, data any) error {
	decoder := json.NewDecoder(r)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(data); err != nil {
		return err
	}

	if decoder.More() {
		return errors.New("BODY MUST ONLY CONTAIN A SINGLE JSON VALUE")
	}

	return nil
}
