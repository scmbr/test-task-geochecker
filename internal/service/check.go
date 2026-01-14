package service

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/scmbr/test-task-geochecker/internal/domain"
	"github.com/scmbr/test-task-geochecker/internal/repository"
	"github.com/scmbr/test-task-geochecker/internal/service/dto"
)

type CheckSvc struct {
	checkRepo    repository.CheckRepository
	incidentRepo repository.IncidentRepository
	radius       uint16
}

func NewCheckService(checkRepo repository.CheckRepository, incidentRepo repository.IncidentRepository, radius uint16) *CheckSvc {
	return &CheckSvc{checkRepo: checkRepo, incidentRepo: incidentRepo, radius: radius}
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
