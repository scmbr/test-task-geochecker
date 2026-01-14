package models

import (
	"fmt"
	"time"

	"github.com/scmbr/test-task-geochecker/internal/domain"
)

type Incident struct {
	IncidentID string     `gorm:"primaryKey;column:incident_id"`
	OperatorID string     `gorm:"column:operator_id;not null"`
	Location   string     `gorm:"column:location;type:geometry(POINT,4326);not null"`
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

func IncidentModelToDomain(m *Incident) (*domain.Incident, error) {
	lon, lat, err := ParsePointWKT(m.Location)
	if err != nil {
		return nil, err
	}
	return &domain.Incident{
		IncidentID: m.IncidentID,
		OperatorID: m.OperatorID,
		Longitude:  lon,
		Latitude:   lat,
		Radius:     m.Radius,
		CreatedAt:  m.CreatedAt,
		DeletedAt:  m.DeletedAt,
		UpdatedAt:  m.UpdatedAt,
	}, nil
}

func IncidentDomainToModel(i *domain.Incident) *Incident {
	now := time.Now().UTC()
	return &Incident{
		IncidentID: i.IncidentID,
		OperatorID: i.OperatorID,
		Location:   PointWKT(i.Longitude, i.Latitude),
		Radius:     i.Radius,
		CreatedAt:  now,
		UpdatedAt:  &now,
	}
}
func PointWKT(lon, lat float64) string {
	return fmt.Sprintf("SRID=4326;POINT(%f %f)", lon, lat)
}
func ParsePointWKT(wkt string) (lon, lat float64, err error) {
	_, err = fmt.Sscanf(wkt, "SRID=4326;POINT(%f %f)", &lon, &lat)
	return
}
