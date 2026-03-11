package httpmux

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestNewServer_Timeouts(t *testing.T) {
	readTimeout := 5 * time.Second
	writeTimeout := 10 * time.Second
	idleTimeout := 120 * time.Second
	readHeaderTimeout := 2 * time.Second

	srv := NewServer(":8080",
		WithReadTimeout(readTimeout),
		WithWriteTimeout(writeTimeout),
		WithIdleTimeout(idleTimeout),
		WithReadHeaderTimeout(readHeaderTimeout),
	)

	require.Equal(t, readTimeout, srv.httpServer.ReadTimeout)
	require.Equal(t, writeTimeout, srv.httpServer.WriteTimeout)
	require.Equal(t, idleTimeout, srv.httpServer.IdleTimeout)
	require.Equal(t, readHeaderTimeout, srv.httpServer.ReadHeaderTimeout)
}
