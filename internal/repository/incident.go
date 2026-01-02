package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/scmbr/test-task-geochecker/internal/domain/models"
	"gorm.io/gorm"
)

type IncidentRepo struct {
	db *gorm.DB
}

func NewIncidentRepository(db *gorm.DB) *IncidentRepo {
	return &IncidentRepo{db: db}
}

func (r *IncidentRepo) Create(ctx context.Context, incident models.Incident) error {
	if err := r.db.WithContext(ctx).Create(&incident).Error; err != nil {
		return err
	}
	return nil
}
func (r *IncidentRepo) GetAll(ctx context.Context, offset, limit int) ([]models.Incident, uint32, error) {
	var incidents []models.Incident
	var total int64
	if err := r.db.WithContext(ctx).Model(&models.Incident{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := r.db.WithContext(ctx).Limit(limit).Offset(offset).Order("created_at DESC").Find(&incidents).Error; err != nil {
		return nil, 0, err
	}
	return incidents, uint32(total), nil
}
func (r *IncidentRepo) GetById(ctx context.Context, id string) (*models.Incident, error) {
	var incident models.Incident
	res := r.db.WithContext(ctx).Where("incident_id = ?", id).First(&incident)
	if res.Error == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if res.Error != nil {
		return nil, res.Error
	}
	return &incident, nil
}
func (r *IncidentRepo) Update(ctx context.Context, id string, incident models.Incident) error {
	res := r.db.WithContext(ctx).Where("incident_id = ?", id).
		Updates(map[string]interface{}{
			"operator_id": incident.OperatorID,
			"latitude":    incident.Latitude,
			"longitude":   incident.Longitude,
			"radius":      incident.Radius,
		})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return fmt.Errorf("incident with id: %s not found", id)
	}
	return nil
}
func (r *IncidentRepo) Delete(ctx context.Context, id string) error {
	res := r.db.WithContext(ctx).Where("incident_id = ?", id).Update("deleted_at", time.Now())
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return fmt.Errorf("incident with id: %s not found", id)
	}
	return nil
}
func (r *IncidentRepo) CountUniqueUsers(ctx context.Context, incidentID string, since time.Time) (int, error) {
	var count int64
	if err := r.db.WithContext(ctx).Table("checks").Where("incident_id = ?", incidentID).Where("created_at >= ?", since).Distinct("user_id").Count(&count).Error; err != nil {
		return 0, err
	}
	return int(count), nil
}
