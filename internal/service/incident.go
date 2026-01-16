package service

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/scmbr/test-task-geochecker/internal/domain"
	"github.com/scmbr/test-task-geochecker/internal/repository"
	"github.com/scmbr/test-task-geochecker/internal/repository/models"
	"github.com/scmbr/test-task-geochecker/internal/service/dto"
	"github.com/scmbr/test-task-geochecker/pkg/cache"
	"github.com/scmbr/test-task-geochecker/pkg/logger"
)

type IncidentSvc struct {
	incidentRepo repository.IncidentRepository
	cache        cache.Cache
}

func NewIncidentService(incidentRepo repository.IncidentRepository, cache cache.Cache) *IncidentSvc {
	return &IncidentSvc{
		incidentRepo: incidentRepo,
		cache:        cache,
	}
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
	filterJSON, err := json.Marshal(input)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal filter: %w", err)
	}
	hash := sha256.Sum256(filterJSON)
	cacheKey := fmt.Sprintf("incidents:%x", hash)
	if cached, err := s.cache.Get(ctx, cacheKey); err == nil && cached != "" {
		var incidents dto.GetAllIncidentsOutput
		if jsonErr := json.Unmarshal([]byte(cached), &incidents); jsonErr == nil {
			return &incidents, nil
		} else {
			logger.Error("failed to unmarshal cached data for key ", jsonErr, map[string]interface{}{
				"cache_key": cacheKey,
			})
		}
	}
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
	payload, err := json.Marshal(dto.GetAllIncidentsOutput{
		Total:     total,
		Incidents: incidents})
	if err != nil {
		logger.Error("failed to marshal payload", err, nil)
	} else if err := s.cache.Set(ctx, cacheKey, string(payload), time.Minute*2); err != nil {
		logger.Error("failed to set cache", err, nil)
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
