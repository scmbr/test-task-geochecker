package models

import (
	"time"

	"github.com/scmbr/test-task-geochecker/internal/domain"
)

type Incident struct {
	IncidentID string     `gorm:"primaryKey;column:incident_id"`
	OperatorID string     `gorm:"column:operator_id;not null"`
	Longitude  float64    `gorm:"column:longitude;not null"`
	Latitude   float64    `gorm:"column:latitude;not null"`
	Radius     uint16     `gorm:"column:radius;not null"`
	CreatedAt  time.Time  `gorm:"column:created_at;autoCreateTime"`
	DeletedAt  *time.Time `gorm:"column:deleted_at;default:null"`
	UpdatedAt  *time.Time `gorm:"column:updated_at;default:null"`
}
type UpdateIncidentInput struct {
	OperatorID *string  `json:"operator_id"`
	Latitude   *float64 `json:"latitude"`
	Longitude  *float64 `json:"longitude"`
	Radius     *uint16  `json:"radius"`
}

func IncidentModelToDomain(m *Incident) *domain.Incident {
	return &domain.Incident{
		IncidentID: m.IncidentID,
		OperatorID: m.OperatorID,
		Longitude:  m.Longitude,
		Latitude:   m.Latitude,
		Radius:     uint16(m.Radius),
		CreatedAt:  m.CreatedAt,
		DeletedAt:  m.DeletedAt,
		UpdatedAt:  m.UpdatedAt,
	}
}

func IncidentDomainToModel(i *domain.Incident) *Incident {
	now := time.Now().UTC()
	return &Incident{
		IncidentID: i.IncidentID,
		OperatorID: i.OperatorID,
		Longitude:  i.Longitude,
		Latitude:   i.Latitude,
		Radius:     i.Radius,
		CreatedAt:  now,
		UpdatedAt:  &now,
	}
}
