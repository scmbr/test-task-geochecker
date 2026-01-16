package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisQueue struct {
	client           *redis.Client
	streamKey        string
	deadLetterStream string
	group            string
	consumer         string
	maxAttempts      uint
}

func NewRedisQueue(client *redis.Client, streamKey, group, consumer string, maxAttempts uint, deadLetterStream string) *RedisQueue {
	return &RedisQueue{
		client:           client,
		streamKey:        streamKey,
		group:            group,
		consumer:         consumer,
		maxAttempts:      maxAttempts,
		deadLetterStream: deadLetterStream,
	}
}

func (r *RedisQueue) Enqueue(ctx context.Context, task Task) error {
	taskJSON, err := json.Marshal(task)
	if err != nil {
		return fmt.Errorf("failed to marshal task: %w", err)
	}
	_, err = r.client.XAdd(ctx, &redis.XAddArgs{
		Stream: r.streamKey,
		Values: map[string]interface{}{"task": string(taskJSON)},
		ID:     "*",
	}).Result()
	if err != nil {
		return fmt.Errorf("failed to enqueue task: %w", err)
	}
	return nil
}

func (r *RedisQueue) Dequeue(ctx context.Context, blockMS int64) (Task, string, error) {
	var empty Task

	streams, err := r.client.XReadGroup(ctx, &redis.XReadGroupArgs{
		Group:    r.group,
		Consumer: r.consumer,
		Streams:  []string{r.streamKey, ">"},
		Count:    1,
		Block:    time.Duration(blockMS) * time.Millisecond,
	}).Result()
	if err != nil {
		if err == redis.Nil {
			return empty, "", nil
		}
		return empty, "", fmt.Errorf("failed to read from stream: %w", err)
	}

	msg := streams[0].Messages[0]
	dataStr, ok := msg.Values["task"].(string)
	if !ok {
		return empty, "", fmt.Errorf("task field missing or invalid")
	}

	var task Task
	if err := json.Unmarshal([]byte(dataStr), &task); err != nil {
		return empty, "", fmt.Errorf("failed to unmarshal task: %w", err)
	}

	return task, msg.ID, nil
}

func (r *RedisQueue) Ack(ctx context.Context, taskID string) error {
	_, err := r.client.XAck(ctx, r.streamKey, r.group, taskID).Result()
	if err != nil {
		return fmt.Errorf("failed to ack task: %w", err)
	}
	return nil
}

func (r *RedisQueue) Nack(ctx context.Context, taskID string, task Task, reason string) error {
	task.Attempts++

	if task.Attempts < r.maxAttempts {

		time.Sleep(100 * time.Millisecond)
		taskJSON, _ := json.Marshal(task)
		_, err := r.client.XAdd(ctx, &redis.XAddArgs{
			Stream: r.streamKey,
			Values: map[string]interface{}{"task": string(taskJSON)},
			ID:     "*",
		}).Result()
		if err != nil {
			return fmt.Errorf("failed to requeue task: %w", err)
		}
	} else {

		taskJSON, _ := json.Marshal(task)
		_, err := r.client.XAdd(ctx, &redis.XAddArgs{
			Stream: r.deadLetterStream,
			Values: map[string]interface{}{"task": string(taskJSON)},
			ID:     "*",
		}).Result()
		if err != nil {
			return fmt.Errorf("failed to send task to dead-letter queue: %w", err)
		}
	}

	if _, err := r.client.XAck(ctx, r.streamKey, r.group, taskID).Result(); err != nil {
		return fmt.Errorf("failed to ack task after nack: %w", err)
	}

	return nil
}
