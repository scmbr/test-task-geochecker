package domain

import (
	"fmt"
	"time"
)

type Incident struct {
	IncidentID string
	OperatorID string
	Longitude  float64
	Latitude   float64
	Radius     uint16
	CreatedAt  time.Time
	DeletedAt  *time.Time
	UpdatedAt  *time.Time
}

func NewIncident(id, operatorID string, lat, lon float64, radius uint16) (*Incident, error) {
	if lat < -90 || lat > 90 {
		return nil, fmt.Errorf("invalid latitude")
	}
	if lon < -180 || lon > 180 {
		return nil, fmt.Errorf("invalid longitude")
	}
	if radius == 0 {
		return nil, fmt.Errorf("radius must be positive")
	}
	return &Incident{
		IncidentID: id,
		OperatorID: operatorID,
		Latitude:   lat,
		Longitude:  lon,
		Radius:     radius,
	}, nil
}
