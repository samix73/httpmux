package httpmux

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type traceIDCtxKey struct{}

// TraceIDFromContext returns the trace ID from the context.
func TraceIDFromContext(ctx context.Context) string {
	return ctx.Value(traceIDCtxKey{}).(string)
}

// ContextWithTraceID returns a new context with the given trace ID.
func ContextWithTraceID(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, traceIDCtxKey{}, traceID)
}

// ErrorResponse represents an HTTP error response body.
type ErrorResponse struct {
	Error   string `json:"error"`
	TraceID string `json:"trace_id,omitempty"`
}

// JSON writes a JSON response to the ResponseWriter with the given status code.
func JSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

// Error writes a JSON error response with the given status code and message.
// It includes the trace ID from the context.
func Error(ctx context.Context, w http.ResponseWriter, status int, msg string) {
	resp := ErrorResponse{
		Error:   msg,
		TraceID: TraceIDFromContext(ctx),
	}
	JSON(w, status, resp)
}

// Decode unmarshals a JSON request body into a new instance of type T.
func Decode[T any](r *http.Request) (T, error) {
	var v T
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		return v, fmt.Errorf("failed to decode request: %w", err)
	}

	return v, nil
}
