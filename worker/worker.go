package worker

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/scmbr/test-task-geochecker/pkg/logger"
	"github.com/scmbr/test-task-geochecker/pkg/queue"
)

type Worker struct {
	queue queue.Queue
	ctx   context.Context
}

func NewWorker(ctx context.Context, q queue.Queue) *Worker {
	return &Worker{
		queue: q,
		ctx:   ctx,
	}
}

func (w *Worker) Run() {
	for {

		task, taskID, err := w.queue.Dequeue(w.ctx, 1000)
		if err != nil {
			log.Printf("failed to dequeue task: %v", err)
			time.Sleep(time.Second)
			continue
		}

		if taskID == "" {
			continue
		}

		err = w.handleTask(task)
		if err != nil {
			reason := err.Error()
			task.Attempts++
			if err := w.queue.Nack(w.ctx, taskID, task, reason); err != nil {
				logger.Error("failed to nack task", err, map[string]interface{}{
					"task_id": taskID,
				})
			}
			continue
		}

		if err := w.queue.Ack(w.ctx, taskID); err != nil {
			logger.Error("failed to ack task", err, map[string]interface{}{
				"task_id": taskID,
			})
		}
	}
}

func (w *Worker) handleTask(task queue.Task) error {
	req, err := http.NewRequestWithContext(w.ctx, http.MethodPost, task.TargetURL, bytes.NewBuffer([]byte(task.Payload)))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		logger.Error("failed to handle task", fmt.Errorf("status %d", resp.StatusCode), map[string]interface{}{
			"task_id":    task.TaskID,
			"target_url": task.TargetURL,
		})
		return fmt.Errorf("webhook returned status %d", resp.StatusCode)
	}

	return nil
}
