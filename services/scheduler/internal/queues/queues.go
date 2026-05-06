package queues

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

// Queue wraps a go-redis client for simple list-based job queuing.
type Queue struct {
	client *redis.Client
}

// New creates a Queue from a Redis URL.
func New(redisURL string) (*Queue, error) {
	opts, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, fmt.Errorf("parse redis url: %w", err)
	}
	return &Queue{client: redis.NewClient(opts)}, nil
}

// Enqueue pushes a JSON payload to the tail of the named list (RPUSH).
func (q *Queue) Enqueue(ctx context.Context, queueName, payload string) error {
	if err := q.client.RPush(ctx, queueName, payload).Err(); err != nil {
		return fmt.Errorf("enqueue %s: %w", queueName, err)
	}
	return nil
}

// Dequeue pops a payload from the head of the named list (LPOP).
// Returns ("", nil) when the queue is empty.
func (q *Queue) Dequeue(ctx context.Context, queueName string) (string, error) {
	val, err := q.client.LPop(ctx, queueName).Result()
	if err == redis.Nil {
		return "", nil
	}
	if err != nil {
		return "", fmt.Errorf("dequeue %s: %w", queueName, err)
	}
	return val, nil
}
