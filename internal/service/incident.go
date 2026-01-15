package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/scmbr/test-task-geochecker/internal/domain"
	"github.com/scmbr/test-task-geochecker/internal/repository"
	"github.com/scmbr/test-task-geochecker/internal/repository/models"
	"github.com/scmbr/test-task-geochecker/internal/service/dto"
)

type IncidentSvc struct {
	incidentRepo repository.IncidentRepository
}

func NewIncidentService(incidentRepo repository.IncidentRepository) *IncidentSvc {
	return &IncidentSvc{incidentRepo: incidentRepo}
}
func (s *IncidentSvc) Create(ctx context.Context, input *dto.CreateIncidentInput) error {
	incident, err := domain.NewIncident(uuid.NewString(), input.OperatorID, input.Latitude, input.Longitude, uint16(input.Radius))
	if err != nil {
		return err
	}
	err = s.incidentRepo.Create(ctx, incident)

	if err != nil {
		return err
	}
	return nil
}
func (s *IncidentSvc) GetAll(ctx context.Context, input *dto.GetAllIncidentsInput) (*dto.GetAllIncidentsOutput, error) {
	res, total, err := s.incidentRepo.GetAll(ctx, input.Offset, input.Limit)
	if err != nil {
		return nil, err
	}
	incidents := make([]dto.GetIncidentOutput, len(res))
	for idx, i := range res {
		incidents[idx] = dto.GetIncidentOutput{
			ID:         i.IncidentID,
			OperatorID: i.OperatorID,
			Latitude:   i.Latitude,
			Longitude:  i.Longitude,
			Radius:     i.Radius,
			CreatedAt:  i.CreatedAt,
		}
	}
	return &dto.GetAllIncidentsOutput{
		Total:     total,
		Incidents: incidents}, nil
}
func (s *IncidentSvc) GetById(ctx context.Context, id string) (*dto.GetIncidentOutput, error) {
	res, err := s.incidentRepo.GetById(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrIncidentNotFound
		}
		return nil, err
	}
	return &dto.GetIncidentOutput{
		ID:         res.IncidentID,
		OperatorID: res.OperatorID,
		Latitude:   res.Latitude,
		Longitude:  res.Longitude,
		Radius:     res.Radius,
		CreatedAt:  res.CreatedAt,
	}, nil
}

func (s *IncidentSvc) Update(ctx context.Context, id string, input *dto.UpdateIncidentInput) error {
	incidentToUpdate := models.UpdateIncidentInput{}

	if input.OperatorID != nil {
		incidentToUpdate.OperatorID = input.OperatorID
	}
	if input.Latitude != nil {
		incidentToUpdate.Latitude = input.Latitude
	}
	if input.Longitude != nil {
		incidentToUpdate.Longitude = input.Longitude
	}
	if input.Radius != nil {
		incidentToUpdate.Radius = input.Radius
	}

	err := s.incidentRepo.Update(ctx, id, incidentToUpdate)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrIncidentNotFound
		}
		return err
	}

	return nil
}

func (s *IncidentSvc) Delete(ctx context.Context, id string) error {
	err := s.incidentRepo.Delete(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrIncidentNotFound
		}
		return err
	}
	return nil
}
func (s *IncidentSvc) GetStats(ctx context.Context, incidentID string, since time.Time) (int, error) {
	return s.incidentRepo.CountUniqueUsers(ctx, incidentID, since)
}
