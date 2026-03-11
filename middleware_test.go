package httpmux

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCORSMiddleware(t *testing.T) {
	tests := []struct {
		name           string
		allowedOrigins []string
		requestOrigin  string
		expectedOrigin string
		expectedVary   string
	}{
		{
			name:           "Wildcard",
			allowedOrigins: []string{"*"},
			requestOrigin:  "http://example.com",
			expectedOrigin: "*",
			expectedVary:   "",
		},
		{
			name:           "Allowed Origin",
			allowedOrigins: []string{"http://example.com", "http://localhost:3000"},
			requestOrigin:  "http://example.com",
			expectedOrigin: "http://example.com",
			expectedVary:   "Origin",
		},
		{
			name:           "Disallowed Origin",
			allowedOrigins: []string{"http://localhost:3000"},
			requestOrigin:  "http://example.com",
			expectedOrigin: "",
			expectedVary:   "Origin",
		},
		{
			name:           "Empty Origin Header",
			allowedOrigins: []string{"http://localhost:3000"},
			requestOrigin:  "",
			expectedOrigin: "",
			expectedVary:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := CORSMiddleware(tt.allowedOrigins)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			}))

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			if tt.requestOrigin != "" {
				req.Header.Set("Origin", tt.requestOrigin)
			}
			rec := httptest.NewRecorder()

			handler.ServeHTTP(rec, req)

			assert.Equal(t, tt.expectedOrigin, rec.Header().Get("Access-Control-Allow-Origin"))
			assert.Equal(t, tt.expectedVary, rec.Header().Get("Vary"))
		})
	}
}

func TestCORSMiddleware_Options(t *testing.T) {
	called := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	})

	handler := CORSMiddleware([]string{"*"})(next)
	req := httptest.NewRequest(http.MethodOptions, "/", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.False(t, called)
	assert.Equal(t, http.StatusNoContent, rec.Code)
	assert.Equal(t, "GET, POST, PUT, DELETE, OPTIONS", rec.Header().Get("Access-Control-Allow-Methods"))
}
