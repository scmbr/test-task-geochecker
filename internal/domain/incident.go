package domain

import "time"

type Incident struct {
	IncidentID string
	OperatorID string
	Longitude  float64
	Latitude   float64
	Radius     uint8
	CreatedAt  time.Time
	DeletedAt  *time.Time
	UpdatedAt  *time.Time
}
