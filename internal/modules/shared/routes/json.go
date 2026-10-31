package routes

import (
	"encoding/json"
	"io"
	"net/http"
)

func DecodeJSON[T any](r *http.Request) (T, error) {
	var value T

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return value, err
	}

	if err := json.Unmarshal(body, &value); err != nil {
		return value, err
	}

	return value, nil
}

func WriteJSON(
	w http.ResponseWriter,
	status int,
	data any,
) error {
	body, err := json.Marshal(data)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	_, err = w.Write(body)
	return err
}
