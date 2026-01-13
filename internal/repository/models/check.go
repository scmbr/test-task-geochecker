package models

import (
	"time"

	"github.com/scmbr/test-task-geochecker/internal/domain"
)

type Check struct {
	CheckID   string    `gorm:"primaryKey;column:check_id"`
	UserID    string    `gorm:"column:user_id;not null"`
	Longitude float64   `gorm:"column:longitude;not null"`
	Latitude  float64   `gorm:"column:latitude;not null"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
}

func CheckDomainToModel(c *domain.Check) *Check {
	return &Check{
		CheckID:   c.CheckID,
		UserID:    c.UserID,
		Longitude: c.Longitude,
		Latitude:  c.Latitude,
		CreatedAt: c.CreatedAt,
	}
}
func CheckModelToDomain(m *Check) *domain.Check {
	return &domain.Check{
		CheckID:   m.CheckID,
		UserID:    m.UserID,
		Longitude: m.Longitude,
		Latitude:  m.Latitude,
		CreatedAt: m.CreatedAt,
	}
}
