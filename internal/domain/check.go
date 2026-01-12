package domain

import "time"

type Check struct {
	CheckID   string
	UserID    string
	Longitude float64
	Latitude  float64
	CreatedAt time.Time
}
