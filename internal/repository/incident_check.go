package repository

import (
	"context"
	"fmt"

	"github.com/scmbr/test-task-geochecker/internal/repository/models"
	"gorm.io/gorm"
)

type IncidentCheckRepo struct {
	db *gorm.DB
}

func NewIncidentCheckRepository(db *gorm.DB) *IncidentCheckRepo {
	return &IncidentCheckRepo{db: db}
}
func (r *IncidentCheckRepo) Create(ctx context.Context, incidentID, checkID string) error {
	if err := r.db.WithContext(ctx).Create(&models.IncidentsCheck{
		IncidentID: incidentID,
		CheckID:    checkID,
	}).Error; err != nil {
		return fmt.Errorf("repo.Create: %w", err)
	}
	return nil
}
