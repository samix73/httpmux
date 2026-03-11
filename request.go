package httpmux

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// DecodeJSON decodes a JSON request body into a struct.
// It uses DisallowUnknownFields() to prevent unknown fields from being decoded.
func DecodeJSON[T any](r *http.Request) (T, error) {
	var v T
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&v); err != nil {
		return v, fmt.Errorf("failed to decode request: %w", err)
	}

	return v, nil
}
