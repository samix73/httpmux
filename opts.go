package httpmux

import (
	"os"
	"time"
)

// Option is a function that configures server options.
type Option func(*opts)

type opts struct {
	shutdownTimeout    time.Duration
	readTimeout        time.Duration
	writeTimeout       time.Duration
	idleTimeout        time.Duration
	readHeaderTimeout  time.Duration
	terminationSignals []os.Signal
}

// WithShutdownTimeout sets the timeout for shutting down the server.
func WithShutdownTimeout(timeout time.Duration) Option {
	return func(o *opts) {
		o.shutdownTimeout = timeout
	}
}

// WithReadTimeout sets the maximum duration for reading the entire request, including the body.
func WithReadTimeout(timeout time.Duration) Option {
	return func(o *opts) {
		o.readTimeout = timeout
	}
}

// WithWriteTimeout sets the maximum duration before timing out writes of the response.
func WithWriteTimeout(timeout time.Duration) Option {
	return func(o *opts) {
		o.writeTimeout = timeout
	}
}

// WithIdleTimeout sets the maximum amount of time to wait for the next request when keep-alives are enabled.
func WithIdleTimeout(timeout time.Duration) Option {
	return func(o *opts) {
		o.idleTimeout = timeout
	}
}

// WithReadHeaderTimeout sets the amount of time allowed to read request headers.
func WithReadHeaderTimeout(timeout time.Duration) Option {
	return func(o *opts) {
		o.readHeaderTimeout = timeout
	}
}

// WithTerminationSignals sets the OS signals that will trigger a graceful server shutdown.
func WithTerminationSignals(signals ...os.Signal) Option {
	return func(o *opts) {
		o.terminationSignals = signals
	}
}
