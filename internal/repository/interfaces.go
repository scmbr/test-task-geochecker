package repository

import (
	"context"
	"time"

	"github.com/scmbr/test-task-geochecker/internal/domain"
	"github.com/scmbr/test-task-geochecker/internal/repository/models"
	"gorm.io/gorm"
)

type IncidentRepository interface {
	Create(ctx context.Context, incident *domain.Incident) error
	GetAll(ctx context.Context, offset, limit int) ([]*domain.Incident, uint32, error)
	GetById(ctx context.Context, id string) (*domain.Incident, error)
	Update(ctx context.Context, id string, input models.UpdateIncidentInput) error
	Delete(ctx context.Context, id string) error
	CountUniqueUsers(ctx context.Context, incidentID string, since time.Time) (int, error)
	FindNearbyIncidents(ctx context.Context, lat, lon float64, radius uint16) ([]*domain.Incident, error)
}
type CheckRepository interface {
	Create(ctx context.Context, check *domain.Check) error
	GetById(ctx context.Context, id string) (*domain.Check, error)
}
type OperatorRepository interface {
	GetActiveByAPIKeyHash(ctx context.Context, apiKeyHash string) (*domain.Operator, error)
	Create(ctx context.Context, operator *domain.Operator) error
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
