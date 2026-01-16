package queue

import "context"

type Queue interface {
	Enqueue(ctx context.Context, task Task) error
	Dequeue(ctx context.Context, blockMS int64) (Task, string, error)
	Ack(ctx context.Context, taskID string) error
	Nack(ctx context.Context, taskID string, task Task, reason string) error
}
