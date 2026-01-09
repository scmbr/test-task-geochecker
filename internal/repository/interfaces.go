package repository

import (
	"context"
	"time"

	"github.com/scmbr/test-task-geochecker/internal/domain/models"
	"gorm.io/gorm"
)

type IncidentRepository interface {
	Create(ctx context.Context, incident models.Incident) error
	GetAll(ctx context.Context, offset, limit int) ([]models.Incident, uint32, error)
	GetById(ctx context.Context, id string) (*models.Incident, error)
	Update(ctx context.Context, id string, input models.UpdateIncidentInput) error
	Delete(ctx context.Context, id string) error
	CountUniqueUsers(ctx context.Context, incidentID string, since time.Time) (int, error)
}
type CheckRepository interface {
	Create(ctx context.Context, check models.Check) error
	GetById(ctx context.Context, id string) (*models.Check, error)
}
type OperatorRepository interface {
	GetActiveByAPIKeyHash(ctx context.Context, apiKeyHash string) (*models.Operator, error)
	Create(ctx context.Context, operator models.Operator) error
	Revoke(ctx context.Context, operatorID string) error
}
type Repository struct {
	Incident IncidentRepository
	Check    CheckRepository
	Operator OperatorRepository
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		Incident: NewIncidentRepository(db),
		Check:    NewCheckRepository(db),
		Operator: NewOperatorRepository(db),
	}
}
