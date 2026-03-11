package httpmux

import (
	"net/http"
	"slices"
)

// Middleware represents a function that wraps an http.Handler with additional behavior.
type Middleware func(next http.Handler) http.Handler

// CORSMiddleware returns a middleware that adds CORS headers to responses.
func CORSMiddleware(allowedOrigins []string) Middleware {
	allowAll := slices.Contains(allowedOrigins, "*")
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if allowAll {
				w.Header().Set("Access-Control-Allow-Origin", "*")
			} else {
				origin := r.Header.Get("Origin")
				if origin != "" {
					if slices.Contains(allowedOrigins, origin) {
						w.Header().Set("Access-Control-Allow-Origin", origin)
					}
					w.Header().Set("Vary", "Origin")
				}
			}

			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Project-ID")
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// Use registers a middleware to the server.
func (s *Server) Use(m ...Middleware) {
	s.httpServer.Handler = UseHandler(s.httpServer.Handler, m...)
}

// UseHandler wraps a handler with a list of middleware.
func UseHandler(handler http.Handler, m ...Middleware) http.Handler {
	for i := len(m) - 1; i >= 0; i-- {
		handler = m[i](handler)
	}

	return handler
}
