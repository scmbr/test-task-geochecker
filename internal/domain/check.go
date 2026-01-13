package domain

import (
	"fmt"
	"time"
)

type Check struct {
	CheckID   string
	UserID    string
	Longitude float64
	Latitude  float64
	CreatedAt time.Time
}

func NewCheck(id, userID string, lat, lon float64) (*Check, error) {
	if lat < -90 || lat > 90 {
		return nil, fmt.Errorf("latitude must be between -90 and 90")
	}
	if lon < -180 || lon > 180 {
		return nil, fmt.Errorf("longitude must be between -180 and 180")
	}
	if userID == "" {
		return nil, fmt.Errorf("userID cannot be empty")
	}
	return &Check{
		CheckID:   id,
		UserID:    userID,
		Latitude:  lat,
		Longitude: lon,
	}, nil
}
