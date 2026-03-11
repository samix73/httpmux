package httpmux

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

// Default timeouts for the HTTP server.
const (
	DefaultShutdownTimeout   = 30 * time.Second
	DefaultReadTimeout       = 10 * time.Second
	DefaultWriteTimeout      = 10 * time.Second
	DefaultIdleTimeout       = 15 * time.Second
	DefaultReadHeaderTimeout = 10 * time.Second
)

// Server is a wrapper around http.Server that provides additional features.
type Server struct {
	httpServer *http.Server
	router     *ServeMux
	opts       *opts
	shutdownFn []func(ctx context.Context) error
}

// NewServer creates a new server.
func NewServer(serverAddress string, options ...Option) *Server {
	o := opts{
		shutdownTimeout:   DefaultShutdownTimeout,
		readTimeout:       DefaultReadTimeout,
		writeTimeout:      DefaultWriteTimeout,
		idleTimeout:       DefaultIdleTimeout,
		readHeaderTimeout: DefaultReadHeaderTimeout,
		terminationSignals: []os.Signal{
			syscall.SIGINT,
			syscall.SIGTERM,
		},
	}
	for _, opt := range options {
		opt(&o)
	}

	router := NewServeMux()

	return &Server{
		httpServer: &http.Server{
			Addr:              serverAddress,
			Handler:           router,
			ReadTimeout:       o.readTimeout,
			WriteTimeout:      o.writeTimeout,
			IdleTimeout:       o.idleTimeout,
			ReadHeaderTimeout: o.readHeaderTimeout,
		},
		router: router,
		opts:   &o,
	}
}

// ListenAndServe starts the server. It will block until the server is shut down.
func (s *Server) ListenAndServe() error {
	slog.Info("server.Server.ListenAndServe", slog.String("address", s.httpServer.Addr))
	return s.httpServer.ListenAndServe()
}

// AtShutdown registers a function to be called when the server is shut down.
func (s *Server) AtShutdown(fn ...func(ctx context.Context) error) {
	s.shutdownFn = append(s.shutdownFn, fn...)
}

// NotifyTermination returns a channel that will receive termination signals.
func (s *Server) NotifyTermination() <-chan os.Signal {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, s.opts.terminationSignals...)

	return sigChan
}

// Shutdown shuts down the server gracefully.
func (s *Server) Shutdown(ctx context.Context) error {
	slog.Info("server.Server.Shutdown signal received")
	ctx, cancel := context.WithTimeout(ctx, s.opts.shutdownTimeout)
	defer cancel()

	for _, fn := range s.shutdownFn {
		if err := fn(ctx); err != nil {
			slog.Error("server.Server.Shutdown fn", slog.String("error", err.Error()))
		}
	}

	if err := s.httpServer.Shutdown(ctx); err != nil {
		return fmt.Errorf("server.Server.Shutdown httpServer.Shutdown error: %w", err)
	}

	slog.Info("Shutdown completed")

	return nil
}

// Group creates a new router group with the given prefix and middleware.
func (s *Server) Group(prefix string, fn func(s *ServeMux), middleware ...Middleware) {
	mux := NewServeMux()

	prefix = strings.TrimSuffix(prefix, "/")

	s.Handle(prefix+"/", http.StripPrefix(prefix, UseHandler(mux, middleware...)))
	fn(mux)
}

// Handle registers a handler for the given pattern.
func (s *Server) Handle(pattern string, handler http.Handler) {
	s.router.Handle(pattern, handler)
}

// HandleFunc registers a handler function for the given pattern.
func (s *Server) HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	s.router.HandleFunc(pattern, handler)
}
