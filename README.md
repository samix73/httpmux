# httpmux

This package is a wrapper around `net/http` to provide additional features and a more ergonomic API for building HTTP servers in Go.

## Features

- **Middleware Support:** Easy-to-use middleware chains for the server and specific route groups (`Use`, `Group`).
- **Route Groups:** Group routes under a common prefix and apply middleware locally.
- **Graceful Shutdown:** Built-in signal handling and coordinated shutdown of your server and background tasks.
- **Request Decoding:** Strongly-typed JSON request decoding with `Decode` and `DecodeJSON`.
- **Response Writing:** Helper functions for writing standard JSON responses (`JSON`) and structured errors (`Error`).
- **Trace ID Support:** Built-in context handling for request trace IDs.
- **Configurable Timeouts:** Sensible defaults for read, write, idle, and shutdown timeouts, customizable via functional options.
- **CORS Middleware:** Built-in middleware for handling Cross-Origin Resource Sharing.

## Usage

### Creating a Server

```go
package main

import (
	"log/slog"
	"net/http"

	"github.com/samix73/httpmux"
)

func main() {
	// Create a new server on port 8080
	s := httpmux.NewServer(":8080")

	// Apply global middleware
	s.Use(httpmux.CORSMiddleware([]string{"*"}))

	// Register a simple handler
	s.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		httpmux.JSON(w, http.StatusOK, map[string]string{"message": "pong"})
	})

	// Start the server
	if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		slog.Error("server failed", "error", err)
	}
}
```

### Route Groups and Middleware

You can group routes under a common prefix and apply specific middleware to that group only.

```go
s.Group("/api/v1", func(mux *httpmux.ServeMux) {
	mux.HandleFunc("/users", listUsersHandler)
	mux.HandleFunc("/users/{id}", getUserHandler)
}, authMiddleware) // authMiddleware applies only to /api/v1/*
```

### Request Decoding & Response Encoding

The package provides handy generic helpers for decoding JSON requests and encoding JSON responses.

```go
type CreateUserRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func createUserHandler(w http.ResponseWriter, r *http.Request) {
	// Decode JSON request body into struct (disallows unknown fields by default)
	req, err := httpmux.DecodeJSON[CreateUserRequest](r)
	if err != nil {
		httpmux.Error(r.Context(), w, http.StatusBadRequest, "invalid request body")
		return
	}

	// ... process request ...
	
	// Write standard JSON response
	user := map[string]string{"id": "123", "name": req.Name}
	httpmux.JSON(w, http.StatusCreated, user)
}
```

### Graceful Shutdown

The `Server` makes it easy to handle OS termination signals properly and register shutdown hooks for cleanup (like database closing).

```go
import (
	"context"
	"time"
)

// Initialize server with customized timeouts
s := httpmux.NewServer(":8080", httpmux.WithShutdownTimeout(15 * time.Second))

// Register cleanup tasks
s.AtShutdown(func(ctx context.Context) error {
	slog.Info("cleaning up resources...")
	// e.g. return db.Close()
	return nil
})

// Run the server in a goroutine
go func() {
	if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		slog.Error("server failed", "error", err)
	}
}()

// Block until a termination signal (SIGINT, SIGTERM) is received
<-s.NotifyTermination()

// Initiate graceful shutdown
if err := s.Shutdown(context.Background()); err != nil {
	slog.Error("shutdown failed", "error", err)
}
```
