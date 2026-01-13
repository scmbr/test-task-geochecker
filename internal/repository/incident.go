package repository

import (
	"context"
	"time"

	"github.com/scmbr/test-task-geochecker/internal/domain"
	"github.com/scmbr/test-task-geochecker/internal/repository/models"
	"gorm.io/gorm"
)

type IncidentRepo struct {
	db *gorm.DB
}

func NewIncidentRepository(db *gorm.DB) *IncidentRepo {
	return &IncidentRepo{db: db}
}

func (r *IncidentRepo) Create(ctx context.Context, incident *domain.Incident) error {
	if err := r.db.WithContext(ctx).Create(&incident).Error; err != nil {
		if err == gorm.ErrDuplicatedKey {
			return ErrAlreadyExists
		}
		return err
	}
	return nil
}
func (r *IncidentRepo) GetAll(ctx context.Context, offset, limit int) ([]domain.Incident, uint32, error) {
	var incidents []domain.Incident
	var total int64
	q := r.db.WithContext(ctx).Model(&domain.Incident{}).Where("deleted_at IS NULL")

	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := q.Order("created_at DESC").Limit(limit).Offset(offset).Find(&incidents).Error; err != nil {
		return nil, 0, err
	}

	return incidents, uint32(total), nil
}
func (r *IncidentRepo) GetById(ctx context.Context, id string) (*domain.Incident, error) {
	var incident domain.Incident
	res := r.db.WithContext(ctx).Where("incident_id = ?", id).First(&incident)
	if res.Error == gorm.ErrRecordNotFound {
		return nil, ErrNotFound
	}
	if res.Error != nil {
		return nil, res.Error
	}
	return &incident, nil
}
func (r *IncidentRepo) Update(ctx context.Context, id string, input models.UpdateIncidentInput) error {
	updates := map[string]interface{}{
		"updated_at": time.Now(),
	}

	if input.OperatorID != nil {
		updates["operator_id"] = *input.OperatorID
	}
	if input.Latitude != nil {
		updates["latitude"] = *input.Latitude
	}
	if input.Longitude != nil {
		updates["longitude"] = *input.Longitude
	}
	if input.Radius != nil {
		updates["radius"] = *input.Radius
	}

	res := r.db.WithContext(ctx).Where("incident_id = ?", id).Updates(updates)

	if res.RowsAffected == 0 {
		return ErrNotFound
	}
	if res.Error != nil {
		return res.Error
	}

	return nil
}

func (r *IncidentRepo) Delete(ctx context.Context, id string) error {
	res := r.db.WithContext(ctx).Where("incident_id = ?", id).Update("deleted_at", time.Now())
	if res.RowsAffected == 0 {
		return ErrNotFound
	}
	if res.Error != nil {
		return res.Error
	}

	return nil
}
func (r *IncidentRepo) CountUniqueUsers(ctx context.Context, incidentID string, since time.Time) (int, error) {
	var count int64
	res := r.db.WithContext(ctx).Table("checks").Where("incident_id = ?", incidentID).Where("created_at >= ?", since).Distinct("user_id").Count(&count)
	if res.Error != nil {
		return 0, res.Error
	}
	return int(count), nil
}
