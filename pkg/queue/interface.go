package queue

import "context"

type Queue interface {
	Enqueue(ctx context.Context, task Task) error
	Dequeue(ctx context.Context) (Task, error)
	Ack(ctx context.Context, task Task) error
	Nack(ctx context.Context, task Task, reason string) error
}
