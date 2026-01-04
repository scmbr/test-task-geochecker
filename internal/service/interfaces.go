package service

import (
	"context"
	"time"

	"github.com/scmbr/test-task-geochecker/internal/repository"
	"github.com/scmbr/test-task-geochecker/internal/service/dto"
)

type IncidentService interface {
	Create(ctx context.Context, input *dto.CreateIncidentInput) error
	GetAll(ctx context.Context, offset, limit int) (*dto.GetAllIncidentOutput, error)
	GetById(ctx context.Context, id string) (*dto.GetIncidentOutput, error)
	Update(ctx context.Context, id string, input *dto.UpdateIncidentInput) error
	Delete(ctx context.Context, id string) error
	GetStats(ctx context.Context, incidentID string, since time.Time) (int, error)
}
type CheckService interface {
	Check(ctx context.Context, input *dto.CheckInput) error
	GetById(ctx context.Context, id string) (*dto.GetCheckOutput, error)
}
type Service struct {
	Incident IncidentService
	Check    CheckService
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		Incident: NewIncidentService(repo.Incident),
		Check:    NewCheckService(repo.Check),
	}
}
