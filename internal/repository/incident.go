package repository

import (
	"context"
	"fmt"
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
	incidentModel := models.IncidentDomainToModel(incident)
	if err := r.db.WithContext(ctx).Create(&incidentModel).Error; err != nil {
		if err == gorm.ErrDuplicatedKey {
			return ErrAlreadyExists
		}
		return fmt.Errorf("repo.Create: %w", err)
	}
	return nil
}
func (r *IncidentRepo) GetAll(ctx context.Context, offset, limit int) ([]*domain.Incident, uint32, error) {
	var incidents []models.Incident
	var total int64
	q := r.db.WithContext(ctx).Model(&models.Incident{}).Where("deleted_at IS NULL")

	if err := q.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("repo.GetAll: %w", err)
	}
	if err := q.Order("created_at DESC").Limit(limit).Offset(offset).Find(&incidents).Error; err != nil {
		return nil, 0, fmt.Errorf("repo.GetAll: %w", err)
	}
	incidentsDomain := make([]*domain.Incident, len(incidents))
	for idx := range incidentsDomain {
		var err error
		incidentsDomain[idx], err = models.IncidentModelToDomain(&incidents[idx])
		if err != nil {
			return nil, 0, err
		}
	}
	return incidentsDomain, uint32(total), nil
}
func (r *IncidentRepo) GetById(ctx context.Context, id string) (*domain.Incident, error) {
	var incident models.Incident
	res := r.db.WithContext(ctx).Where("incident_id = ?", id).First(&incident)
	if res.Error == gorm.ErrRecordNotFound {
		return nil, ErrNotFound
	}
	if res.Error != nil {
		return nil, fmt.Errorf("repo.GetById: %w", res.Error)
	}
	incidentDomain, err := models.IncidentModelToDomain(&incident)
	if err != nil {
		return nil, err
	}
	return incidentDomain, nil
}
func (r *IncidentRepo) Update(ctx context.Context, id string, input models.UpdateIncidentInput) error {
	updates := map[string]interface{}{
		"updated_at": time.Now(),
	}

	if input.OperatorID != nil {
		updates["operator_id"] = *input.OperatorID
	}

	if input.Latitude != nil && input.Longitude != nil {
		updates["location"] = models.PointWKT(*input.Longitude, *input.Latitude)
	}
	if input.Radius != nil {
		updates["radius"] = *input.Radius
	}

	res := r.db.WithContext(ctx).Model(&models.Incident{}).Where("incident_id = ?", id).Updates(updates)

	if res.RowsAffected == 0 {
		return ErrNotFound
	}
	if res.Error != nil {
		return fmt.Errorf("repo.Update: %w", res.Error)
	}

	return nil
}

func (r *IncidentRepo) Delete(ctx context.Context, id string) error {
	res := r.db.WithContext(ctx).Model(&models.Incident{}).Where("incident_id = ?", id).Update("deleted_at", time.Now())
	if res.RowsAffected == 0 {
		return ErrNotFound
	}
	if res.Error != nil {
		return fmt.Errorf("repo.Delete: %w", res.Error)
	}

	return nil
}
func (r *IncidentRepo) CountUniqueUsers(ctx context.Context, incidentID string, since time.Time) (int, error) {
	var count int64

	err := r.db.WithContext(ctx).
		Table("checks").
		Joins("JOIN incidents_checks ic ON ic.check_id = checks.check_id").
		Where("ic.incident_id = ?", incidentID).
		Where("checks.created_at >= ?", since).
		Select("COUNT(DISTINCT checks.user_id)").
		Scan(&count).
		Error

	if err != nil {
		return 0, fmt.Errorf("repo.CountUniqueUsers: %w", err)
	}

	return int(count), nil
}
func (r *IncidentRepo) FindNearbyIncidents(ctx context.Context, lat, lon float64, radius uint16) ([]*domain.Incident, error) {
	var incidents []*models.Incident
	point := models.PointWKT(lon, lat)
	radiusMeters := float64(radius)
	if err := r.db.WithContext(ctx).
		Where("ST_DWithin(location, ST_GeomFromText(?, 4326)::geography, ? + radius)", point, radiusMeters).
		Find(&incidents).Error; err != nil {
		return nil, fmt.Errorf("repo.FindNearbyIncidents: %w", err)
	}
	incidentsDomain := make([]*domain.Incident, len(incidents))
	for idx := range incidentsDomain {
		var err error
		incidentsDomain[idx], err = models.IncidentModelToDomain(incidents[idx])
		if err != nil {
			return nil, err
		}
	}
	return incidentsDomain, nil
}
