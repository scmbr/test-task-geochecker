package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/scmbr/test-task-geochecker/internal/domain"
	"github.com/scmbr/test-task-geochecker/internal/repository"
	"github.com/scmbr/test-task-geochecker/internal/service/dto"
	"github.com/scmbr/test-task-geochecker/pkg/logger"
	"github.com/scmbr/test-task-geochecker/pkg/queue"
)

type CheckSvc struct {
	checkRepo         repository.CheckRepository
	incidentRepo      repository.IncidentRepository
	incidentCheckRepo repository.IncidentCheckRepository
	radius            uint16
	queue             queue.Queue
	webhookURL        string
}

func NewCheckService(checkRepo repository.CheckRepository, incidentRepo repository.IncidentRepository, incidentCheckRepo repository.IncidentCheckRepository, radius uint16, queue queue.Queue, webhookURL string) *CheckSvc {
	return &CheckSvc{
		checkRepo:         checkRepo,
		incidentRepo:      incidentRepo,
		incidentCheckRepo: incidentCheckRepo,
		radius:            radius,
		queue:             queue,
		webhookURL:        webhookURL,
	}
}
func (s *CheckSvc) Check(ctx context.Context, input *dto.CheckInput) ([]*dto.GetIncidentOutput, error) {
	check, err := domain.NewCheck(uuid.NewString(), input.UserID, input.Latitude, input.Longitude)
	if err != nil {
		return nil, err
	}
	err = s.checkRepo.Create(ctx, check)
	if err != nil {
		return nil, err
	}
	res, err := s.incidentRepo.FindNearbyIncidents(ctx, input.Latitude, input.Longitude, s.radius)
	if err != nil {
		return nil, err
	}
	if len(res) > 0 {
		for _, i := range res {
			payload, err := json.Marshal(i)
			if err != nil {
				return nil, fmt.Errorf("service.Check: %w", err)
			}
			taskID := uuid.NewString()
			go func(task queue.Task) {
				if err := s.queue.Enqueue(ctx, task); err != nil {
					logger.Error("failed to put task in queue: %w", err, map[string]interface{}{
						"task_id":    task.TaskID,
						"payload":    task.Payload,
						"target_url": task.TargetURL,
						"attempts":   task.Attempts,
					})
				}
			}(queue.Task{
				TaskID:    taskID,
				Payload:   string(payload),
				TargetURL: s.webhookURL,
				Attempts:  0,
			})
		}
	}
	for _, incident := range res {
		if err := s.incidentCheckRepo.Create(ctx, incident.IncidentID, check.CheckID); err != nil {
			return nil, err
		}
	}
	incidents := make([]*dto.GetIncidentOutput, len(res))
	for idx, i := range res {
		incidents[idx] = &dto.GetIncidentOutput{
			ID:         i.IncidentID,
			OperatorID: i.OperatorID,
			Latitude:   i.Latitude,
			Longitude:  i.Longitude,
			Radius:     i.Radius,
			CreatedAt:  i.CreatedAt,
		}
	}
	return incidents, nil
}
func (s *CheckSvc) GetById(ctx context.Context, id string) (*dto.GetCheckOutput, error) {
	check, err := s.checkRepo.GetById(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrCheckNotFound
		}
	}
	return &dto.GetCheckOutput{
		ID:        check.CheckID,
		UserID:    check.UserID,
		Latitude:  check.Latitude,
		Longitude: check.Longitude,
	}, nil
}
