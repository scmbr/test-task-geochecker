package models

import (
	"time"

	"github.com/scmbr/test-task-geochecker/internal/domain"
)

type Check struct {
	CheckID   string    `gorm:"primaryKey;column:check_id"`
	UserID    string    `gorm:"column:user_id;not null"`
	Location  string    `gorm:"column:location;type:geometry(POINT,4326);not null"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
}

func CheckDomainToModel(c *domain.Check) *Check {
	return &Check{
		CheckID:   c.CheckID,
		UserID:    c.UserID,
		Location:  PointWKT(c.Longitude, c.Latitude),
		CreatedAt: c.CreatedAt,
	}
}
func CheckModelToDomain(m *Check) (*domain.Check, error) {
	lon, lat, err := ParseEWKBPoint(m.Location)
	if err != nil {
		return nil, err
	}
	return &domain.Check{
		CheckID:   m.CheckID,
		UserID:    m.UserID,
		Longitude: lon,
		Latitude:  lat,
		CreatedAt: m.CreatedAt,
	}, nil
}
