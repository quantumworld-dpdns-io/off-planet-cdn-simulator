package redis

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type Client struct {
	rdb *redis.Client
}

func New(addr string) *Client {
	rdb := redis.NewClient(&redis.Options{Addr: addr})
	return &Client{rdb: rdb}
}

func (c *Client) Ping(ctx context.Context) error {
	return c.rdb.Ping(ctx).Err()
}

func (c *Client) Close() error {
	return c.rdb.Close()
}

// EnqueuePreloadJob pushes a job ID to the right of the preload queue.
func (c *Client) EnqueuePreloadJob(ctx context.Context, jobID string) error {
	return c.rdb.RPush(ctx, "preload:jobs:pending", jobID).Err()
}

// DequeuePreloadJob pops a job ID from the left of the preload queue.
// Returns ("", nil) if queue is empty.
func (c *Client) DequeuePreloadJob(ctx context.Context) (string, error) {
	result, err := c.rdb.LPop(ctx, "preload:jobs:pending").Result()
	if err == redis.Nil {
		return "", nil
	}
	return result, err
}

// CancelPreloadJob stores a cancellation signal in a Redis set.
func (c *Client) CancelPreloadJob(ctx context.Context, jobID string) error {
	return c.rdb.SAdd(ctx, "preload:jobs:cancelled", jobID).Err()
}

// IsJobCancelled checks whether a job was cancelled.
func (c *Client) IsJobCancelled(ctx context.Context, jobID string) (bool, error) {
	return c.rdb.SIsMember(ctx, "preload:jobs:cancelled", jobID).Result()
}
