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
	checkRepo repository.CheckRepository
}

func NewCheckService(checkRepo repository.CheckRepository) *CheckSvc {
	return &CheckSvc{checkRepo: checkRepo}
}
func (s *CheckSvc) Check(ctx context.Context, input *dto.CheckInput) error {
	err := s.checkRepo.Create(ctx, domain.Check{
		CheckID:   uuid.NewString(),
		UserID:    input.UserID,
		Longitude: input.Longitude,
		Latitude:  input.Latitude,
	})
	if err != nil {
		if errors.Is(err, repository.ErrAlreadyExists) {
			return ErrCheckAlreadyExists
		}
		return err
	}
	return nil
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
