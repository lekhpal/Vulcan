package transport

import (
	"context"
	"time"
)

// Config holds common configuration for connections.
type Config struct {
	Timeout       time.Duration
	BufferSize    int
	// Add more common configuration parameters here as needed
}

// Connection defines the interface for managing network connections.
type Connection interface {
	Connect(ctx context.Context) error
	Disconnect(ctx context.Context) error
	Send(ctx context.Context, data []byte) error
	Receive(ctx context.Context) ([]byte, error)
	IsConnected() bool
}

// Middleware defines an interface for intercepting and modifying connection operations.
type Middleware interface {
	Handle(ctx context.Context, conn Connection, next Connection) (Connection, error)
}

// Client defines a high-level interface for interacting with a remote service
// using the underlying connection and middleware.
type Client interface {
	Call(ctx context.Context, method string, request []byte) ([]byte, error)
	// Add more client-specific methods here as needed
}


