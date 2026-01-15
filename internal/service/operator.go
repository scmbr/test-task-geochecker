package service

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/scmbr/test-task-geochecker/internal/domain"
	"github.com/scmbr/test-task-geochecker/internal/repository"
	"github.com/scmbr/test-task-geochecker/internal/service/dto"
	"github.com/scmbr/test-task-geochecker/pkg/hasher"
)

type OperatorSvc struct {
	operatorRepo repository.OperatorRepository
	apiKeySecret string
}

func NewOperatorService(operatorRepo repository.OperatorRepository, apiKeySecret string) *OperatorSvc {
	return &OperatorSvc{operatorRepo: operatorRepo, apiKeySecret: apiKeySecret}
}

func (s *OperatorSvc) Create(ctx context.Context, input *dto.CreateOperatorInput) error {
	hash := hasher.HashAPIKey(s.apiKeySecret, input.APIKey)

	operator, err := domain.NewOperator(uuid.NewString(), input.Name, string(hash))
	if err != nil {
		return err
	}
	if err := s.operatorRepo.Create(ctx, operator); err != nil {
		return err
	}
	return nil
}

func (s *OperatorSvc) ValidateAPIKey(ctx context.Context, apiKey string) (*dto.ValidateOperatorOutput, error) {
	apiKeyHash := hasher.HashAPIKey(s.apiKeySecret, apiKey)
	op, err := s.operatorRepo.GetActiveByAPIKeyHash(ctx, apiKeyHash)
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
