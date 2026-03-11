package httpmux

import "net/http"

// ServeMux is a wrapper around http.ServeMux that supports middleware.
type ServeMux struct {
	*http.ServeMux
}

// NewServeMux creates a new ServeMux wrapper.
func NewServeMux() *ServeMux {
	return &ServeMux{ServeMux: http.NewServeMux()}
}

// Handle registers the handler for the given pattern with optional middleware.
func (s *ServeMux) Handle(pattern string, handler http.Handler, middleware ...Middleware) {
	s.ServeMux.Handle(pattern, UseHandler(handler, middleware...))
}

// HandleFunc registers the handler function for the given pattern with optional middleware.
func (s *ServeMux) HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request), middleware ...Middleware) {
	s.ServeMux.Handle(pattern, UseHandler(http.HandlerFunc(handler), middleware...))
}
