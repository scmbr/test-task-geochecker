package service

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/scmbr/test-task-geochecker/internal/domain"
	"github.com/scmbr/test-task-geochecker/internal/repository"
	"github.com/scmbr/test-task-geochecker/internal/service/dto"
	"golang.org/x/crypto/bcrypt"
)

type OperatorSvc struct {
	operatorRepo repository.OperatorRepository
}

func NewOperatorService(operatorRepo repository.OperatorRepository) *OperatorSvc {
	return &OperatorSvc{operatorRepo: operatorRepo}
}

func (s *OperatorSvc) Create(ctx context.Context, input *dto.CreateOperatorInput) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(input.APIKey), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	operator, err := domain.NewOperator(uuid.NewString(), input.Name, string(hash))
	if err != nil {
		return err
	}
	if err := s.operatorRepo.Create(ctx, operator); err != nil {
		if errors.Is(err, repository.ErrAlreadyExists) {
			return ErrOperatorAlreadyExists
		}
		return err
	}
	return nil
}

func (s *OperatorSvc) ValidateAPIKey(ctx context.Context, apiKey string) (*dto.ValidateOperatorOutput, error) {
	op, err := s.operatorRepo.GetActiveByAPIKeyHash(ctx, apiKey)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrInvalidAPIKey
		}
		return nil, err
	}

	return &dto.ValidateOperatorOutput{
		OperatorID: op.OperatorID,
		Name:       op.Name,
		CreatedAt:  op.CreatedAt,
		RevokedAt:  op.RevokedAt,
	}, nil
}

func (s *OperatorSvc) Revoke(ctx context.Context, operatorID string) error {
	err := s.operatorRepo.Revoke(ctx, operatorID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrOperatorNotFound
		}
		return err
	}
	return nil
}
