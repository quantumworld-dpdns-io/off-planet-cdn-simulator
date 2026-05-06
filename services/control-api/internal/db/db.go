package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Client wraps a pgxpool.Pool for database access.
type Client struct {
	Pool *pgxpool.Pool
}

// New creates a new database client, opens the connection pool, and pings the server.
func New(ctx context.Context, connString string) (*Client, error) {
	pool, err := pgxpool.New(ctx, connString)
	if err != nil {
		return nil, fmt.Errorf("create pgx pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("ping database: %w", err)
	}

	return &Client{Pool: pool}, nil
}

// Close releases all connections in the pool.
func (c *Client) Close() {
	c.Pool.Close()
}
