package queues

import (
	"context"
	"github.com/redis/go-redis/v9"
)

const (
	PendingJobsQueue = "preload:jobs:pending"
	CancelledJobsSet = "preload:jobs:cancelled"
)

type Client struct {
	rdb *redis.Client
}

func New(addr string) *Client {
	rdb := redis.NewClient(&redis.Options{Addr: addr})
	return &Client{rdb: rdb}
}

func (c *Client) Close() error {
	return c.rdb.Close()
}

func (c *Client) Ping(ctx context.Context) error {
	return c.rdb.Ping(ctx).Err()
}

// Enqueue pushes a job ID onto the pending queue.
func (c *Client) Enqueue(ctx context.Context, jobID string) error {
	return c.rdb.RPush(ctx, PendingJobsQueue, jobID).Err()
}

// Dequeue pops a job ID from the pending queue. Returns ("", nil) if empty.
func (c *Client) Dequeue(ctx context.Context) (string, error) {
	result, err := c.rdb.LPop(ctx, PendingJobsQueue).Result()
	if err == redis.Nil {
		return "", nil
	}
	return result, err
}

// IsCancelled returns true if the job was cancelled via the cancellation set.
func (c *Client) IsCancelled(ctx context.Context, jobID string) (bool, error) {
	return c.rdb.SIsMember(ctx, CancelledJobsSet, jobID).Result()
}

// QueueLength returns the current length of the pending jobs queue.
func (c *Client) QueueLength(ctx context.Context) (int64, error) {
	return c.rdb.LLen(ctx, PendingJobsQueue).Result()
}
